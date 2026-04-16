package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// RunScore computes and prints a weighted drift score between two secret versions.
func RunScore(args []string) error {
	fs := flag.NewFlagSet("score", flag.ContinueOnError)
	path := fs.String("path", "", "secret path (required)")
	addr := fs.String("addr", os.Getenv("VAULT_ADDR"), "Vault address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	mount := fs.String("mount", "secret", "KV v2 mount")
	versionA := fs.Int("version-a", 0, "first version (0 = previous)")
	versionB := fs.Int("version-b", 0, "second version (0 = latest)")

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

	client, err := vault.NewClient(*addr, *token, *mount)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secretA, err := client.ReadSecretVersion(*path, *versionA)
	if err != nil {
		return fmt.Errorf("reading version-a: %w", err)
	}
	secretB, err := client.ReadSecretVersion(*path, *versionB)
	if err != nil {
		return fmt.Errorf("reading version-b: %w", err)
	}

	entries := diff.Compare(secretA, secretB)
	ds := diff.ScoreDrift(entries)
	fmt.Println(diff.FormatScore(ds))
	return nil
}
