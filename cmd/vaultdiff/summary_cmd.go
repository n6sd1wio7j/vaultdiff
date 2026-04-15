package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// RunSummary fetches two secret versions, computes a diff, and prints a
// concise drift summary. With --json the output is machine-readable.
func RunSummary(args []string) error {
	fs := flag.NewFlagSet("summary", flag.ContinueOnError)
	addr := fs.String("addr", os.Getenv("VAULT_ADDR"), "Vault address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	mount := fs.String("mount", "secret", "KV v2 mount path")
	path := fs.String("path", "", "Secret path (required)")
	v1 := fs.Int("v1", 0, "First version (0 = latest)")
	v2 := fs.Int("v2", 0, "Second version (0 = latest)")
	jsonOut := fs.Bool("json", false, "Output summary as JSON")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *path == "" {
		return fmt.Errorf("--path is required")
	}
	if *addr == "" {
		return fmt.Errorf("--addr or VAULT_ADDR is required")
	}
	if *token == "" {
		return fmt.Errorf("--token or VAULT_TOKEN is required")
	}

	client, err := vault.NewClient(*addr, *token, *mount)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secA, err := client.ReadSecretVersion(*path, *v1)
	if err != nil {
		return fmt.Errorf("reading version %d: %w", *v1, err)
	}
	secB, err := client.ReadSecretVersion(*path, *v2)
	if err != nil {
		return fmt.Errorf("reading version %d: %w", *v2, err)
	}

	entries := diff.Compare(secA, secB)
	s := diff.Summarize(entries)

	if *jsonOut {
		return json.NewEncoder(os.Stdout).Encode(s)
	}
	fmt.Println(s.String())
	if s.HasDrift() {
		os.Exit(1)
	}
	return nil
}
