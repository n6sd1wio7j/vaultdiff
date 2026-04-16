package diff

import (
	"strings"
	"testing"
)

func sampleAnnotateEntries() []Entry {
	return []Entry{
		{Key: "DB_PASS", Change: Modified, OldValue: "old", NewValue: "new"},
		{Key: "API_KEY", Change: Added, NewValue: "abc123"},
		{Key: "HOST", Change: Unchanged, NewValue: "localhost"},
	}
}

func TestAnnotate_ReturnsMatchingNotes(t *testing.T) {
	entries := sampleAnnotateEntries()
	opts := AnnotateOptions{
		Annotations: []Annotation{
			{Key: "DB_PASS", Note: "rotated by ops"},
			{Key: "API_KEY", Note: "new integration"},
		},
	}
	notes := Annotate(entries, opts)
	if notes["DB_PASS"] != "rotated by ops" {
		t.Errorf("expected note for DB_PASS")
	}
	if notes["API_KEY"] != "new integration" {
		t.Errorf("expected note for API_KEY")
	}
}

func TestAnnotate_NoMatchReturnsEmpty(t *testing.T) {
	entries := sampleAnnotateEntries()
	opts := AnnotateOptions{
		Annotations: []Annotation{{Key: "UNKNOWN", Note: "ghost"}},
	}
	notes := Annotate(entries, opts)
	if len(notes) != 0 {
		t.Errorf("expected empty notes map, got %d entries", len(notes))
	}
}

func TestAnnotate_EmptyAnnotations(t *testing.T) {
	entries := sampleAnnotateEntries()
	notes := Annotate(entries, AnnotateOptions{})
	if len(notes) != 0 {
		t.Errorf("expected empty map")
	}
}

func TestFormatAnnotated_WithNote(t *testing.T) {
	e := Entry{Key: "DB_PASS", Change: Modified, NewValue: "secret"}
	line := FormatAnnotated(e, "rotated", false)
	if !strings.Contains(line, "# rotated") {
		t.Errorf("expected annotation in line: %s", line)
	}
}

func TestFormatAnnotated_NoNote(t *testing.T) {
	e := Entry{Key: "HOST", Change: Unchanged, NewValue: "localhost"}
	line := FormatAnnotated(e, "", false)
	if strings.Contains(line, "#") {
		t.Errorf("unexpected annotation marker in line: %s", line)
	}
}

func TestFormatAnnotated_MasksValue(t *testing.T) {
	e := Entry{Key: "SECRET", Change: Added, NewValue: "topsecret"}
	line := FormatAnnotated(e, "sensitive", true)
	if strings.Contains(line, "topsecret") {
		t.Errorf("value should be masked")
	}
}
