package diff

import (
	"fmt"
	"strings"
)

// DeprecateOptions controls which keys are flagged as deprecated.
type DeprecateOptions struct {
	// DeprecatedKeys is a list of exact key names considered deprecated.
	DeprecatedKeys []string
	// DeprecatedPrefixes is a list of key prefixes considered deprecated.
	DeprecatedPrefixes []string
	// IncludeUnchanged also flags unchanged entries that match deprecated rules.
	IncludeUnchanged bool
}

// DeprecatedEntry pairs a diff entry with a deprecation reason.
type DeprecatedEntry struct {
	Entry  Entry
	Reason string
}

// DefaultDeprecateOptions returns sensible defaults for deprecation detection.
func DefaultDeprecateOptions() DeprecateOptions {
	return DeprecateOptions{
		DeprecatedPrefixes: []string{"legacy_", "old_", "deprecated_"},
	}
}

// DetectDeprecated scans entries and returns those matching deprecation rules.
func DetectDeprecated(entries []Entry, opts DeprecateOptions) []DeprecatedEntry {
	var results []DeprecatedEntry
	for _, e := range entries {
		if !opts.IncludeUnchanged && e.Status == StatusUnchanged {
			continue
		}
		if reason, ok := matchDeprecated(e.Key, opts); ok {
			results = append(results, DeprecatedEntry{Entry: e, Reason: reason})
		}
	}
	return results
}

func matchDeprecated(key string, opts DeprecateOptions) (string, bool) {
	lower := strings.ToLower(key)
	for _, k := range opts.DeprecatedKeys {
		if strings.ToLower(k) == lower {
			return fmt.Sprintf("key %q is explicitly deprecated", key), true
		}
	}
	for _, p := range opts.DeprecatedPrefixes {
		if strings.HasPrefix(lower, strings.ToLower(p)) {
			return fmt.Sprintf("key %q matches deprecated prefix %q", key, p), true
		}
	}
	return "", false
}

// FormatDeprecated returns a human-readable report of deprecated entries.
func FormatDeprecated(entries []DeprecatedEntry) string {
	if len(entries) == 0 {
		return "no deprecated keys detected\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("deprecated keys detected: %d\n", len(entries)))
	for _, d := range entries {
		sb.WriteString(fmt.Sprintf("  [%s] %s — %s\n", d.Entry.Status, d.Entry.Key, d.Reason))
	}
	return sb.String()
}
