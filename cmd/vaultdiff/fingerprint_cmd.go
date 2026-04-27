package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// RunFingerprint computes and prints a stable fingerprint for a secret version.
func RunFingerprint(args []string) error {
	fs := flag.NewFlagSet("fingerprint", flag.ContinueOnError)

	path := fs.String("path", "", "Vault secret path (required)")
	addr := fs.String("addr", os.Getenv("VAULT_ADDR"), "Vault address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	version := fs.Int("version", 0, "Secret version to fingerprint (0 = latest)")
	mask := fs.Bool("mask", false, "Hash values before fingerprinting")
	skipUnchanged := fs.Bool("skip-unchanged", false, "Exclude unchanged keys from fingerprint")

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

	client, err := vault.NewClient(vault.Config{
		Address: *addr,
		Token:   *token,
	})
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secrets, err := client.ReadSecretVersion(*path, *version)
	if err != nil {
		return fmt.Errorf("read secret: %w", err)
	}

	// Build entries treating all keys as unchanged (snapshot fingerprint)
	entries := make([]diff.Entry, 0, len(secrets))
	for k, v := range secrets {
		entries = append(entries, diff.Entry{
			Key:      k,
			NewValue: v,
			OldValue: v,
			Status:   diff.StatusUnchanged,
		})
	}

	opts := diff.DefaultFingerprintOptions()
	opts.MaskValues = *mask
	opts.IncludeUnchanged = !*skipUnchanged

	result := diff.ComputeFingerprint(*path, *version, entries, opts)
	fmt.Print(diff.FormatFingerprint(result))
	return nil
}
