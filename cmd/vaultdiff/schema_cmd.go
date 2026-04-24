package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// RunSchema validates the current secret version at path against schema rules.
func RunSchema(args []string) error {
	fs := flag.NewFlagSet("schema", flag.ContinueOnError)

	path := fs.String("path", "", "Vault KV path to validate (required)")
	addr := fs.String("addr", os.Getenv("VAULT_ADDR"), "Vault address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	mount := fs.String("mount", "secret", "KV v2 mount path")
	version := fs.Int("version", 0, "Secret version (0 = latest)")

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

	client, err := vault.NewClient(vault.Config{
		Address: *addr,
		Token:   *token,
		Mount:   *mount,
	})
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secrets, err := client.ReadSecretVersion(*path, *version)
	if err != nil {
		return fmt.Errorf("read secret: %w", err)
	}

	// Build entries treating all current keys as "added" for schema validation
	var entries []diff.Entry
	for k, v := range secrets {
		entries = append(entries, diff.Entry{
			Key:      k,
			Status:   diff.StatusAdded,
			NewValue: v,
		})
	}

	rules := diff.DefaultSchemaRules()
	violations := diff.ValidateSchema(entries, rules)
	fmt.Print(diff.FormatSchemaViolations(violations))

	if len(violations) > 0 {
		os.Exit(1)
	}
	return nil
}
