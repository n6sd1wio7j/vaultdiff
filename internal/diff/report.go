package diff

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// ReportOptions controls report generation behavior.
type ReportOptions struct {
	SourceEnv  string
	TargetEnv  string
	SecretPath string
	ShowMasked bool
	Timestamp  time.Time
}

// Report writes a structured audit report of the diff to the given writer.
func Report(w io.Writer, entries []Entry, opts ReportOptions) error {
	ts := opts.Timestamp
	if ts.IsZero() {
		ts = time.Now().UTC()
	}

	header := buildHeader(opts, ts)
	if _, err := fmt.Fprintln(w, header); err != nil {
		return err
	}

	if !HasChanges(entries) {
		_, err := fmt.Fprintln(w, "No changes detected.")
		return err
	}

	var added, removed, modified int
	for _, e := range entries {
		switch e.Status {
		case StatusAdded:
			added++
		case StatusRemoved:
			removed++
		case StatusModified:
			modified++
		}
	}

	_, err := fmt.Fprintf(w, "Summary: +%d added, -%d removed, ~%d modified\n\n", added, removed, modified)
	if err != nil {
		return err
	}

	renderOpts := RenderOptions{
		ShowUnchanged: false,
		MaskValues:    !opts.ShowMasked,
	}
	lines := Render(entries, renderOpts)
	for _, line := range lines {
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func buildHeader(opts ReportOptions, ts time.Time) string {
	var sb strings.Builder
	sb.WriteString(strings.Repeat("-", 60) + "\n")
	sb.WriteString(fmt.Sprintf("VaultDiff Audit Report — %s\n", ts.Format(time.RFC3339)))
	if opts.SecretPath != "" {
		sb.WriteString(fmt.Sprintf("Path:   %s\n", opts.SecretPath))
	}
	if opts.SourceEnv != "" || opts.TargetEnv != "" {
		sb.WriteString(fmt.Sprintf("Source: %s  →  Target: %s\n", opts.SourceEnv, opts.TargetEnv))
	}
	sb.WriteString(strings.Repeat("-", 60))
	return sb.String()
}
