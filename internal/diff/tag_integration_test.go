package diff

import (
	"strings"
	"testing"
)

func TestTagIntegration_EnvAndVersionTogether(t *testing.T) {
	entries := []DiffEntry{
		{Key: "TOKEN", Status: StatusModified, OldValue: "old", NewValue: "new"},
	}
	tagged := TagEntries(entries, TagOptions{Env: "prod", Version: "v5"})
	if len(tagged) != 1 {
		t.Fatalf("expected 1 entry")
	}
	keys := map[string]bool{}
	for _, tg := range tagged[0].Tags {
		keys[tg.Key] = true
	}
	if !keys["env"] || !keys["version"] {
		t.Errorf("expected both env and version tags, got %v", tagged[0].Tags)
	}
}

func TestTagIntegration_FormatAll(t *testing.T) {
	entries := []DiffEntry{
		{Key: "A", Status: StatusAdded, NewValue: "1"},
		{Key: "B", Status: StatusRemoved, OldValue: "2"},
	}
	tagged := TagEntries(entries, TagOptions{Env: "staging"})
	for _, te := range tagged {
		out := FormatTagged(te)
		if !strings.Contains(out, "env=staging") {
			t.Errorf("missing env tag in: %s", out)
		}
	}
}
