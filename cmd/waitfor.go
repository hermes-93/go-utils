package cmd

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

func Waitfor(args []string) {
	fs := flag.NewFlagSet("waitfor", flag.ExitOnError)
	timeout := fs.Int("timeout", 60, "Max wait time in seconds")
	interval := fs.Int("interval", 2, "Poll interval in seconds")
	useHTTP := fs.Bool("http", false, "Treat address as HTTP URL (check for non-5xx response)")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: devops waitfor [flags] <address>")
		fmt.Fprintln(os.Stderr, "  TCP:  devops waitfor postgres:5432")
		fmt.Fprintln(os.Stderr, "  HTTP: devops waitfor --http http://api/health")
		fmt.Fprintln(os.Stderr, "\nFlags:")
		fs.PrintDefaults()
	}
	fs.Parse(args) //nolint:errcheck

	if fs.NArg() != 1 {
		fs.Usage()
		os.Exit(1)
	}

	addr := fs.Arg(0)
	deadline := time.Now().Add(time.Duration(*timeout) * time.Second)
	fmt.Printf("Waiting for %s (timeout %ds)...\n", addr, *timeout)

	for {
		ready := false
		if *useHTTP {
			ready = tryHTTP(addr)
		} else {
			ready = tryTCP(addr)
		}

		if ready {
			fmt.Printf("READY %s\n", addr)
			return
		}

		if time.Now().After(deadline) {
			fmt.Fprintf(os.Stderr, "TIMEOUT: %s not ready after %ds\n", addr, *timeout)
			os.Exit(1)
		}
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}

func tryTCP(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func tryHTTP(url string) bool {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode < 500
}
