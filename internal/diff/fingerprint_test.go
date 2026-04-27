package diff

import (
	"strings"
	"testing"
)

var sampleFingerprintEntries = []Entry{
	{Key: "DB_HOST", OldValue: "localhost", NewValue: "prod-db", Status: StatusModified},
	{Key: "API_KEY", OldValue: "", NewValue: "abc123", Status: StatusAdded},
	{Key: "LOG_LEVEL", OldValue: "info", NewValue: "info", Status: StatusUnchanged},
}

func TestComputeFingerprint_NonEmpty(t *testing.T) {
	opts := DefaultFingerprintOptions()
	r := ComputeFingerprint("secret/app", 3, sampleFingerprintEntries, opts)
	if r.Fingerprint == "" {
		t.Fatal("expected non-empty fingerprint")
	}
	if len(r.Fingerprint) != 64 {
		t.Fatalf("expected 64-char hex fingerprint, got %d", len(r.Fingerprint))
	}
}

func TestComputeFingerprint_SetsPathAndVersion(t *testing.T) {
	r := ComputeFingerprint("secret/app", 7, sampleFingerprintEntries, DefaultFingerprintOptions())
	if r.Path != "secret/app" {
		t.Errorf("expected path 'secret/app', got %q", r.Path)
	}
	if r.Version != 7 {
		t.Errorf("expected version 7, got %d", r.Version)
	}
}

func TestComputeFingerprint_ExcludesUnchanged(t *testing.T) {
	opts := DefaultFingerprintOptions()
	opts.IncludeUnchanged = false
	r := ComputeFingerprint("secret/app", 1, sampleFingerprintEntries, opts)
	// Only 2 entries (modified + added)
	if r.EntryCount != 2 {
		t.Errorf("expected 2 entries, got %d", r.EntryCount)
	}
}

func TestComputeFingerprint_Deterministic(t *testing.T) {
	opts := DefaultFingerprintOptions()
	r1 := ComputeFingerprint("secret/app", 1, sampleFingerprintEntries, opts)
	r2 := ComputeFingerprint("secret/app", 1, sampleFingerprintEntries, opts)
	if r1.Fingerprint != r2.Fingerprint {
		t.Error("fingerprint should be deterministic")
	}
}

func TestComputeFingerprint_MaskValues(t *testing.T) {
	opts := DefaultFingerprintOptions()
	opts.MaskValues = false
	r1 := ComputeFingerprint("secret/app", 1, sampleFingerprintEntries, opts)

	opts.MaskValues = true
	r2 := ComputeFingerprint("secret/app", 1, sampleFingerprintEntries, opts)

	// Fingerprints differ when masking is toggled
	if r1.Fingerprint == r2.Fingerprint {
		t.Error("masked and unmasked fingerprints should differ")
	}
}

func TestFormatFingerprint_ContainsFields(t *testing.T) {
	r := ComputeFingerprint("secret/data", 2, sampleFingerprintEntries, DefaultFingerprintOptions())
	out := FormatFingerprint(r)
	for _, want := range []string{"secret/data", "Version:", "Fingerprint:", "Entries:"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}
