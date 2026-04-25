package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// RunHeatmap fetches a chain of secret versions and prints a change-frequency
// heatmap showing which keys change most often.
func RunHeatmap(args []string) error {
	fs := flag.NewFlagSet("heatmap", flag.ContinueOnError)

	path := fs.String("path", "", "secret path (required)")
	addr := fs.String("addr", os.Getenv("VAULT_ADDR"), "Vault address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	mount := fs.String("mount", "secret", "KV v2 mount")
	topN := fs.Int("top", 10, "show top N keys by change frequency (0 = all)")
	showAll := fs.Bool("show-all", false, "include keys with zero changes")

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
		Mount:   *mount,
	})
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	versions, err := vault.ListVersions(client, *path)
	if err != nil {
		return fmt.Errorf("list versions: %w", err)
	}
	if len(versions) < 2 {
		fmt.Println("not enough versions to build heatmap")
		return nil
	}

	chain, err := diff.BuildChain(client, *path, versions)
	if err != nil {
		return fmt.Errorf("build chain: %w", err)
	}

	opts := diff.DefaultHeatmapOptions()
	opts.TopN = *topN
	opts.ShowAll = *showAll

	entries := diff.BuildHeatmap(chain, opts)
	fmt.Print(diff.FormatHeatmap(entries))
	return nil
}
