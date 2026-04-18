package diff

import (
	"strings"
	"testing"
)

func sampleTagEntries() []DiffEntry {
	return []DiffEntry{
		{Key: "API_KEY", Status: StatusAdded, NewValue: "abc"},
		{Key: "DB_PASS", Status: StatusModified, OldValue: "x", NewValue: "y"},
	}
}

func TestTagEntries_EnvTag(t *testing.T) {
	tagged := TagEntries(sampleTagEntries(), TagOptions{Env: "production"})
	if len(tagged) != 2 {
		t.Fatalf("expected 2 tagged entries, got %d", len(tagged))
	}
	if tagged[0].Tags[0].Value != "production" {
		t.Errorf("expected env=production, got %s", tagged[0].Tags[0].Value)
	}
}

func TestTagEntries_VersionTag(t *testing.T) {
	tagged := TagEntries(sampleTagEntries(), TagOptions{Version: "v3"})
	if tagged[0].Tags[0].Key != "version" {
		t.Errorf("expected version tag, got %s", tagged[0].Tags[0].Key)
	}
}

func TestTagEntries_CustomTags(t *testing.T) {
	opts := TagOptions{Custom: map[string]string{"team": "platform"}}
	tagged := TagEntries(sampleTagEntries(), opts)
	found := false
	for _, tg := range tagged[0].Tags {
		if tg.Key == "team" && tg.Value == "platform" {
			found = true
		}
	}
	if !found {
		t.Error("expected custom tag team=platform")
	}
}

func TestTagEntries_NoOptions(t *testing.T) {
	tagged := TagEntries(sampleTagEntries(), TagOptions{})
	for _, te := range tagged {
		if len(te.Tags) != 0 {
			t.Errorf("expected no tags, got %v", te.Tags)
		}
	}
}

func TestFormatTagged_ContainsKey(t *testing.T) {
	te := TaggedEntry{
		Entry: DiffEntry{Key: "SECRET", Status: StatusAdded},
		Tags:  []Tag{{Key: "env", Value: "staging"}},
	}
	out := FormatTagged(te)
	if !strings.Contains(out, "SECRET") {
		t.Errorf("expected key in output: %s", out)
	}
	if !strings.Contains(out, "env=staging") {
		t.Errorf("expected tag in output: %s", out)
	}
}

func TestFormatTagged_StatusPrefix(t *testing.T) {
	te := TaggedEntry{
		Entry: DiffEntry{Key: "X", Status: StatusRemoved},
	}
	out := FormatTagged(te)
	if !strings.Contains(out, string(StatusRemoved)) {
		t.Errorf("expected status in output: %s", out)
	}
}
