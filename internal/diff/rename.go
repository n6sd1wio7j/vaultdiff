package diff

import "fmt"

// RenameEntry represents a detected key rename between two secret versions.
type RenameEntry struct {
	OldKey string
	NewKey string
	Value  string
}

// RenameOptions controls how rename detection is performed.
type RenameOptions struct {
	// MinValueLength is the minimum value length to consider for rename matching.
	MinValueLength int
	// MaskValues hides the matched value in output.
	MaskValues bool
}

// DefaultRenameOptions returns sensible defaults for rename detection.
func DefaultRenameOptions() RenameOptions {
	return RenameOptions{
		MinValueLength: 1,
		MaskValues:     true,
	}
}

// DetectRenames scans diff entries for keys that were removed in one version
// and added in another with the same value, suggesting a rename occurred.
func DetectRenames(entries []Entry, opts RenameOptions) []RenameEntry {
	if opts.MinValueLength < 1 {
		opts.MinValueLength = 1
	}

	removed := make(map[string]string) // value -> key
	added := make(map[string]string)   // value -> key

	for _, e := range entries {
		switch e.Status {
		case StatusRemoved:
			if len(e.OldValue) >= opts.MinValueLength {
				removed[e.OldValue] = e.Key
			}
		case StatusAdded:
			if len(e.NewValue) >= opts.MinValueLength {
				added[e.NewValue] = e.Key
			}
		}
	}

	var renames []RenameEntry
	for val, oldKey := range removed {
		if newKey, ok := added[val]; ok {
			display := val
			if opts.MaskValues {
				display = "***"
			}
			renames = append(renames, RenameEntry{
				OldKey: oldKey,
				NewKey: newKey,
				Value:  display,
			})
		}
	}
	return renames
}

// FormatRenames returns a human-readable summary of detected renames.
func FormatRenames(renames []RenameEntry) string {
	if len(renames) == 0 {
		return "no renames detected\n"
	}
	out := fmt.Sprintf("detected %d rename(s):\n", len(renames))
	for _, r := range renames {
		out += fmt.Sprintf("  %s -> %s (value: %s)\n", r.OldKey, r.NewKey, r.Value)
	}
	return out
}
