package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// RunLint reads two secret versions and reports lint violations.
func RunLint(args []string) error {
	fs := flag.NewFlagSet("lint", flag.ContinueOnError)
	path := fs.String("path", "", "secret path (required)")
	addr := fs.String("addr", os.Getenv("VAULT_ADDR"), "Vault address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	v1 := fs.Int("v1", 0, "first version")
	v2 := fs.Int("v2", 0, "second version")

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

	client, err := vault.NewClient(vault.Config{Address: *addr, Token: *token})
	if err != nil {
		return err
	}

	secA, err := client.ReadSecretVersion(*path, *v1)
	if err != nil {
		return fmt.Errorf("read v%d: %w", *v1, err)
	}
	secB, err := client.ReadSecretVersion(*path, *v2)
	if err != nil {
		return fmt.Errorf("read v%d: %w", *v2, err)
	}

	entries := diff.Compare(secA, secB)
	violations := diff.Lint(entries, diff.DefaultLintRules())
	fmt.Print(diff.FormatLint(violations))
	if len(violations) > 0 {
		os.Exit(1)
	}
	return nil
}
