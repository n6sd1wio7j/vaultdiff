package diff

import (
	"strings"
	"testing"
	"time"
)

var blameTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

var sampleBlameEntries = []Entry{
	{Key: "DB_PASS", OldValue: "", NewValue: "secret", Change: ChangeAdded},
	{Key: "API_KEY", OldValue: "old", NewValue: "new", Change: ChangeModified},
	{Key: "HOST", OldValue: "localhost", NewValue: "localhost", Change: ChangeUnchanged},
	{Key: "TOKEN", OldValue: "tok", NewValue: "", Change: ChangeRemoved},
}

func TestBlame_ExcludesUnchanged(t *testing.T) {
	r := Blame("secret/app", sampleBlameEntries, 3, blameTime)
	if _, ok := r.Entries["HOST"]; ok {
		t.Error("expected unchanged key HOST to be excluded from blame")
	}
}

func TestBlame_IncludesChanged(t *testing.T) {
	r := Blame("secret/app", sampleBlameEntries, 3, blameTime)
	for _, key := range []string{"DB_PASS", "API_KEY", "TOKEN"} {
		if _, ok := r.Entries[key]; !ok {
			t.Errorf("expected key %s in blame report", key)
		}
	}
}

func TestBlame_SetsVersion(t *testing.T) {
	r := Blame("secret/app", sampleBlameEntries, 5, blameTime)
	for _, e := range r.Entries {
		if e.Version != 5 {
			t.Errorf("expected version 5, got %d", e.Version)
		}
	}
}

func TestBlame_SetsTimestamp(t *testing.T) {
	r := Blame("secret/app", sampleBlameEntries, 1, blameTime)
	for _, e := range r.Entries {
		if !e.ChangedAt.Equal(blameTime) {
			t.Errorf("expected timestamp %v, got %v", blameTime, e.ChangedAt)
		}
	}
}

func TestFormatBlame_NoChanges(t *testing.T) {
	r := Blame("secret/app", []Entry{}, 1, blameTime)
	out := FormatBlame(r)
	if !strings.Contains(out, "no changes") {
		t.Errorf("expected 'no changes' in output, got: %s", out)
	}
}

func TestFormatBlame_ContainsPath(t *testing.T) {
	r := Blame("secret/myapp", sampleBlameEntries, 2, blameTime)
	out := FormatBlame(r)
	if !strings.Contains(out, "secret/myapp") {
		t.Errorf("expected path in output, got: %s", out)
	}
}

func TestFormatBlame_ContainsKey(t *testing.T) {
	r := Blame("secret/app", sampleBlameEntries, 2, blameTime)
	out := FormatBlame(r)
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in output, got: %s", out)
	}
}
