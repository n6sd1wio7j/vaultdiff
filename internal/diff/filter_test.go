package diff

import (
	"testing"
)

func TestFilter_OnlyChanged(t *testing.T) {
	entries := []Entry{
		{Key: "a", Status: StatusAdded},
		{Key: "b", Status: StatusUnchanged},
		{Key: "c", Status: StatusRemoved},
	}
	result := Filter(entries, FilterOptions{OnlyChanged: true})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	for _, e := range result {
		if e.Status == StatusUnchanged {
			t.Errorf("unexpected unchanged entry: %s", e.Key)
		}
	}
}

func TestFilter_KeyPrefix(t *testing.T) {
	entries := []Entry{
		{Key: "db_host", Status: StatusAdded},
		{Key: "db_pass", Status: StatusModified},
		{Key: "api_key", Status: StatusUnchanged},
	}
	result := Filter(entries, FilterOptions{KeyPrefix: "db_"})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	for _, e := range result {
		if e.Key == "api_key" {
			t.Error("api_key should have been filtered out")
		}
	}
}

func TestFilter_ExcludeKeys(t *testing.T) {
	entries := []Entry{
		{Key: "secret", Status: StatusAdded},
		{Key: "token", Status: StatusAdded},
		{Key: "name", Status: StatusAdded},
	}
	result := Filter(entries, FilterOptions{ExcludeKeys: []string{"secret", "token"}})
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].Key != "name" {
		t.Errorf("expected key 'name', got %s", result[0].Key)
	}
}

func TestFilter_CombinedOptions(t *testing.T) {
	entries := []Entry{
		{Key: "db_host", Status: StatusUnchanged},
		{Key: "db_pass", Status: StatusModified},
		{Key: "db_user", Status: StatusAdded},
	}
	result := Filter(entries, FilterOptions{
		OnlyChanged:  true,
		KeyPrefix:    "db_",
		ExcludeKeys:  []string{"db_user"},
	})
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].Key != "db_pass" {
		t.Errorf("expected db_pass, got %s", result[0].Key)
	}
}

func TestFilter_EmptyOptions(t *testing.T) {
	entries := []Entry{
		{Key: "a", Status: StatusAdded},
		{Key: "b", Status: StatusUnchanged},
	}
	result := Filter(entries, FilterOptions{})
	if len(result) != 2 {
		t.Fatalf("expected all 2 entries, got %d", len(result))
	}
}
