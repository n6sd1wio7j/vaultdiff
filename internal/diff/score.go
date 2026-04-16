package diff

import "fmt"

// DriftScore represents a numeric risk score for a set of diff entries.
type DriftScore struct {
	Total    int
	Added    int
	Removed  int
	Modified int
	Score    float64
}

// weights assigned to each change type
const (
	weightAdded    = 1.0
	weightRemoved  = 2.0
	weightModified = 1.5
)

// ScoreDrift computes a weighted drift score from a slice of DiffEntry.
// Higher scores indicate more significant secret drift.
func ScoreDrift(entries []Entry) DriftScore {
	var added, removed, modified int
	for _, e := range entries {
		switch e.Status {
		case StatusAdded:
			added++
		case StatusRemoved:
			removed++
		case StatusModified:
			modified++
		}
	}
	score := float64(added)*weightAdded +
		float64(removed)*weightRemoved +
		float64(modified)*weightModified
	return DriftScore{
		Total:    len(entries),
		Added:    added,
		Removed:  removed,
		Modified: modified,
		Score:    score,
	}
}

// FormatScore returns a human-readable summary of the drift score.
func FormatScore(ds DriftScore) string {
	return fmt.Sprintf(
		"drift score: %.1f (added=%d removed=%d modified=%d total=%d)",
		ds.Score, ds.Added, ds.Removed, ds.Modified, ds.Total,
	)
}
