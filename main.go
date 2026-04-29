package main

import (
	"fmt"
	"os"

	"github.com/hermes-93/go-utils/cmd"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "healthcheck":
		cmd.Healthcheck(os.Args[2:])
	case "portcheck":
		cmd.Portcheck(os.Args[2:])
	case "logparse":
		cmd.Logparse(os.Args[2:])
	case "envcheck":
		cmd.Envcheck(os.Args[2:])
	case "waitfor":
		cmd.Waitfor(os.Args[2:])
	case "version", "--version", "-v":
		fmt.Println("devops", version)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`devops — Go CLI utilities for DevOps automation

Usage:
  devops <command> [flags]

Commands:
  healthcheck   Check HTTP endpoint health with retries
  portcheck     Check TCP port availability
  logparse      Parse and filter JSON structured logs
  envcheck      Validate required environment variables
  waitfor       Wait for a service to become available

Run 'devops <command> --help' for more information.`)
}
