package diff

import (
	"strings"
	"testing"
)

var sampleDeprecateEntries = []Entry{
	{Key: "legacy_token", Status: StatusModified, OldValue: "abc", NewValue: "xyz"},
	{Key: "old_password", Status: StatusAdded, NewValue: "secret"},
	{Key: "database_url", Status: StatusUnchanged, OldValue: "postgres://", NewValue: "postgres://"},
	{Key: "api_key", Status: StatusRemoved, OldValue: "key123"},
	{Key: "deprecated_host", Status: StatusAdded, NewValue: "host"},
}

func TestDetectDeprecated_MatchesPrefix(t *testing.T) {
	opts := DefaultDeprecateOptions()
	results := DetectDeprecated(sampleDeprecateEntries, opts)
	if len(results) != 2 {
		t.Fatalf("expected 2 deprecated entries, got %d", len(results))
	}
}

func TestDetectDeprecated_ExcludesUnchangedByDefault(t *testing.T) {
	opts := DefaultDeprecateOptions()
	opts.DeprecatedKeys = []string{"database_url"}
	results := DetectDeprecated(sampleDeprecateEntries, opts)
	for _, r := range results {
		if r.Entry.Key == "database_url" {
			t.Error("unchanged entry should be excluded by default")
		}
	}
}

func TestDetectDeprecated_IncludesUnchangedWhenEnabled(t *testing.T) {
	opts := DefaultDeprecateOptions()
	opts.DeprecatedKeys = []string{"database_url"}
	opts.IncludeUnchanged = true
	results := DetectDeprecated(sampleDeprecateEntries, opts)
	found := false
	for _, r := range results {
		if r.Entry.Key == "database_url" {
			found = true
		}
	}
	if !found {
		t.Error("expected database_url to be flagged when IncludeUnchanged is true")
	}
}

func TestDetectDeprecated_ExactKeyMatch(t *testing.T) {
	opts := DeprecateOptions{
		DeprecatedKeys: []string{"api_key"},
	}
	results := DetectDeprecated(sampleDeprecateEntries, opts)
	if len(results) != 1 || results[0].Entry.Key != "api_key" {
		t.Errorf("expected api_key to be flagged, got %+v", results)
	}
}

func TestDetectDeprecated_ReasonContainsKey(t *testing.T) {
	opts := DeprecateOptions{
		DeprecatedKeys: []string{"api_key"},
	}
	results := DetectDeprecated(sampleDeprecateEntries, opts)
	if len(results) == 0 {
		t.Fatal("expected at least one result")
	}
	if !strings.Contains(results[0].Reason, "api_key") {
		t.Errorf("reason should contain key name, got: %s", results[0].Reason)
	}
}

func TestFormatDeprecated_NoEntries(t *testing.T) {
	out := FormatDeprecated(nil)
	if !strings.Contains(out, "no deprecated") {
		t.Errorf("expected no-deprecated message, got: %s", out)
	}
}

func TestFormatDeprecated_ContainsSummary(t *testing.T) {
	opts := DefaultDeprecateOptions()
	results := DetectDeprecated(sampleDeprecateEntries, opts)
	out := FormatDeprecated(results)
	if !strings.Contains(out, "deprecated keys detected") {
		t.Errorf("expected summary line, got: %s", out)
	}
	if !strings.Contains(out, "legacy_token") {
		t.Errorf("expected legacy_token in output, got: %s", out)
	}
}
