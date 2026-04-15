package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/your-org/vaultdiff/internal/audit"
	"github.com/your-org/vaultdiff/internal/diff"
	"github.com/your-org/vaultdiff/internal/vault"
)

// Run executes the vaultdiff logic using the provided Config and writer.
func Run(cfg *Config, out io.Writer) error {
	client, err := vault.NewClient(vault.Options{
		Address: cfg.Address,
		Token:   cfg.Token,
		Mount:   cfg.Mount,
	})
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secretA, err := client.ReadSecretVersion(cfg.Path, cfg.VersionA)
	if err != nil {
		return fmt.Errorf("reading version %d of %q: %w", cfg.VersionA, cfg.Path, err)
	}

	secretB, err := client.ReadSecretVersion(cfg.Path, cfg.VersionB)
	if err != nil {
		return fmt.Errorf("reading version %d of %q: %w", cfg.VersionB, cfg.Path, err)
	}

	entries := diff.Compare(secretA, secretB)

	// Audit logging
	var logger *audit.Logger
	if cfg.AuditLog != "" {
		fl, ferr := audit.OpenFileLogger(cfg.AuditLog)
		if ferr != nil {
			return fmt.Errorf("opening audit log: %w", ferr)
		}
		logger = fl
	} else {
		logger = audit.NewLogger(os.Stdout)
	}

	if err := logger.Record(cfg.Path, cfg.VersionA, cfg.VersionB, entries); err != nil {
		return fmt.Errorf("audit log: %w", err)
	}

	if cfg.OutputJSON {
		return json.NewEncoder(out).Encode(entries)
	}

	report := diff.Report(cfg.Path, cfg.VersionA, cfg.VersionB, entries)
	render := diff.Render(entries, diff.RenderOptions{
		ShowUnchanged: cfg.ShowAll,
		MaskValues:    cfg.MaskValues,
	})

	fmt.Fprintln(out, report)
	fmt.Fprintln(out, render)
	return nil
}
