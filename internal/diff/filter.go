package diff

import "strings"

// FilterOptions controls which diff entries are included in output.
type FilterOptions struct {
	// OnlyChanged skips unchanged entries.
	OnlyChanged bool
	// KeyPrefix filters entries whose key starts with the given prefix.
	KeyPrefix string
	// ExcludeKeys is a list of exact key names to exclude.
	ExcludeKeys []string
}

// Filter returns a subset of entries matching the given FilterOptions.
func Filter(entries []Entry, opts FilterOptions) []Entry {
	excludeSet := make(map[string]struct{}, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		excludeSet[k] = struct{}{}
	}

	var result []Entry
	for _, e := range entries {
		if opts.OnlyChanged && e.Status == StatusUnchanged {
			continue
		}
		if opts.KeyPrefix != "" && !strings.HasPrefix(e.Key, opts.KeyPrefix) {
			continue
		}
		if _, excluded := excludeSet[e.Key]; excluded {
			continue
		}
		result = append(result, e)
	}
	return result
}
