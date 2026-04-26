package diff

import (
	"strings"
	"testing"
)

var sampleClampEntries = []DiffEntry{
	{Key: "short_key", Status: StatusAdded, NewValue: "hi"},
	{Key: "medium_key", Status: StatusModified, OldValue: "old", NewValue: "newvalue"},
	{Key: "long_key", Status: StatusAdded, NewValue: "this_is_a_very_long_secret_value"},
	{Key: "removed_key", Status: StatusRemoved, OldValue: "gone"},
}

func TestClamp_NoConstraints(t *testing.T) {
	results := Clamp(sampleClampEntries, DefaultClampOptions())
	if len(results) != len(sampleClampEntries) {
		t.Fatalf("expected %d results, got %d", len(sampleClampEntries), len(results))
	}
}

func TestClamp_MinLengthFiltersShort(t *testing.T) {
	opts := DefaultClampOptions()
	opts.MinLength = 5
	results := Clamp(sampleClampEntries, opts)
	for _, r := range results {
		if r.Length < 5 {
			t.Errorf("entry %q has length %d, below min 5", r.Entry.Key, r.Length)
		}
	}
}

func TestClamp_MaxLengthFiltersLong(t *testing.T) {
	opts := DefaultClampOptions()
	opts.MaxLength = 8
	results := Clamp(sampleClampEntries, opts)
	for _, r := range results {
		if r.Length > 8 {
			t.Errorf("entry %q has length %d, above max 8", r.Entry.Key, r.Length)
		}
	}
}

func TestClamp_MaskValuesHidesValue(t *testing.T) {
	opts := DefaultClampOptions()
	opts.MaskValues = true
	results := Clamp(sampleClampEntries, opts)
	for _, r := range results {
		if r.Value != "***" {
			t.Errorf("expected masked value for %q, got %q", r.Entry.Key, r.Value)
		}
	}
}

func TestClamp_UnmaskedShowsValue(t *testing.T) {
	opts := DefaultClampOptions()
	opts.MaskValues = false
	opts.MaxLength = 4
	results := Clamp(sampleClampEntries, opts)
	for _, r := range results {
		if r.Value == "***" {
			t.Errorf("expected real value for %q, got masked", r.Entry.Key)
		}
	}
}

func TestClamp_EmptyEntries(t *testing.T) {
	results := Clamp([]DiffEntry{}, DefaultClampOptions())
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestFormatClamp_NoResults(t *testing.T) {
	out := FormatClamp([]ClampResult{})
	if !strings.Contains(out, "no entries") {
		t.Errorf("expected no-entries message, got: %s", out)
	}
}

func TestFormatClamp_ContainsKey(t *testing.T) {
	opts := DefaultClampOptions()
	results := Clamp(sampleClampEntries, opts)
	out := FormatClamp(results)
	if !strings.Contains(out, "short_key") {
		t.Errorf("expected key in output, got: %s", out)
	}
}
