package main

import (
	"fmt"
	"os"

	"github.com/yourusername/vaultdiff/internal/diff"
)

// FilterFlags holds CLI flags specific to filtered diff output.
type FilterFlags struct {
	OnlyChanged bool
	KeyPrefix   string
	ExcludeKeys []string
}

// RunFilteredDiff executes a diff between two secret versions and applies
// the provided filter options before rendering.
func RunFilteredDiff(cfg *Config, ff FilterFlags) error {
	client, err := newVaultClient(cfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	before, err := client.ReadSecretVersion(cfg.Path, cfg.VersionA)
	if err != nil {
		return fmt.Errorf("reading version %d: %w", cfg.VersionA, err)
	}

	after, err := client.ReadSecretVersion(cfg.Path, cfg.VersionB)
	if err != nil {
		return fmt.Errorf("reading version %d: %w", cfg.VersionB, err)
	}

	entries := diff.Compare(before, after)
	entries = diff.Filter(entries, diff.FilterOptions{
		OnlyChanged: ff.OnlyChanged,
		KeyPrefix:   ff.KeyPrefix,
		ExcludeKeys: ff.ExcludeKeys,
	})

	output := diff.Render(entries, diff.RenderOptions{MaskValues: cfg.Mask})
	fmt.Fprint(os.Stdout, output)
	return nil
}
