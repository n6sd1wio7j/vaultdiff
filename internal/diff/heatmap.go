package diff

import (
	"fmt"
	"sort"
	"strings"
)

// HeatmapEntry represents a key and its change frequency across versions.
type HeatmapEntry struct {
	Key       string
	Changes   int
	Total     int
	Frequency float64 // Changes / Total
}

// HeatmapOptions controls heatmap generation.
type HeatmapOptions struct {
	TopN    int  // 0 means all
	ShowAll bool // include keys with zero changes
}

// DefaultHeatmapOptions returns sensible defaults.
func DefaultHeatmapOptions() HeatmapOptions {
	return HeatmapOptions{
		TopN:    10,
		ShowAll: false,
	}
}

// BuildHeatmap analyses a slice of ChainStep diffs and computes change
// frequency per key across all steps.
func BuildHeatmap(steps []ChainStep, opts HeatmapOptions) []HeatmapEntry {
	counts := make(map[string]int)
	total := len(steps)

	for _, step := range steps {
		for _, e := range step.Entries {
			if e.Status != StatusUnchanged {
				counts[e.Key]++
			}
		}
	}

	// collect all keys seen across steps
	allKeys := make(map[string]struct{})
	for _, step := range steps {
		for _, e := range step.Entries {
			allKeys[e.Key] = struct{}{}
		}
	}

	var entries []HeatmapEntry
	for key := range allKeys {
		changes := counts[key]
		if !opts.ShowAll && changes == 0 {
			continue
		}
		freq := 0.0
		if total > 0 {
			freq = float64(changes) / float64(total)
		}
		entries = append(entries, HeatmapEntry{
			Key:       key,
			Changes:   changes,
			Total:     total,
			Frequency: freq,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Changes != entries[j].Changes {
			return entries[i].Changes > entries[j].Changes
		}
		return entries[i].Key < entries[j].Key
	})

	if opts.TopN > 0 && len(entries) > opts.TopN {
		entries = entries[:opts.TopN]
	}
	return entries
}

// FormatHeatmap returns a human-readable heatmap table.
func FormatHeatmap(entries []HeatmapEntry) string {
	if len(entries) == 0 {
		return "no heatmap data\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-40s %8s %8s %8s\n", "KEY", "CHANGES", "TOTAL", "FREQ"))
	sb.WriteString(strings.Repeat("-", 68) + "\n")
	for _, e := range entries {
		bar := heatBar(e.Frequency)
		sb.WriteString(fmt.Sprintf("%-40s %8d %8d %7.0f%% %s\n",
			e.Key, e.Changes, e.Total, e.Frequency*100, bar))
	}
	return sb.String()
}

func heatBar(freq float64) string {
	n := int(freq * 10)
	if n > 10 {
		n = 10
	}
	return strings.Repeat("█", n) + strings.Repeat("░", 10-n)
}
