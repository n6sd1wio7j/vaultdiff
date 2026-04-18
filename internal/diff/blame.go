package diff

import (
	"fmt"
	"strings"
	"time"
)

// BlameEntry records who last changed a key and when.
type BlameEntry struct {
	Key       string    `json:"key"`
	ChangedAt time.Time `json:"changed_at"`
	Version   int       `json:"version"`
	ChangeType string   `json:"change_type"`
}

// BlameReport maps keys to their blame entries.
type BlameReport struct {
	Path    string                 `json:"path"`
	Entries map[string]BlameEntry  `json:"entries"`
}

// Blame builds a BlameReport from a slice of DiffEntries and version metadata.
func Blame(path string, entries []Entry, version int, at time.Time) BlameReport {
	report := BlameReport{
		Path:    path,
		Entries: make(map[string]BlameEntry),
	}
	for _, e := range entries {
		if e.Change == ChangeUnchanged {
			continue
		}
		report.Entries[e.Key] = BlameEntry{
			Key:        e.Key,
			ChangedAt:  at,
			Version:    version,
			ChangeType: string(e.Change),
		}
	}
	return report
}

// FormatBlame returns a human-readable blame report.
func FormatBlame(r BlameReport) string {
	if len(r.Entries) == 0 {
		return fmt.Sprintf("blame: no changes recorded for %s\n", r.Path)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("blame report: %s\n", r.Path))
	sb.WriteString(strings.Repeat("-", 48) + "\n")
	for _, e := range r.Entries {
		sb.WriteString(fmt.Sprintf("  %-24s v%-4d %-8s %s\n",
			e.Key, e.Version, e.ChangeType, e.ChangedAt.Format(time.RFC3339)))
	}
	return sb.String()
}
