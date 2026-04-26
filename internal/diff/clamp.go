package diff

import (
	"fmt"
	"strings"
)

// ClampOptions controls which entries are included based on value length constraints.
type ClampOptions struct {
	MinLength int // entries with values shorter than this are excluded (0 = no min)
	MaxLength int // entries with values longer than this are excluded (0 = no max)
	MaskValues bool
}

// DefaultClampOptions returns sensible defaults for clamping.
func DefaultClampOptions() ClampOptions {
	return ClampOptions{
		MinLength:  0,
		MaxLength:  0,
		MaskValues: true,
	}
}

// ClampResult holds an entry that passed the clamp filter along with metadata.
type ClampResult struct {
	Entry  DiffEntry
	Length int
	Value  string
}

// Clamp filters diff entries whose old or new values fall outside the given length bounds.
// Unchanged entries are always included unless they violate the length constraint.
func Clamp(entries []DiffEntry, opts ClampOptions) []ClampResult {
	var results []ClampResult
	for _, e := range entries {
		val := e.NewValue
		if val == "" {
			val = e.OldValue
		}
		l := len(val)
		if opts.MinLength > 0 && l < opts.MinLength {
			continue
		}
		if opts.MaxLength > 0 && l > opts.MaxLength {
			continue
		}
		display := val
		if opts.MaskValues {
			display = "***"
		}
		results = append(results, ClampResult{
			Entry:  e,
			Length: l,
			Value:  display,
		})
	}
	return results
}

// FormatClamp returns a human-readable table of clamped results.
func FormatClamp(results []ClampResult) string {
	if len(results) == 0 {
		return "no entries within clamp bounds\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-40s %-10s %6s\n", "KEY", "STATUS", "LEN"))
	sb.WriteString(strings.Repeat("-", 60) + "\n")
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("%-40s %-10s %6d\n", r.Entry.Key, string(r.Entry.Status), r.Length))
	}
	return sb.String()
}
