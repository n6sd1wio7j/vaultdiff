package diff

import (
	"fmt"
	"strings"
)

// PatchEntry represents a single line in a unified-style patch.
type PatchEntry struct {
	Op    string // "+", "-", or " "
	Key   string
	Value string
}

// Patch generates a unified-style patch from a slice of DiffEntries.
func Patch(entries []Entry) []PatchEntry {
	var patch []PatchEntry
	for _, e := range entries {
		switch e.Status {
		case StatusAdded:
			patch = append(patch, PatchEntry{Op: "+", Key: e.Key, Value: e.NewValue})
		case StatusRemoved:
			patch = append(patch, PatchEntry{Op: "-", Key: e.Key, Value: e.OldValue})
		case StatusModified:
			patch = append(patch, PatchEntry{Op: "-", Key: e.Key, Value: e.OldValue})
			patch = append(patch, PatchEntry{Op: "+", Key: e.Key, Value: e.NewValue})
		case StatusUnchanged:
			patch = append(patch, PatchEntry{Op: " ", Key: e.Key, Value: e.NewValue})
		}
	}
	return patch
}

// FormatPatch renders a patch slice as a human-readable string.
func FormatPatch(patch []PatchEntry, mask bool) string {
	var sb strings.Builder
	for _, p := range patch {
		val := p.Value
		if mask && p.Op != " " {
			val = "***"
		}
		sb.WriteString(fmt.Sprintf("%s %s=%s\n", p.Op, p.Key, val))
	}
	return sb.String()
}
