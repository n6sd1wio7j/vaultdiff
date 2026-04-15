package diff

import "sort"

// SecretData represents a map of key-value pairs from a Vault secret version.
type SecretData map[string]interface{}

// ChangeType describes the kind of change detected for a key.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// DiffEntry represents a single key-level difference between two secret versions.
type DiffEntry struct {
	Key      string
	Change   ChangeType
	OldValue interface{}
	NewValue interface{}
}

// Compare computes the diff between two SecretData maps and returns a sorted
// slice of DiffEntry values describing each key's change status.
func Compare(old, new SecretData) []DiffEntry {
	keys := unionKeys(old, new)
	sort.Strings(keys)

	entries := make([]DiffEntry, 0, len(keys))
	for _, k := range keys {
		oldVal, inOld := old[k]
		newVal, inNew := new[k]

		var change ChangeType
		switch {
		case inOld && !inNew:
			change = Removed
		case !inOld && inNew:
			change = Added
		case fmt(oldVal) != fmt(newVal):
			change = Modified
		default:
			change = Unchanged
		}

		entries = append(entries, DiffEntry{
			Key:      k,
			Change:   change,
			OldValue: oldVal,
			NewValue: newVal,
		})
	}
	return entries
}

// HasChanges returns true if any entry in the diff is not Unchanged.
func HasChanges(entries []DiffEntry) bool {
	for _, e := range entries {
		if e.Change != Unchanged {
			return true
		}
	}
	return false
}

func unionKeys(a, b SecretData) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}

func fmt(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	default:
		return string(rune(0)) // force non-equal for non-string types that differ
	}
}
