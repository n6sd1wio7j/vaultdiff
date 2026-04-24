package diff

import (
	"strings"
	"testing"
)

func sampleClassifyEntries() []Entry {
	return []Entry{
		{Key: "DATABASE_URL", Status: StatusModified, OldValue: "postgres://old", NewValue: "postgres://new"},
		{Key: "APP_ENV", Status: StatusUnchanged, OldValue: "production", NewValue: "production"},
		{Key: "API_TOKEN", Status: StatusAdded, NewValue: "tok-abc123"},
		{Key: "DEPRECATED_KEY", Status: StatusRemoved, OldValue: "old-value"},
	}
}

func TestDefaultClassifyOptions(t *testing.T) {
	opts := DefaultClassifyOptions()
	if opts.ShowAll {
		t.Error("expected ShowAll to be false by default")
	}
	if !opts.MaskSecrets {
		t.Error("expected MaskSecrets to be true by default")
	}
}

func TestFormatClassifiedSummary_Empty(t *testing.T) {
	result := FormatClassifiedSummary(nil, DefaultClassifyOptions())
	if !strings.Contains(result, "no entries") {
		t.Errorf("expected empty message, got: %s", result)
	}
}

func TestFormatClassifiedSummary_HidesUnchangedByDefault(t *testing.T) {
	entries := Classify(sampleClassifyEntries())
	opts := DefaultClassifyOptions()
	result := FormatClassifiedSummary(entries, opts)
	if strings.Contains(result, "APP_ENV") {
		t.Error("expected unchanged entry APP_ENV to be hidden by default")
	}
}

func TestFormatClassifiedSummary_ShowAllIncludesUnchanged(t *testing.T) {
	entries := Classify(sampleClassifyEntries())
	opts := DefaultClassifyOptions()
	opts.ShowAll = true
	result := FormatClassifiedSummary(entries, opts)
	if !strings.Contains(result, "APP_ENV") {
		t.Error("expected APP_ENV to appear when ShowAll is true")
	}
}

func TestFormatClassifiedSummary_MasksSecrets(t *testing.T) {
	entries := Classify(sampleClassifyEntries())
	opts := DefaultClassifyOptions()
	opts.ShowAll = true
	opts.MaskSecrets = true
	result := FormatClassifiedSummary(entries, opts)
	if strings.Contains(result, "tok-abc123") {
		t.Error("expected sensitive value to be masked")
	}
	if !strings.Contains(result, "***") {
		t.Error("expected masked placeholder '***' in output")
	}
}

func TestFormatClassifiedSummary_ContainsSummarySection(t *testing.T) {
	entries := Classify(sampleClassifyEntries())
	opts := DefaultClassifyOptions()
	opts.ShowAll = true
	result := FormatClassifiedSummary(entries, opts)
	if !strings.Contains(result, "summary:") {
		t.Errorf("expected summary section in output, got: %s", result)
	}
}

func TestFormatClassifiedSummary_CountMatchesEntries(t *testing.T) {
	entries := Classify(sampleClassifyEntries())
	opts := DefaultClassifyOptions()
	result := FormatClassifiedSummary(entries, opts)
	if !strings.Contains(result, "classified 4 entries") {
		t.Errorf("expected count of 4, got: %s", result)
	}
}
