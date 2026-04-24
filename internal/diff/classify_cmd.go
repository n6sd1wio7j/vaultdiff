package diff

import (
	"fmt"
	"strings"
)

// ClassifyOptions controls how entries are classified.
type ClassifyOptions struct {
	ShowAll     bool
	MaskSecrets bool
}

// DefaultClassifyOptions returns sensible defaults for classification.
func DefaultClassifyOptions() ClassifyOptions {
	return ClassifyOptions{
		ShowAll:     false,
		MaskSecrets: true,
	}
}

// FormatClassifiedSummary returns a compact summary of classified entries.
func FormatClassifiedSummary(entries []ClassifiedEntry, opts ClassifyOptions) string {
	if len(entries) == 0 {
		return "no entries to classify\n"
	}

	counts := make(map[string]int)
	for _, e := range entries {
		counts[string(e.Category)]++
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("classified %d entries:\n", len(entries)))

	for _, e := range entries {
		if !opts.ShowAll && e.Entry.Status == StatusUnchanged {
			continue
		}
		val := e.Entry.NewValue
		if opts.MaskSecrets && e.Sensitive {
			val = "***"
		}
		sb.WriteString(fmt.Sprintf("  [%s] %s (%s) = %s\n",
			e.Category, e.Entry.Key, e.Entry.Status, val))
	}

	sb.WriteString("summary:\n")
	for cat, count := range counts {
		sb.WriteString(fmt.Sprintf("  %s: %d\n", cat, count))
	}

	return sb.String()
}
