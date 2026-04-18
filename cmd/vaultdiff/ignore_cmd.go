package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// RunIgnoreDiff runs a diff with ignore rules applied.
func RunIgnoreDiff(args []string) error {
	fs := flag.NewFlagSet("ignore", flag.ContinueOnError)
	path := fs.String("path", "", "secret path")
	addr := fs.String("addr", "", "vault address")
	token := fs.String("token", "", "vault token")
	v1 := fs.Int("v1", 1, "version 1")
	v2 := fs.Int("v2", 2, "version 2")
	ignoreKeys := fs.String("ignore-keys", "", "comma-separated keys to ignore")
	ignorePrefixes := fs.String("ignore-prefixes", "", "comma-separated key prefixes to ignore")
	ignoreStatuses := fs.String("ignore-statuses", "", "comma-separated statuses to ignore")

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

	c, err := vault.NewClient(*addr, *token, "")
	if err != nil {
		return err
	}

	secA, err := c.ReadSecretVersion(*path, *v1)
	if err != nil {
		return err
	}
	secB, err := c.ReadSecretVersion(*path, *v2)
	if err != nil {
		return err
	}

	entries := diff.Compare(secA, secB)

	opts := diff.IgnoreOptions{}
	if *ignoreKeys != "" {
		opts.Keys = strings.Split(*ignoreKeys, ",")
	}
	if *ignorePrefixes != "" {
		opts.Prefixes = strings.Split(*ignorePrefixes, ",")
	}
	if *ignoreStatuses != "" {
		opts.Statuses = strings.Split(*ignoreStatuses, ",")
	}

	filtered := diff.ApplyIgnore(entries, opts)
	fmt.Fprint(os.Stdout, diff.Render(filtered, diff.RenderOptions{}))
	return nil
}
