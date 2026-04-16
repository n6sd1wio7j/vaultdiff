package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// RunAnnotate diffs two secret versions and renders output with inline annotations.
func RunAnnotate(args []string) error {
	fs := flag.NewFlagSet("annotate", flag.ContinueOnError)
	path := fs.String("path", "", "secret path")
	addr := fs.String("addr", "", "vault address")
	token := fs.String("token", "", "vault token")
	v1 := fs.Int("v1", 0, "first version")
	v2 := fs.Int("v2", 0, "second version")
	annotFile := fs.String("annotations", "", "JSON file with annotations [{key,note}]")
	mask := fs.Bool("mask", true, "mask secret values")

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

	client, err := vault.NewClient(*addr, *token, "")
	if err != nil {
		return err
	}

	secA, err := client.ReadSecretVersion(*path, *v1)
	if err != nil {
		return fmt.Errorf("reading v%d: %w", *v1, err)
	}
	secB, err := client.ReadSecretVersion(*path, *v2)
	if err != nil {
		return fmt.Errorf("reading v%d: %w", *v2, err)
	}

	entries := diff.Compare(secA, secB)

	var opts diff.AnnotateOptions
	if *annotFile != "" {
		data, err := os.ReadFile(*annotFile)
		if err != nil {
			return fmt.Errorf("reading annotations file: %w", err)
		}
		if err := json.Unmarshal(data, &opts.Annotations); err != nil {
			return fmt.Errorf("parsing annotations: %w", err)
		}
	}

	notes := diff.Annotate(entries, opts)
	for _, e := range entries {
		fmt.Println(diff.FormatAnnotated(e, notes[e.Key], *mask))
	}
	return nil
}
