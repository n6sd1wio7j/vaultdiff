package diff

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

var referenceNow = time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

func sampleStaleEntries() []Entry {
	// One entry carries a timestamp suffix indicating rotation 120 days ago.
	old := referenceNow.Add(-120 * 24 * time.Hour).Unix()
	// Another was rotated 10 days ago — well within the window.
	recent := referenceNow.Add(-10 * 24 * time.Hour).Unix()
	return []Entry{
		{Key: fmt.Sprintf("OLD_PASSWORD@%d", old), Status: StatusUnchanged},
		{Key: fmt.Sprintf("NEW_TOKEN@%d", recent), Status: StatusUnchanged},
		{Key: "MODIFIED_KEY", Status: StatusModified},
	}
}

func TestDetectStale_MarksOldEntryStale(t *testing.T) {
	entries := sampleStaleEntries()
	opts := DefaultStaleOptions()
	opts.IncludeUnchanged = true
	results := DetectStale(entries, referenceNow, opts)

	var found *StaleEntry
	for i := range results {
		if strings.HasPrefix(results[i].Key, "OLD_PASSWORD") {
			found = &results[i]
		}
	}
	if found == nil {
		t.Fatal("expected OLD_PASSWORD entry in results")
	}
	if !found.Stale {
		t.Errorf("expected OLD_PASSWORD to be stale")
	}
}

func TestDetectStale_RecentEntryNotStale(t *testing.T) {
	entries := sampleStaleEntries()
	opts := DefaultStaleOptions()
	opts.IncludeUnchanged = true
	results := DetectStale(entries, referenceNow, opts)

	for _, r := range results {
		if strings.HasPrefix(r.Key, "NEW_TOKEN") && r.Stale {
			t.Errorf("expected NEW_TOKEN to not be stale")
		}
	}
}

func TestDetectStale_ExcludesUnchangedByDefault(t *testing.T) {
	entries := sampleStaleEntries()
	opts := DefaultStaleOptions() // IncludeUnchanged = false
	results := DetectStale(entries, referenceNow, opts)

	for _, r := range results {
		if r.Status == StatusUnchanged {
			t.Errorf("unexpected unchanged entry %q in results", r.Key)
		}
	}
}

func TestDetectStale_ModifiedEntryIncluded(t *testing.T) {
	entries := sampleStaleEntries()
	opts := DefaultStaleOptions()
	results := DetectStale(entries, referenceNow, opts)

	found := false
	for _, r := range results {
		if r.Key == "MODIFIED_KEY" {
			found = true
		}
	}
	if !found {
		t.Error("expected MODIFIED_KEY in results")
	}
}

func TestDetectStale_ReasonContainsDays(t *testing.T) {
	old := referenceNow.Add(-200 * 24 * time.Hour).Unix()
	entries := []Entry{
		{Key: fmt.Sprintf("ANCIENT_SECRET@%d", old), Status: StatusModified},
	}
	opts := DefaultStaleOptions()
	results := DetectStale(entries, referenceNow, opts)
	if len(results) == 0 {
		t.Fatal("expected at least one result")
	}
	if !strings.Contains(results[0].Reason, "200") {
		t.Errorf("expected reason to contain day count, got: %s", results[0].Reason)
	}
}

func TestFormatStale_NoEntries(t *testing.T) {
	out := FormatStale(nil)
	if !strings.Contains(out, "no stale") {
		t.Errorf("expected no-stale message, got: %s", out)
	}
}

func TestFormatStale_ContainsKey(t *testing.T) {
	entries := []StaleEntry{
		{Key: "MY_SECRET", Status: StatusModified, Stale: true, Reason: "not rotated in 100 days (max 90)"},
	}
	out := FormatStale(entries)
	if !strings.Contains(out, "MY_SECRET") {
		t.Errorf("expected key in output, got: %s", out)
	}
	if !strings.Contains(out, "!") {
		t.Errorf("expected stale marker in output, got: %s", out)
	}
}
