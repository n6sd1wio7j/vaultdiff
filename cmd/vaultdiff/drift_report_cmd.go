package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/your-org/vaultdiff/internal/diff"
	"github.com/your-org/vaultdiff/internal/vault"
)

// RunDriftReport reads two secret versions and emits a full drift report.
func RunDriftReport(args []string) error {
	fs := flag.NewFlagSet("drift-report", flag.ContinueOnError)

	path := fs.String("path", "", "Secret path (required)")
	addr := fs.String("addr", os.Getenv("VAULT_ADDR"), "Vault address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	mount := fs.String("mount", "secret", "KV mount")
	versionA := fs.Int("version-a", 0, "First version (required)")
	versionB := fs.Int("version-b", 0, "Second version (required)")
	env := fs.String("env", "", "Environment label (optional)")
	showScore := fs.Bool("score", true, "Include drift score")
	showSummary := fs.Bool("summary", true, "Include summary line")

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
	if *versionA == 0 || *versionB == 0 {
		return fmt.Errorf("--version-a and --version-b are required")
	}

	client, err := vault.NewClient(vault.Config{
		Address: *addr,
		Token:   *token,
		Mount:   *mount,
	})
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secA, err := client.ReadSecretVersion(*path, *versionA)
	if err != nil {
		return fmt.Errorf("read version %d: %w", *versionA, err)
	}
	secB, err := client.ReadSecretVersion(*path, *versionB)
	if err != nil {
		return fmt.Errorf("read version %d: %w", *versionB, err)
	}

	entries := diff.Compare(secA, secB)

	opts := diff.DriftReportOptions{
		Env:            *env,
		Path:           *path,
		VersionA:       *versionA,
		VersionB:       *versionB,
		IncludeScore:   *showScore,
		IncludeSummary: *showSummary,
	}

	report := diff.BuildDriftReport(entries, opts)
	fmt.Print(diff.FormatDriftReport(report, opts))
	return nil
}
