package diff

import "fmt"

// Summary holds aggregated statistics about a diff result.
type Summary struct {
	Added     int
	Removed   int
	Modified  int
	Unchanged int
	Total     int
}

// Summarize computes statistics from a slice of DiffEntry values.
func Summarize(entries []Entry) Summary {
	s := Summary{}
	for _, e := range entries {
		s.Total++
		switch e.Status {
		case StatusAdded:
			s.Added++
		case StatusRemoved:
			s.Removed++
		case StatusModified:
			s.Modified++
		case StatusUnchanged:
			s.Unchanged++
		}
	}
	return s
}

// String returns a human-readable one-line summary.
func (s Summary) String() string {
	return fmt.Sprintf(
		"total=%d added=%d removed=%d modified=%d unchanged=%d",
		s.Total, s.Added, s.Removed, s.Modified, s.Unchanged,
	)
}

// HasDrift returns true when any keys were added, removed, or modified.
func (s Summary) HasDrift() bool {
	return s.Added > 0 || s.Removed > 0 || s.Modified > 0
}
