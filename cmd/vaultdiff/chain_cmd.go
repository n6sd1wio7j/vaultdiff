package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// RunChain reads multiple versions of a secret and renders a diff chain.
func RunChain(args []string) error {
	fs := flag.NewFlagSet("chain", flag.ContinueOnError)
	path := fs.String("path", "", "Secret path (required)")
	addr := fs.String("addr", "", "Vault address (required)")
	token := fs.String("token", "", "Vault token (required)")
	mount := fs.String("mount", "secret", "KV v2 mount point")
	versionsFlag := fs.String("versions", "", "Comma-separated version numbers to chain, e.g. 1,2,3")

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
	if *versionsFlag == "" {
		return fmt.Errorf("--versions is required (e.g. 1,2,3)")
	}

	parts := strings.Split(*versionsFlag, ",")
	var versionNums []int
	for _, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return fmt.Errorf("invalid version number %q: %w", p, err)
		}
		versionNums = append(versionNums, n)
	}
	if len(versionNums) < 2 {
		return fmt.Errorf("at least two versions are required to build a chain")
	}

	client, err := vault.NewClient(vault.Config{
		Address: *addr,
		Token:   *token,
		Mount:   *mount,
	})
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	var versionMaps []map[string]string
	for _, v := range versionNums {
		secrets, err := client.ReadSecretVersion(*path, v)
		if err != nil {
			return fmt.Errorf("reading version %d: %w", v, err)
		}
		versionMaps = append(versionMaps, secrets)
	}

	chain := diff.BuildChain(*path, versionMaps)
	fmt.Fprint(os.Stdout, diff.FormatChain(chain))
	return nil
}
