package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// RunBlame reads two versions of a secret and prints a blame report.
func RunBlame(args []string) error {
	fs := flag.NewFlagSet("blame", flag.ContinueOnError)
	path := fs.String("path", "", "secret path (required)")
	addr := fs.String("addr", os.Getenv("VAULT_ADDR"), "Vault address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	v1 := fs.Int("v1", 1, "first version")
	v2 := fs.Int("v2", 2, "second version (blamed)")

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

	client, err := vault.NewClient(*addr, *token, "")
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secrets1, err := client.ReadSecretVersion(*path, *v1)
	if err != nil {
		return fmt.Errorf("read version %d: %w", *v1, err)
	}
	secrets2, err := client.ReadSecretVersion(*path, *v2)
	if err != nil {
		return fmt.Errorf("read version %d: %w", *v2, err)
	}

	entries := diff.Compare(secrets1, secrets2)
	report := diff.Blame(*path, entries, *v2, time.Now().UTC())
	fmt.Print(diff.FormatBlame(report))
	return nil
}
