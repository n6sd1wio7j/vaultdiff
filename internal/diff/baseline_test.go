package diff

import (
	"os"
	"path/filepath"
	"testing"
)

func sampleBaselineEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Status: StatusUnchanged, OldValue: "localhost", NewValue: "localhost"},
		{Key: "API_KEY", Status: StatusUnchanged, OldValue: "secret", NewValue: "secret"},
	}
}

func TestSaveAndLoadBaseline_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "baseline.json")
	entries := sampleBaselineEntries()

	if err := SaveBaseline(file, "secret/myapp", entries); err != nil {
		t.Fatalf("SaveBaseline error: %v", err)
	}

	b, err := LoadBaseline(file)
	if err != nil {
		t.Fatalf("LoadBaseline error: %v", err)
	}
	if b.Path != "secret/myapp" {
		t.Errorf("expected path secret/myapp, got %s", b.Path)
	}
	if len(b.Entries) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(b.Entries))
	}
	if b.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestLoadBaseline_MissingFile(t *testing.T) {
	_, err := LoadBaseline("/nonexistent/baseline.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadBaseline_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "bad.json")
	os.WriteFile(file, []byte("not json"), 0644)
	_, err := LoadBaseline(file)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestCompareToBaseline_DetectsChanges(t *testing.T) {
	baseline := &Baseline{
		Path: "secret/myapp",
		Entries: []Entry{
			{Key: "DB_HOST", NewValue: "localhost"},
			{Key: "API_KEY", NewValue: "old-secret"},
		},
	}
	current := map[string]string{
		"DB_HOST": "localhost",
		"API_KEY": "new-secret",
	}
	entries := CompareToBaseline(baseline, current)
	changed := 0
	for _, e := range entries {
		if e.Status == StatusModified {
			changed++
		}
	}
	if changed != 1 {
		t.Errorf("expected 1 modified entry, got %d", changed)
	}
}
