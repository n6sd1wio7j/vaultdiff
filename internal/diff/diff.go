package diff

import "sort"

// Status represents the change state of a single secret key.
type Status string

const (
	StatusAdded     Status = "added"
	StatusRemoved   Status = "removed"
	StatusModified  Status = "modified"
	StatusUnchanged Status = "unchanged"
)

// Entry holds the diff result for a single key.
type Entry struct {
	Key      string
	Status   Status
	OldValue string
	NewValue string
}

// Compare produces a sorted slice of Entry values by comparing two secret maps.
func Compare(before, after map[string]string) []Entry {
	keys := unionKeys(before, after)
	sort.Strings(keys)

	entries := make([]Entry, 0, len(keys))
	for _, k := range keys {
		oldVal, inBefore := before[k]
		newVal, inAfter := after[k]

		var status Status
		switch {
		case inBefore && !inAfter:
			status = StatusRemoved
		case !inBefore && inAfter:
			status = StatusAdded
		case oldVal != newVal:
			status = StatusModified
		default:
			status = StatusUnchanged
		}

		entries = append(entries, Entry{
			Key:      k,
			Status:   status,
			OldValue: oldVal,
			NewValue: newVal,
		})
	}
	return entries
}

// HasChanges returns true if any entry is not StatusUnchanged.
func HasChanges(entries []Entry) bool {
	for _, e := range entries {
		if e.Status != StatusUnchanged {
			return true
		}
	}
	return false
}

func unionKeys(a, b map[string]string) []string {
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
