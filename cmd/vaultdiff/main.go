package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: vaultdiff <command> [flags]")
		fmt.Fprintln(os.Stderr, "commands: diff, versions, filter, export, summary, watch")
		os.Exit(1)
	}
	cmd, args := os.Args[1], os.Args[2:]
	var err error
	switch cmd {
	case "diff":
		err = Run(args)
	case "versions":
		err = RunListVersions(args)
	case "filter":
		err = RunFilteredDiff(args)
	case "export":
		err = RunExport(args)
	case "summary":
		err = RunSummary(args)
	case "watch":
		err = RunWatch(args)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
