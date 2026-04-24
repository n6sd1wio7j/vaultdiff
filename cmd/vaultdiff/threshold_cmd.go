package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// RunThreshold compares two secret versions and exits non-zero if thresholds are breached.
func RunThreshold(args []string) error {
	fs := flag.NewFlagSet("threshold", flag.ContinueOnError)

	path := fs.String("path", "", "KV secret path (required)")
	addr := fs.String("addr", os.Getenv("VAULT_ADDR"), "Vault address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	v1 := fs.Int("v1", 0, "first version (0 = latest-1)")
	v2 := fs.Int("v2", 0, "second version (0 = latest)")
	maxScore := fs.Float64("max-score", 100.0, "maximum allowed weighted drift score")
	maxAdded := fs.Int("max-added", 50, "maximum allowed added keys")
	maxRemoved := fs.Int("max-removed", 50, "maximum allowed removed keys")
	maxModified := fs.Int("max-modified", 50, "maximum allowed modified keys")

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

	secA, err := client.ReadSecretVersion(*path, *v1)
	if err != nil {
		return fmt.Errorf("reading version %d: %w", *v1, err)
	}
	secB, err := client.ReadSecretVersion(*path, *v2)
	if err != nil {
		return fmt.Errorf("reading version %d: %w", *v2, err)
	}

	entries := diff.Compare(secA, secB)
	score := diff.ScoreDrift(entries)

	opts := diff.ThresholdOptions{
		MaxScore:    *maxScore,
		MaxAdded:    *maxAdded,
		MaxRemoved:  *maxRemoved,
		MaxModified: *maxModified,
	}

	result := diff.CheckThreshold(score, opts)
	fmt.Print(diff.FormatThresholdResult(result))

	if result.Breached {
		os.Exit(1)
	}
	return nil
}
