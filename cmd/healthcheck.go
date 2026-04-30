package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/hermes-93/go-utils/internal/health"
)

func Healthcheck(args []string) {
	fs := flag.NewFlagSet("healthcheck", flag.ExitOnError)
	timeout  := fs.Int("timeout", 5, "HTTP timeout in seconds")
	retries  := fs.Int("retries", 3, "Number of retries on non-2xx")
	interval := fs.Int("interval", 2, "Seconds between retries")
	expected := fs.Int("expected", 200, "Expected HTTP status code")
	jsonOut  := fs.Bool("json", false, "Output result as JSON")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: devops healthcheck [flags] <url>")
		fmt.Fprintln(os.Stderr, "\nFlags:")
		fs.PrintDefaults()
	}
	fs.Parse(args) //nolint:errcheck

	if fs.NArg() != 1 {
		fs.Usage()
		os.Exit(1)
	}

	result := health.Check(fs.Arg(0), health.Options{
		Timeout:       *timeout,
		Retries:       *retries,
		RetryInterval: *interval,
		ExpectedCode:  *expected,
	})

	if *jsonOut {
		printHealthJSON(result)
	} else {
		fmt.Println(result.String())
	}
	if !result.Healthy {
		os.Exit(1)
	}
}

func printHealthJSON(r health.Result) {
	out := struct {
		URL       string `json:"url"`
		Healthy   bool   `json:"healthy"`
		Status    int    `json:"status_code,omitempty"`
		LatencyMs int64  `json:"latency_ms"`
		Error     string `json:"error,omitempty"`
		Attempts  int    `json:"attempts"`
	}{
		URL:       r.URL,
		Healthy:   r.Healthy,
		Status:    r.StatusCode,
		LatencyMs: r.Latency.Milliseconds(),
		Error:     r.Error,
		Attempts:  r.Attempts,
	}
	data, _ := json.Marshal(out)
	fmt.Println(string(data))
}
