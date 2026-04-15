package main

import (
	"fmt"
	"os"
)

func main() {
	cfg, err := ParseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Usage: vaultdiff -addr <url> -token <tok> -path <secret-path> [options]")
		os.Exit(1)
	}

	if err := Run(cfg, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "vaultdiff: %v\n", err)
		os.Exit(1)
	}
}
