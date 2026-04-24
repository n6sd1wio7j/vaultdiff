package diff

import (
	"strings"
	"testing"
)

var sampleRenameEntries = []Entry{
	{Key: "old_db_pass", Status: StatusRemoved, OldValue: "s3cr3t", NewValue: ""},
	{Key: "db_password", Status: StatusAdded, OldValue: "", NewValue: "s3cr3t"},
	{Key: "api_key", Status: StatusRemoved, OldValue: "key123", NewValue: ""},
	{Key: "service_api_key", Status: StatusAdded, OldValue: "", NewValue: "key123"},
	{Key: "unchanged_key", Status: StatusUnchanged, OldValue: "abc", NewValue: "abc"},
	{Key: "modified_key", Status: StatusModified, OldValue: "foo", NewValue: "bar"},
}

func TestDetectRenames_FindsRenamedKeys(t *testing.T) {
	opts := DefaultRenameOptions()
	opts.MaskValues = false
	renames := DetectRenames(sampleRenameEntries, opts)
	if len(renames) != 2 {
		t.Fatalf("expected 2 renames, got %d", len(renames))
	}
}

func TestDetectRenames_MasksValuesByDefault(t *testing.T) {
	opts := DefaultRenameOptions()
	renames := DetectRenames(sampleRenameEntries, opts)
	for _, r := range renames {
		if r.Value != "***" {
			t.Errorf("expected masked value, got %q", r.Value)
		}
	}
}

func TestDetectRenames_ExcludesUnchangedAndModified(t *testing.T) {
	entries := []Entry{
		{Key: "x", Status: StatusUnchanged, OldValue: "v", NewValue: "v"},
		{Key: "y", Status: StatusModified, OldValue: "a", NewValue: "b"},
	}
	renames := DetectRenames(entries, DefaultRenameOptions())
	if len(renames) != 0 {
		t.Errorf("expected no renames, got %d", len(renames))
	}
}

func TestDetectRenames_RespectsMinValueLength(t *testing.T) {
	entries := []Entry{
		{Key: "old", Status: StatusRemoved, OldValue: "x", NewValue: ""},
		{Key: "new", Status: StatusAdded, OldValue: "", NewValue: "x"},
	}
	opts := RenameOptions{MinValueLength: 5, MaskValues: false}
	renames := DetectRenames(entries, opts)
	if len(renames) != 0 {
		t.Errorf("expected no renames due to min length, got %d", len(renames))
	}
}

func TestFormatRenames_NoRenames(t *testing.T) {
	out := FormatRenames(nil)
	if !strings.Contains(out, "no renames") {
		t.Errorf("expected 'no renames' message, got %q", out)
	}
}

func TestFormatRenames_ContainsOldAndNewKey(t *testing.T) {
	renames := []RenameEntry{
		{OldKey: "old_key", NewKey: "new_key", Value: "***"},
	}
	out := FormatRenames(renames)
	if !strings.Contains(out, "old_key") {
		t.Errorf("expected old_key in output")
	}
	if !strings.Contains(out, "new_key") {
		t.Errorf("expected new_key in output")
	}
	if !strings.Contains(out, "1 rename") {
		t.Errorf("expected rename count in output")
	}
}
