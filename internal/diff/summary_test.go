package diff

import (
	"strings"
	"testing"
)

func sampleSummaryEntries() []Entry {
	return []Entry{
		{Key: "A", Status: StatusAdded},
		{Key: "B", Status: StatusAdded},
		{Key: "C", Status: StatusRemoved},
		{Key: "D", Status: StatusModified},
		{Key: "E", Status: StatusUnchanged},
		{Key: "F", Status: StatusUnchanged},
	}
}

func TestSummarize_Counts(t *testing.T) {
	s := Summarize(sampleSummaryEntries())
	if s.Added != 2 {
		t.Errorf("expected 2 added, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", s.Removed)
	}
	if s.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", s.Modified)
	}
	if s.Unchanged != 2 {
		t.Errorf("expected 2 unchanged, got %d", s.Unchanged)
	}
	if s.Total != 6 {
		t.Errorf("expected total 6, got %d", s.Total)
	}
}

func TestSummarize_Empty(t *testing.T) {
	s := Summarize([]Entry{})
	if s.Total != 0 || s.HasDrift() {
		t.Error("expected zero summary with no drift")
	}
}

func TestSummary_HasDrift_True(t *testing.T) {
	s := Summarize(sampleSummaryEntries())
	if !s.HasDrift() {
		t.Error("expected HasDrift to be true")
	}
}

func TestSummary_HasDrift_False(t *testing.T) {
	entries := []Entry{
		{Key: "X", Status: StatusUnchanged},
	}
	s := Summarize(entries)
	if s.HasDrift() {
		t.Error("expected HasDrift to be false when only unchanged entries")
	}
}

func TestSummary_String_ContainsFields(t *testing.T) {
	s := Summarize(sampleSummaryEntries())
	str := s.String()
	for _, want := range []string{"total=", "added=", "removed=", "modified=", "unchanged="} {
		if !strings.Contains(str, want) {
			t.Errorf("summary string missing field %q: %s", want, str)
		}
	}
}
