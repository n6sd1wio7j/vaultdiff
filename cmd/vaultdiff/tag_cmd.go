package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// RunTag fetches two secret versions, diffs them, and prints tagged output.
func RunTag(args []string) error {
	fs := flag.NewFlagSet("tag", flag.ContinueOnError)
	path := fs.String("path", "", "secret path (required)")
	addr := fs.String("addr", "", "Vault address (required)")
	token := fs.String("token", "", "Vault token (required)")
	v1 := fs.Int("v1", 1, "first version")
	v2 := fs.Int("v2", 2, "second version")
	env := fs.String("env", "", "environment tag")
	version := fs.String("version", "", "version label tag")

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

	old, err := client.ReadSecretVersion(*path, *v1)
	if err != nil {
		return fmt.Errorf("reading v%d: %w", *v1, err)
	}
	new, err := client.ReadSecretVersion(*path, *v2)
	if err != nil {
		return fmt.Errorf("reading v%d: %w", *v2, err)
	}

	entries := diff.Compare(old, new)
	tagged := diff.TagEntries(entries, diff.TagOptions{
		Env:     *env,
		Version: *version,
	})

	for _, te := range tagged {
		fmt.Fprintln(os.Stdout, diff.FormatTagged(te))
	}
	return nil
}
