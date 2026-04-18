package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// RunPromote diffs two paths and prints a promotion plan.
func RunPromote(args []string) error {
	fs := flag.NewFlagSet("promote", flag.ContinueOnError)
	srcPath := fs.String("src", "", "source secret path (required)")
	dstPath := fs.String("dst", "", "destination secret path (required)")
	addr := fs.String("addr", os.Getenv("VAULT_ADDR"), "Vault address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	mount := fs.String("mount", "secret", "KV mount path")
	srcVer := fs.Int("src-version", 0, "source secret version (0 = latest)")
	dstVer := fs.Int("dst-version", 0, "destination secret version (0 = latest)")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *srcPath == "" {
		return fmt.Errorf("--src is required")
	}
	if *dstPath == "" {
		return fmt.Errorf("--dst is required")
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

	srcSecrets, err := client.ReadSecretVersion(*srcPath, *srcVer)
	if err != nil {
		return fmt.Errorf("read src: %w", err)
	}
	dstSecrets, err := client.ReadSecretVersion(*dstPath, *dstVer)
	if err != nil {
		return fmt.Errorf("read dst: %w", err)
	}

	entries := diff.Compare(srcSecrets, dstSecrets)
	plan := diff.BuildPromotePlan(*srcPath, *dstPath, entries)
	fmt.Print(diff.FormatPromotePlan(plan))
	return nil
}
