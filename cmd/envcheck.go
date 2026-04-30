package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/hermes-93/go-utils/internal/env"
)

func Envcheck(args []string) {
	fs := flag.NewFlagSet("envcheck", flag.ExitOnError)
	file  := fs.String("file", "", "Path to .env.example listing required variable names")
	quiet := fs.Bool("quiet", false, "Only print missing variables (suppress OK lines)")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: devops envcheck [flags] [VAR...]")
		fmt.Fprintln(os.Stderr, "\nFlags:")
		fs.PrintDefaults()
	}
	fs.Parse(args) //nolint:errcheck

	vars := fs.Args()
	if *file != "" {
		fromFile, err := env.ParseFile(*file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		vars = append(vars, fromFile...)
	}

	if len(vars) == 0 {
		fmt.Fprintln(os.Stderr, "error: specify variable names or --file <path>")
		fs.Usage()
		os.Exit(1)
	}

	missing := env.Check(vars)
	exitCode := 0
	for _, v := range vars {
		if contains(missing, v) {
			fmt.Printf("MISSING %s\n", v)
			exitCode = 1
		} else if !*quiet {
			fmt.Printf("OK      %s\n", v)
		}
	}
	os.Exit(exitCode)
}

func contains(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}
