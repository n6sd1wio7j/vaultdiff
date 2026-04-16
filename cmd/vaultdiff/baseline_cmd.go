package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// RunBaseline handles the `baseline` subcommand: save or compare against a snapshot.
func RunBaseline(args []string) error {
	fs := flag.NewFlagSet("baseline", flag.ContinueOnError)
	path := fs.String("path", "", "secret path in Vault (required)")
	addr := fs.String("addr", os.Getenv("VAULT_ADDR"), "Vault address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	version := fs.Int("version", 0, "secret version (0 = latest)")
	baselineFile := fs.String("file", "baseline.json", "path to baseline file")
	save := fs.Bool("save", false, "save current state as baseline")

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

	client, err := vault.NewClient(*addr, *token, "")
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	current, err := client.ReadSecretVersion(*path, *version)
	if err != nil {
		return fmt.Errorf("read secret: %w", err)
	}

	if *save {
		entries := diff.Compare(current, current)
		if err := diff.SaveBaseline(*baselineFile, *path, entries); err != nil {
			return fmt.Errorf("save baseline: %w", err)
		}
		fmt.Printf("Baseline saved to %s\n", *baselineFile)
		return nil
	}

	baseline, err := diff.LoadBaseline(*baselineFile)
	if err != nil {
		return fmt.Errorf("load baseline: %w", err)
	}

	entries := diff.CompareToBaseline(baseline, current)
	fmt.Print(diff.Render(entries, diff.RenderOptions{ShowUnchanged: false, MaskValues: false}))
	return nil
}
