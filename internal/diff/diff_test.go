package diff

import (
	"testing"
)

func TestCompare_Added(t *testing.T) {
	old := SecretData{}
	new := SecretData{"key": "value"}
	entries := Compare(old, new)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Change != Added {
		t.Errorf("expected Added, got %s", entries[0].Change)
	}
	if entries[0].Key != "key" {
		t.Errorf("expected key 'key', got '%s'", entries[0].Key)
	}
}

func TestCompare_Removed(t *testing.T) {
	old := SecretData{"key": "value"}
	new := SecretData{}
	entries := Compare(old, new)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Change != Removed {
		t.Errorf("expected Removed, got %s", entries[0].Change)
	}
}

func TestCompare_Modified(t *testing.T) {
	old := SecretData{"db_pass": "old_secret"}
	new := SecretData{"db_pass": "new_secret"}
	entries := Compare(old, new)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Change != Modified {
		t.Errorf("expected Modified, got %s", entries[0].Change)
	}
}

func TestCompare_Unchanged(t *testing.T) {
	old := SecretData{"api_key": "abc123"}
	new := SecretData{"api_key": "abc123"}
	entries := Compare(old, new)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Change != Unchanged {
		t.Errorf("expected Unchanged, got %s", entries[0].Change)
	}
}

func TestCompare_SortedKeys(t *testing.T) {
	old := SecretData{"z": "1", "a": "2", "m": "3"}
	new := SecretData{"z": "1", "a": "2", "m": "3"}
	entries := Compare(old, new)
	expected := []string{"a", "m", "z"}
	for i, e := range entries {
		if e.Key != expected[i] {
			t.Errorf("index %d: expected key '%s', got '%s'", i, expected[i], e.Key)
		}
	}
}

func TestHasChanges_True(t *testing.T) {
	entries := []DiffEntry{{Key: "x", Change: Added}}
	if !HasChanges(entries) {
		t.Error("expected HasChanges to return true")
	}
}

func TestHasChanges_False(t *testing.T) {
	entries := []DiffEntry{{Key: "x", Change: Unchanged}}
	if HasChanges(entries) {
		t.Error("expected HasChanges to return false")
	}
}
