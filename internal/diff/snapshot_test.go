package diff

import (
	"os"
	"path/filepath"
	"testing"
)

var sampleSnapshotEntries = []Entry{
	{Key: "API_KEY", OldValue: "abc", NewValue: "abc", Status: StatusUnchanged},
	{Key: "DB_PASS", OldValue: "", NewValue: "secret", Status: StatusAdded},
}

func TestCaptureSnapshot_SetsFields(t *testing.T) {
	s := CaptureSnapshot("secret/app", 3, sampleSnapshotEntries)
	if s.Path != "secret/app" {
		t.Errorf("expected path secret/app, got %s", s.Path)
	}
	if s.Version != 3 {
		t.Errorf("expected version 3, got %d", s.Version)
	}
	if len(s.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(s.Entries))
	}
	if s.CapturedAt.IsZero() {
		t.Error("expected non-zero CapturedAt")
	}
}

func TestSaveAndLoadSnapshot_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "snap.json")

	orig := CaptureSnapshot("secret/svc", 1, sampleSnapshotEntries)
	if err := SaveSnapshot(file, orig); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	loaded, err := LoadSnapshot(file)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}
	if loaded.Path != orig.Path {
		t.Errorf("path mismatch: got %s", loaded.Path)
	}
	if loaded.Version != orig.Version {
		t.Errorf("version mismatch: got %d", loaded.Version)
	}
	if len(loaded.Entries) != len(orig.Entries) {
		t.Errorf("entries mismatch: got %d", len(loaded.Entries))
	}
}

func TestLoadSnapshot_MissingFile(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadSnapshot_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "bad.json")
	os.WriteFile(file, []byte("not-json"), 0644)
	_, err := LoadSnapshot(file)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestDiffSnapshot_DetectsChanges(t *testing.T) {
	old := CaptureSnapshot("secret/app", 1, []Entry{
		{Key: "HOST", OldValue: "localhost", NewValue: "localhost", Status: StatusUnchanged},
	})
	new := CaptureSnapshot("secret/app", 2, []Entry{
		{Key: "HOST", OldValue: "localhost", NewValue: "prod.host", Status: StatusModified},
	})
	entries := DiffSnapshot(old, new)
	if len(entries) == 0 {
		t.Fatal("expected diff entries")
	}
}
