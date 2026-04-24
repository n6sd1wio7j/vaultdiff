package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// RunNormalizedDiff reads two secret versions, normalizes both, then diffs them.
func RunNormalizedDiff(args []string) error {
	fs := flag.NewFlagSet("normalize", flag.ContinueOnError)

	path := fs.String("path", "", "Secret path (required)")
	addr := fs.String("addr", os.Getenv("VAULT_ADDR"), "Vault address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	versionA := fs.Int("version-a", 0, "First version to compare (0 = latest-1)")
	versionB := fs.Int("version-b", 0, "Second version to compare (0 = latest)")
	lowerKeys := fs.Bool("lowercase-keys", false, "Normalize keys to lowercase before diff")
	noTrim := fs.Bool("no-trim", false, "Disable whitespace trimming")
	showUnchanged := fs.Bool("show-unchanged", false, "Show unchanged keys in output")

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

	secA, err := client.ReadSecretVersion(*path, *versionA)
	if err != nil {
		return fmt.Errorf("reading version-a: %w", err)
	}
	secB, err := client.ReadSecretVersion(*path, *versionB)
	if err != nil {
		return fmt.Errorf("reading version-b: %w", err)
	}

	opts := diff.NormalizeOptions{
		TrimSpace:     !*noTrim,
		LowercaseKeys: *lowerKeys,
		StripNonPrint: true,
	}
	nA, nB := diff.NormalizeSecrets(secA, secB, opts)

	entries := diff.Compare(nA, nB)
	renderOpts := diff.RenderOptions{
		ShowUnchanged: *showUnchanged,
		MaskSecrets:   true,
	}
	fmt.Print(diff.Render(entries, renderOpts))
	return nil
}
