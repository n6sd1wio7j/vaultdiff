package diff

import (
	"strings"
	"testing"
)

func sampleHeatmapSteps() []ChainStep {
	step1 := ChainStep{
		Entries: []Entry{
			{Key: "DB_PASSWORD", Status: StatusModified},
			{Key: "API_KEY", Status: StatusModified},
			{Key: "LOG_LEVEL", Status: StatusUnchanged},
		},
	}
	step2 := ChainStep{
		Entries: []Entry{
			{Key: "DB_PASSWORD", Status: StatusModified},
			{Key: "API_KEY", Status: StatusUnchanged},
			{Key: "LOG_LEVEL", Status: StatusAdded},
		},
	}
	step3 := ChainStep{
		Entries: []Entry{
			{Key: "DB_PASSWORD", Status: StatusUnchanged},
			{Key: "API_KEY", Status: StatusRemoved},
			{Key: "LOG_LEVEL", Status: StatusUnchanged},
		},
	}
	return []ChainStep{step1, step2, step3}
}

func TestBuildHeatmap_CountsChanges(t *testing.T) {
	steps := sampleHeatmapSteps()
	entries := BuildHeatmap(steps, DefaultHeatmapOptions())
	freq := map[string]int{}
	for _, e := range entries {
		freq[e.Key] = e.Changes
	}
	if freq["DB_PASSWORD"] != 2 {
		t.Errorf("expected DB_PASSWORD changes=2, got %d", freq["DB_PASSWORD"])
	}
	if freq["API_KEY"] != 2 {
		t.Errorf("expected API_KEY changes=2, got %d", freq["API_KEY"])
	}
	if freq["LOG_LEVEL"] != 1 {
		t.Errorf("expected LOG_LEVEL changes=1, got %d", freq["LOG_LEVEL"])
	}
}

func TestBuildHeatmap_ExcludesUnchangedByDefault(t *testing.T) {
	steps := []ChainStep{
		{Entries: []Entry{{Key: "STATIC", Status: StatusUnchanged}}},
	}
	entries := BuildHeatmap(steps, DefaultHeatmapOptions())
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestBuildHeatmap_ShowAllIncludesUnchanged(t *testing.T) {
	steps := []ChainStep{
		{Entries: []Entry{{Key: "STATIC", Status: StatusUnchanged}}},
	}
	opts := DefaultHeatmapOptions()
	opts.ShowAll = true
	entries := BuildHeatmap(steps, opts)
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
}

func TestBuildHeatmap_TopNLimitsResults(t *testing.T) {
	steps := sampleHeatmapSteps()
	opts := DefaultHeatmapOptions()
	opts.TopN = 2
	entries := BuildHeatmap(steps, opts)
	if len(entries) > 2 {
		t.Errorf("expected at most 2 entries, got %d", len(entries))
	}
}

func TestBuildHeatmap_FrequencyCalculation(t *testing.T) {
	steps := sampleHeatmapSteps() // 3 steps
	entries := BuildHeatmap(steps, DefaultHeatmapOptions())
	for _, e := range entries {
		if e.Key == "DB_PASSWORD" {
			expected := 2.0 / 3.0
			if e.Frequency < expected-0.001 || e.Frequency > expected+0.001 {
				t.Errorf("expected frequency ~%.3f, got %.3f", expected, e.Frequency)
			}
		}
	}
}

func TestFormatHeatmap_Empty(t *testing.T) {
	out := FormatHeatmap(nil)
	if !strings.Contains(out, "no heatmap") {
		t.Errorf("expected no-data message, got: %s", out)
	}
}

func TestFormatHeatmap_ContainsKeys(t *testing.T) {
	steps := sampleHeatmapSteps()
	entries := BuildHeatmap(steps, DefaultHeatmapOptions())
	out := FormatHeatmap(entries)
	if !strings.Contains(out, "DB_PASSWORD") {
		t.Errorf("expected DB_PASSWORD in output")
	}
	if !strings.Contains(out, "KEY") {
		t.Errorf("expected header in output")
	}
}
