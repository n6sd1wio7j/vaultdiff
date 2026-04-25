package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// RunEnvDiff compares the same secret path across two Vault environments.
// It reads the secret at the given path from both Vault instances and prints
// a diff of the key-value pairs, optionally ignoring keys by prefix.
func RunEnvDiff(args []string) error {
	fs := flag.NewFlagSet("env", flag.ContinueOnError)

	path := fs.String("path", "", "secret path to compare (required)")
	addrA := fs.String("addr-a", "", "Vault address for env A (required)")
	tokenA := fs.String("token-a", "", "Vault token for env A (required)")
	addrB := fs.String("addr-b", "", "Vault address for env B (required)")
	tokenB := fs.String("token-b", "", "Vault token for env B (required)")
	envA := fs.String("env-a", "staging", "label for environment A")
	envB := fs.String("env-b", "production", "label for environment B")
	mount := fs.String("mount", "secret", "KV v2 mount path")
	ignore := fs.String("ignore-prefix", "", "comma-separated key prefixes to ignore")
	versionA := fs.Int("version-a", 0, "secret version for env A (0 = latest)")
	versionB := fs.Int("version-b", 0, "secret version for env B (0 = latest)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *path == "" {
		return fmt.Errorf("--path is required")
	}
	if *addrA == "" || *tokenA == "" {
		return fmt.Errorf("--addr-a and --token-a are required")
	}
	if *addrB == "" || *tokenB == "" {
		return fmt.Errorf("--addr-b and --token-b are required")
	}

	clientA, err := vault.NewClient(vault.Config{Address: *addrA, Token: *tokenA, Mount: *mount})
	if err != nil {
		return fmt.Errorf("client A: %w", err)
	}
	clientB, err := vault.NewClient(vault.Config{Address: *addrB, Token: *tokenB, Mount: *mount})
	if err != nil {
		return fmt.Errorf("client B: %w", err)
	}

	secrets A, err := clientA.ReadSecretVersion(*path, *versionA)
	if err != nil {
		return fmt.Errorf("read env A (%s): %w", *addrA, err)
	}
	secrets B, err := clientB.ReadSecretVersion(*path, *versionB)
	if err != nil {
		return fmt.Errorf("read env B (%s): %w", *addrB, err)
	}

	var ignorePrefixes []string
	if *ignore != "" {
		for _, p := range strings.Split(*ignore, ",") {
			if trimmed := strings.TrimSpace(p); trimmed != "" {
				ignorePrefixes = append(ignorePrefixes, trimmed)
			}
		}
	}

	opts := diff.EnvCompareOptions{
		EnvA:   *envA,
		EnvB:   *envB,
		Ignore: ignorePrefixes,
	}

	result := diff.CompareEnvs(secretsA, secretsB, opts)
	fmt.Fprint(os.Stdout, diff.FormatEnvDiff(result))
	return nil
}
