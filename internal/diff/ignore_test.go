package diff

import (
	"testing"
)

func sampleIgnoreEntries() []Entry {
	return []Entry{
		{Key: "API_KEY", Status: StatusAdded},
		{Key: "DB_PASSWORD", Status: StatusModified},
		{Key: "LOG_LEVEL", Status: StatusUnchanged},
		{Key: "secret_token", Status: StatusRemoved},
		{Key: "feature_flag", Status: StatusUnchanged},
	}
}

func TestApplyIgnore_ExactKey(t *testing.T) {
	entries := sampleIgnoreEntries()
	result := ApplyIgnore(entries, IgnoreOptions{Keys: []string{"API_KEY"}})
	for _, e := range result {
		if e.Key == "API_KEY" {
			t.Fatal("expected API_KEY to be ignored")
		}
	}
}

func TestApplyIgnore_Prefix(t *testing.T) {
	entries := sampleIgnoreEntries()
	result := ApplyIgnore(entries, IgnoreOptions{Prefixes: []string{"secret_"}})
	for _, e := range result {
		if e.Key == "secret_token" {
			t.Fatal("expected secret_token to be ignored")
		}
	}
}

func TestApplyIgnore_Status(t *testing.T) {
	entries := sampleIgnoreEntries()
	result := ApplyIgnore(entries, IgnoreOptions{Statuses: []string{"unchanged"}})
	for _, e := range result {
		if e.Status == StatusUnchanged {
			t.Fatalf("expected unchanged entries to be ignored, got %s", e.Key)
		}
	}
}

func TestApplyIgnore_EmptyOptions(t *testing.T) {
	entries := sampleIgnoreEntries()
	result := ApplyIgnore(entries, IgnoreOptions{})
	if len(result) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(result))
	}
}

func TestApplyIgnore_CaseInsensitiveKey(t *testing.T) {
	entries := sampleIgnoreEntries()
	result := ApplyIgnore(entries, IgnoreOptions{Keys: []string{"api_key"}})
	for _, e := range result {
		if e.Key == "API_KEY" {
			t.Fatal("expected case-insensitive match to ignore API_KEY")
		}
	}
}
