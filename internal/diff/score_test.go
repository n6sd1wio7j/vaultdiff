package diff

import (
	"strings"
	"testing"
)

var sampleScoreEntries = []Entry{
	{Key: "DB_HOST", OldValue: "", NewValue: "localhost", Status: StatusAdded},
	{Key: "DB_PASS", OldValue: "secret", NewValue: "", Status: StatusRemoved},
	{Key: "API_KEY", OldValue: "old", NewValue: "new", Status: StatusModified},
	{Key: "PORT", OldValue: "8080", NewValue: "8080", Status: StatusUnchanged},
}

func TestScoreDrift_Counts(t *testing.T) {
	ds := ScoreDrift(sampleScoreEntries)
	if ds.Added != 1 {
		t.Errorf("expected 1 added, got %d", ds.Added)
	}
	if ds.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", ds.Removed)
	}
	if ds.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", ds.Modified)
	}
	if ds.Total != 4 {
		t.Errorf("expected total 4, got %d", ds.Total)
	}
}

func TestScoreDrift_WeightedScore(t *testing.T) {
	ds := ScoreDrift(sampleScoreEntries)
	expected := 1.0 + 2.0 + 1.5 // added + removed + modified
	if ds.Score != expected {
		t.Errorf("expected score %.1f, got %.1f", expected, ds.Score)
	}
}

func TestScoreDrift_Empty(t *testing.T) {
	ds := ScoreDrift([]Entry{})
	if ds.Score != 0 {
		t.Errorf("expected 0 score for empty entries, got %.1f", ds.Score)
	}
}

func TestScoreDrift_OnlyUnchanged(t *testing.T) {
	entries := []Entry{
		{Key: "A", Status: StatusUnchanged},
		{Key: "B", Status: StatusUnchanged},
	}
	ds := ScoreDrift(entries)
	if ds.Score != 0 {
		t.Errorf("expected 0 score, got %.1f", ds.Score)
	}
}

func TestFormatScore_ContainsFields(t *testing.T) {
	ds := ScoreDrift(sampleScoreEntries)
	out := FormatScore(ds)
	for _, substr := range []string{"drift score", "added=", "removed=", "modified=", "total="} {
		if !strings.Contains(out, substr) {
			t.Errorf("expected %q in output: %s", substr, out)
		}
	}
}
