package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/hermes-93/go-utils/internal/logparse"
)

func Logparse(args []string) {
	fs := flag.NewFlagSet("logparse", flag.ExitOnError)
	level := fs.String("level", "", "Minimum log level to show (debug|info|warn|error|fatal)")
	format := fs.String("format", "text", "Output format: text or json")
	fieldsRaw := fs.String("fields", "", "Comma-separated extra fields to display")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: devops logparse [flags] [file]")
		fmt.Fprintln(os.Stderr, "Reads from stdin if no file is given.")
		fmt.Fprintln(os.Stderr, "\nFlags:")
		fs.PrintDefaults()
	}
	fs.Parse(args) //nolint:errcheck

	var fields []string
	if *fieldsRaw != "" {
		for _, f := range strings.Split(*fieldsRaw, ",") {
			if f = strings.TrimSpace(f); f != "" {
				fields = append(fields, f)
			}
		}
	}

	var scanner *bufio.Scanner
	if fs.NArg() >= 1 {
		f, err := os.Open(fs.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		scanner = bufio.NewScanner(f)
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	p := logparse.Parser{MinLevel: *level, Fields: fields, Format: *format}
	p.Process(scanner, os.Stdout)
}
