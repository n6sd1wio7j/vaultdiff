package main

import (
	"flag"
	"fmt"
	"os"
)

// Config holds all CLI flag values for a vaultdiff run.
type Config struct {
	Address     string
	Token       string
	Mount       string
	Path        string
	VersionA    int
	VersionB    int
	ShowAll     bool
	MaskValues  bool
	AuditLog    string
	OutputJSON  bool
}

// ParseFlags parses command-line flags and returns a Config.
// It exits with a usage message if required flags are missing.
func ParseFlags(args []string) (*Config, error) {
	fs := flag.NewFlagSet("vaultdiff", flag.ContinueOnError)

	cfg := &Config{}

	fs.StringVar(&cfg.Address, "addr", os.Getenv("VAULT_ADDR"), "Vault server address (env: VAULT_ADDR)")
	fs.StringVar(&cfg.Token, "token", os.Getenv("VAULT_TOKEN"), "Vault token (env: VAULT_TOKEN)")
	fs.StringVar(&cfg.Mount, "mount", "secret", "KV v2 mount path")
	fs.StringVar(&cfg.Path, "path", "", "Secret path to diff (required)")
	fs.IntVar(&cfg.VersionA, "version-a", 0, "First version to compare (0 = previous)")
	fs.IntVar(&cfg.VersionB, "version-b", 0, "Second version to compare (0 = latest)")
	fs.BoolVar(&cfg.ShowAll, "show-all", false, "Show unchanged keys as well")
	fs.BoolVar(&cfg.MaskValues, "mask", true, "Mask secret values in output")
	fs.StringVar(&cfg.AuditLog, "audit-log", "", "Path to audit log file (empty = stdout)")
	fs.BoolVar(&cfg.OutputJSON, "json", false, "Output diff as JSON")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if cfg.Path == "" {
		return nil, fmt.Errorf("flag -path is required")
	}
	if cfg.Address == "" {
		return nil, fmt.Errorf("vault address is required (use -addr or VAULT_ADDR)")
	}
	if cfg.Token == "" {
		return nil, fmt.Errorf("vault token is required (use -token or VAULT_TOKEN)")
	}

	return cfg, nil
}
