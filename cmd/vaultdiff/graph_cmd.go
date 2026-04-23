package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// RunGraph fetches multiple secret versions and renders a change graph.
func RunGraph(args []string) error {
	fs := flag.NewFlagSet("graph", flag.ContinueOnError)
	path := fs.String("path", "", "KV secret path (required)")
	addr := fs.String("addr", os.Getenv("VAULT_ADDR"), "Vault address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	versionsFlag := fs.String("versions", "", "Comma-separated version numbers to graph (e.g. 1,2,3)")
	mask := fs.Bool("mask", true, "Mask secret values in output")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *path == "" {
		return fmt.Errorf("--path is required")
	}
	if *addr == "" {
		return fmt.Errorf("--addr is required")
	}
	if *token == "" {
		return fmt.Errorf("--token is required")
	}
	if *versionsFlag == "" {
		return fmt.Errorf("--versions is required (e.g. --versions=1,2,3)")
	}

	parts := strings.Split(*versionsFlag, ",")
	if len(parts) < 2 {
		return fmt.Errorf("--versions must specify at least two version numbers")
	}

	client, err := vault.NewClient(*addr, *token, "")
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	snapshots := map[int][]diff.Entry{}
	for i := 0; i < len(parts)-1; i++ {
		vA, err := strconv.Atoi(strings.TrimSpace(parts[i]))
		if err != nil {
			return fmt.Errorf("invalid version %q: %w", parts[i], err)
		}
		vB, err := strconv.Atoi(strings.TrimSpace(parts[i+1]))
		if err != nil {
			return fmt.Errorf("invalid version %q: %w", parts[i+1], err)
		}

		secA, err := client.ReadSecretVersion(*path, vA)
		if err != nil {
			return fmt.Errorf("read version %d: %w", vA, err)
		}
		secB, err := client.ReadSecretVersion(*path, vB)
		if err != nil {
			return fmt.Errorf("read version %d: %w", vB, err)
		}

		entries := diff.Compare(secA, secB)
		_ = *mask
		snapshots[vA] = entries
		if i == len(parts)-2 {
			snapshots[vB] = entries
		}
	}

	nodes := diff.BuildGraph(snapshots)
	fmt.Print(diff.FormatGraph(nodes))
	return nil
}
