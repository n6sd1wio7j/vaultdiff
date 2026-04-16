package diff

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// TimelineEntry represents a snapshot of drift score at a point in time.
type TimelineEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Path      string    `json:"path"`
	Score     DriftScore `json:"score"`
}

// Timeline holds an ordered list of drift observations.
type Timeline []TimelineEntry

// Append adds a new entry to the timeline.
func (t *Timeline) Append(path string, entries []Entry) TimelineEntry {
	e := TimelineEntry{
		Timestamp: time.Now().UTC(),
		Path:      path,
		Score:     ScoreDrift(entries),
	}
	*t = append(*t, e)
	return e
}

// Sorted returns entries ordered by timestamp ascending.
func (t Timeline) Sorted() Timeline {
	copy := append(Timeline(nil), t...)
	sort.Slice(copy, func(i, j int) bool {
		return copy[i].Timestamp.Before(copy[j].Timestamp)
	})
	return copy
}

// FormatTimeline renders a human-readable timeline table.
func FormatTimeline(t Timeline) string {
	if len(t) == 0 {
		return "no timeline entries\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s %-30s %6s %6s %6s %8s\n",
		"Timestamp", "Path", "Added", "Rmvd", "Mod", "Score"))
	sb.WriteString(strings.Repeat("-", 90) + "\n")
	for _, e := range t.Sorted() {
		sb.WriteString(fmt.Sprintf("%-30s %-30s %6d %6d %6d %8.2f\n",
			e.Timestamp.Format(time.RFC3339),
			e.Path,
			e.Score.Added,
			e.Score.Removed,
			e.Score.Modified,
			e.Score.WeightedScore,
		))
	}
	return sb.String()
}
