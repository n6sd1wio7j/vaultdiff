package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// ExportFlags holds CLI flags for the export subcommand.
type ExportFlags struct {
	Path    string
	Addr    string
	Token   string
	Mount   string
	VerA    int
	VerB    int
	Format  string
	Mask    bool
	Output  string
}

// RunExport executes the export subcommand.
func RunExport(args []string) error {
	fs := flag.NewFlagSet("export", flag.ContinueOnError)
	var f ExportFlags

	fs.StringVar(&f.Path, "path", "", "secret path (required)")
	fs.StringVar(&f.Addr, "addr", os.Getenv("VAULT_ADDR"), "Vault address")
	fs.StringVar(&f.Token, "token", os.Getenv("VAULT_TOKEN"), "Vault token")
	fs.StringVar(&f.Mount, "mount", "secret", "KV mount path")
	fs.IntVar(&f.VerA, "ver-a", 1, "first version")
	fs.IntVar(&f.VerB, "ver-b", 2, "second version")
	fs.StringVar(&f.Format, "format", "json", "export format: json, csv, text")
	fs.BoolVar(&f.Mask, "mask", true, "mask secret values in output")
	fs.StringVar(&f.Output, "output", "", "output file path (default: stdout)")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if f.Path == "" {
		return fmt.Errorf("--path is required")
	}

	client, err := vault.NewClient(vault.Config{Address: f.Addr, Token: f.Token, Mount: f.Mount})
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secA, err := client.ReadSecretVersion(f.Path, f.VerA)
	if err != nil {
		return fmt.Errorf("reading version %d: %w", f.VerA, err)
	}
	secB, err := client.ReadSecretVersion(f.Path, f.VerB)
	if err != nil {
		return fmt.Errorf("reading version %d: %w", f.VerB, err)
	}

	entries := diff.Compare(secA, secB)

	w := os.Stdout
	if f.Output != "" {
		file, err := os.Create(f.Output)
		if err != nil {
			return fmt.Errorf("opening output file: %w", err)
		}
		defer file.Close()
		w = file
	}

	return diff.Export(entries, diff.ExportOptions{
		Format:      diff.ExportFormat(f.Format),
		MaskSecrets: f.Mask,
	}, w)
}
