package diff

import "strings"

// IgnoreOptions configures which entries to ignore during comparison.
type IgnoreOptions struct {
	Keys    []string
	Prefixes []string
	Statuses []string
}

// IgnoreEntry returns true if the entry should be ignored based on the options.
func IgnoreEntry(e Entry, opts IgnoreOptions) bool {
	for _, k := range opts.Keys {
		if strings.EqualFold(e.Key, k) {
			return true
		}
	}
	for _, p := range opts.Prefixes {
		if strings.HasPrefix(e.Key, p) {
			return true
		}
	}
	for _, s := range opts.Statuses {
		if string(e.Status) == s {
			return true
		}
	}
	return false
}

// ApplyIgnore filters out ignored entries from a slice.
func ApplyIgnore(entries []Entry, opts IgnoreOptions) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if !IgnoreEntry(e, opts) {
			out = append(out, e)
		}
	}
	return out
}
