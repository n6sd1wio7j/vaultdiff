package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/your-org/vaultdiff/internal/diff"
t"github.com/your/internal/vault"
)

 polls two secret versions and prints drift reportsargs {
	fs("watch", flag.ContinueOnError)
	addr := fs.String("addr", os.Getenv("VAULT_ADDR"), "Vault address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	path := fs.String("path", "", "Secret path")
	versionA := fs.Int("version-a", 0, "First version")
	versionB := fs.Int("version-b", 0, "Second version")
	interval := fs.Duration("interval", 30*time.Second, "Poll interval")
	mask := fs.Bool("mask", true, "Mask secret values")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *path == "" {
		return fmt.Errorf("--path is required")
	}
	c, err := vault.NewClient(*addr, *token, "")
	if err != nil {
		return err
	}
	fetch := func() ([]diff.Entry, error) {
		a, err := c.ReadSecretVersion(*path, *versionA)
		if err != nil {
			return nil, err
		}
		b, err := c.ReadSecretVersion(*path, *versionB)
		if err != nil {
			return nil, err
		}
		return diff.Compare(a, b), nil
	}
	w := diff.NewWatcher(fetch, diff.WatchOptions{Interval: *interval})
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	fmt.Fprintf(os.Stdout, "Watching %s every %s (Ctrl+C to stop)\n", *path, *interval)
	for result := range w.Run(ctx) {
		fmt.Fprintf(os.Stdout, "\n[%s] drift=%v\n", result.CheckedAt.Format(time.RFC3339), result.HasDrift)
		diff.Render(os.Stdout, result.Entries, diff.RenderOptions{MaskValues: *mask})
	}
	return nil
}
