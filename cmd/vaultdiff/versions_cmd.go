package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// ListVersionsOptions holds configuration for the list-versions sub-command.
type ListVersionsOptions struct {
	Addr   string
	Token  string
	Mount  string
	Path   string
	Output io.Writer
}

// RunListVersions fetches and prints version metadata for a secret path.
func RunListVersions(opts ListVersionsOptions) error {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	client, err := vault.NewClient(vault.Config{
		Address: opts.Addr,
		Token:   opts.Token,
		Mount:   opts.Mount,
	})
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	metas, err := client.ListVersions(context.Background(), opts.Path)
	if err != nil {
		return fmt.Errorf("listing versions: %w", err)
	}

	sort.Slice(metas, func(i, j int) bool {
		return metas[i].Version < metas[j].Version
	})

	w := tabwriter.NewWriter(opts.Output, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "VERSION\tCREATED\tDELETED\tDESTROYED")
	for _, m := range metas {
		deletionTime := m.DeletionTime
		if deletionTime == "" {
			deletionTime = "-"
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%v\n",
			m.Version, m.CreatedTime, deletionTime, m.Destroyed)
	}
	return w.Flush()
}
