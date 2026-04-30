package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/hermes-93/go-utils/internal/port"
)

func Portcheck(args []string) {
	fs := flag.NewFlagSet("portcheck", flag.ExitOnError)
	timeout := fs.Int("timeout", 5, "TCP dial timeout in seconds")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: devops portcheck [flags] <host:port> [host:port...]")
		fmt.Fprintln(os.Stderr, "\nFlags:")
		fs.PrintDefaults()
	}
	fs.Parse(args) //nolint:errcheck

	if fs.NArg() < 1 {
		fs.Usage()
		os.Exit(1)
	}

	failed := 0
	for _, addr := range fs.Args() {
		result := port.Check(addr, *timeout)
		fmt.Println(result.String())
		if !result.Open {
			failed++
		}
	}
	if failed > 0 {
		os.Exit(1)
	}
}
