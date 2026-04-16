package diff

import "strings"

// RedactOptions controls which keys are redacted in output.
type RedactOptions struct {
	// KeySuffixes defines key suffixes that trigger redaction (e.g. "_key", "_secret").
	KeySuffixes []string
	// ExactKeys defines exact key names to always redact.
	ExactKeys []string
}

// DefaultRedactOptions returns sensible defaults for secret redaction.
func DefaultRedactOptions() RedactOptions {
	return RedactOptions{
		KeySuffixes: []string{"_key", "_secret", "_token", "_password", "_pass", "_pwd"},
		ExactKeys:   []string{"password", "secret", "token", "apikey", "api_key"},
	}
}

// ShouldRedact returns true if the given key matches redaction rules.
func ShouldRedact(key string, opts RedactOptions) bool {
	lower := strings.ToLower(key)
	for _, exact := range opts.ExactKeys {
		if lower == exact {
			return true
		}
	}
	for _, suffix := range opts.KeySuffixes {
		if strings.HasSuffix(lower, suffix) {
			return true
		}
	}
	return false
}

// Redact applies redaction to a slice of DiffEntry values, replacing
// sensitive values with "***".
func Redact(entries []DiffEntry, opts RedactOptions) []DiffEntry {
	out := make([]DiffEntry, len(entries))
	for i, e := range entries {
		if ShouldRedact(e.Key, opts) {
			if e.OldValue != "" {
				e.OldValue = "***"
			}
			if e.NewValue != "" {
				e.NewValue = "***"
			}
		}
		out[i] = e
	}
	return out
}
