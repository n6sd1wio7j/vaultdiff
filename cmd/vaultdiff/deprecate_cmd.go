package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// RunDeprecate fetches two secret versions and reports deprecated keys.
func RunDeprecate(args []string) error {
	fs := flag.NewFlagSet("deprecate", flag.ContinueOnError)

	path := fs.String("path", "", "KV secret path (required)")
	addr := fs.String("addr", "", "Vault address (required)")
	token := fs.String("token", "", "Vault token (required)")
	versionA := fs.Int("version-a", 0, "First version to compare (0 = previous)")
	versionB := fs.Int("version-b", 0, "Second version to compare (0 = latest)")
	keys := fs.String("deprecated-keys", "", "Comma-separated list of deprecated key names")
	prefixes := fs.String("deprecated-prefixes", "", "Comma-separated list of deprecated key prefixes")
	includeUnchanged := fs.Bool("include-unchanged", false, "Also flag unchanged deprecated keys")

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

	secretA, err := client.ReadSecretVersion(*path, *versionA)
	if err != nil {
		return fmt.Errorf("reading version %d: %w", *versionA, err)
	}
	secretB, err := client.ReadSecretVersion(*path, *versionB)
	if err != nil {
		return fmt.Errorf("reading version %d: %w", *versionB, err)
	}

	entries := diff.Compare(secretA, secretB)

	opts := diff.DefaultDeprecateOptions()
	opts.IncludeUnchanged = *includeUnchanged
	if *keys != "" {
		opts.DeprecatedKeys = strings.Split(*keys, ",")
	}
	if *prefixes != "" {
		opts.DeprecatedPrefixes = strings.Split(*prefixes, ",")
	}

	results := diff.DetectDeprecated(entries, opts)
	fmt.Fprint(os.Stdout, diff.FormatDeprecated(results))
	return nil
}
