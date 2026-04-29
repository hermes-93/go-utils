package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/hermes-93/go-utils/internal/health"
)

func Healthcheck(args []string) {
	fs := flag.NewFlagSet("healthcheck", flag.ExitOnError)
	timeout := fs.Int("timeout", 5, "HTTP timeout in seconds")
	retries := fs.Int("retries", 3, "Number of retries on non-2xx")
	interval := fs.Int("interval", 2, "Seconds between retries")
	expected := fs.Int("expected", 200, "Expected HTTP status code")
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
	fmt.Println(result.String())
	if !result.Healthy {
		os.Exit(1)
	}
}
