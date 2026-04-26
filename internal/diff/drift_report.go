package diff

import (
	"fmt"
	"strings"
	"time"
)

// DriftReportOptions configures the drift report output.
type DriftReportOptions struct {
	Env         string
	Path        string
	VersionA    int
	VersionB    int
	IncludeScore bool
	IncludeSummary bool
}

// DriftReport is a structured report combining summary, score, and entries.
type DriftReport struct {
	GeneratedAt time.Time
	Env         string
	Path        string
	VersionA    int
	VersionB    int
	Summary     Summary
	Score       DriftScore
	Entries     []Entry
}

// BuildDriftReport assembles a full drift report from a set of diff entries.
func BuildDriftReport(entries []Entry, opts DriftReportOptions) DriftReport {
	return DriftReport{
		GeneratedAt: time.Now().UTC(),
		Env:         opts.Env,
		Path:        opts.Path,
		VersionA:    opts.VersionA,
		VersionB:    opts.VersionB,
		Summary:     Summarize(entries),
		Score:       ScoreDrift(entries),
		Entries:     entries,
	}
}

// FormatDriftReport renders the drift report as a human-readable string.
func FormatDriftReport(r DriftReport, opts DriftReportOptions) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Drift Report — %s\n", r.GeneratedAt.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("Path:     %s\n", r.Path))
	if r.Env != "" {
		sb.WriteString(fmt.Sprintf("Env:      %s\n", r.Env))
	}
	sb.WriteString(fmt.Sprintf("Versions: v%d → v%d\n", r.VersionA, r.VersionB))
	sb.WriteString(strings.Repeat("-", 48) + "\n")

	if opts.IncludeSummary {
		sb.WriteString(fmt.Sprintf("Added: %d  Removed: %d  Modified: %d  Unchanged: %d\n",
			r.Summary.Added, r.Summary.Removed, r.Summary.Modified, r.Summary.Unchanged))
	}

	if opts.IncludeScore {
		sb.WriteString(FormatScore(r.Score) + "\n")
	}

	sb.WriteString(strings.Repeat("-", 48) + "\n")
	for _, e := range r.Entries {
		sb.WriteString(formatLine(e, false) + "\n")
	}

	return sb.String()
}
