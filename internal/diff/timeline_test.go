package diff

import (
	"strings"
	"testing"
	"time"
)

var sampleTimelineEntries = []Entry{
	{Key: "DB_PASS", OldValue: "old", NewValue: "new", Status: Modified},
	{Key: "API_KEY", OldValue: "", NewValue: "abc", Status: Added},
	{Key: "HOST", OldValue: "h", NewValue: "h", Status: Unchanged},
}

func TestTimeline_Append_AddsEntry(t *testing.T) {
	var tl Timeline
	e := tl.Append("/secret/app", sampleTimelineEntries)
	if len(tl) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(tl))
	}
	if e.Path != "/secret/app" {
		t.Errorf("expected path /secret/app, got %s", e.Path)
	}
	if e.Score.Added != 1 {
		t.Errorf("expected 1 added, got %d", e.Score.Added)
	}
	if e.Score.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", e.Score.Modified)
	}
}

func TestTimeline_Sorted_OrdersByTime(t *testing.T) {
	now := time.Now().UTC()
	tl := Timeline{
		{Timestamp: now.Add(2 * time.Minute), Path: "b"},
		{Timestamp: now.Add(1 * time.Minute), Path: "a"},
		{Timestamp: now.Add(3 * time.Minute), Path: "c"},
	}
	sorted := tl.Sorted()
	if sorted[0].Path != "a" || sorted[1].Path != "b" || sorted[2].Path != "c" {
		t.Errorf("unexpected order: %v", sorted)
	}
}

func TestFormatTimeline_Empty(t *testing.T) {
	out := FormatTimeline(Timeline{})
	if !strings.Contains(out, "no timeline") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatTimeline_ContainsPath(t *testing.T) {
	var tl Timeline
	tl.Append("/secret/myapp", sampleTimelineEntries)
	out := FormatTimeline(tl)
	if !strings.Contains(out, "/secret/myapp") {
		t.Errorf("expected path in output, got:\n%s", out)
	}
}

func TestFormatTimeline_ContainsHeader(t *testing.T) {
	var tl Timeline
	tl.Append("/secret/x", sampleTimelineEntries)
	out := FormatTimeline(tl)
	if !strings.Contains(out, "Score") {
		t.Errorf("expected Score header, got:\n%s", out)
	}
}
