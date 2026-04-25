package diff

import (
	"strings"
	"testing"
)

func sampleChecksumSecrets() map[string]string {
	return map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":     "abc123",
		"LOG_LEVEL":   "info",
	}
}

func TestComputeChecksums_EntryCount(t *testing.T) {
	r := ComputeChecksums("secret/app", 3, sampleChecksumSecrets())
	if len(r.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(r.Entries))
	}
}

func TestComputeChecksums_SetsPathAndVersion(t *testing.T) {
	r := ComputeChecksums("secret/app", 7, sampleChecksumSecrets())
	if r.Path != "secret/app" {
		t.Errorf("expected path 'secret/app', got %q", r.Path)
	}
	if r.Version != 7 {
		t.Errorf("expected version 7, got %d", r.Version)
	}
}

func TestComputeChecksums_AggregateNonEmpty(t *testing.T) {
	r := ComputeChecksums("secret/app", 1, sampleChecksumSecrets())
	if r.Aggregate == "" {
		t.Error("expected non-empty aggregate checksum")
	}
}

func TestComputeChecksums_Deterministic(t *testing.T) {
	secrets := sampleChecksumSecrets()
	r1 := ComputeChecksums("secret/app", 1, secrets)
	r2 := ComputeChecksums("secret/app", 1, secrets)
	if r1.Aggregate != r2.Aggregate {
		t.Errorf("expected deterministic aggregate; got %q and %q", r1.Aggregate, r2.Aggregate)
	}
}

func TestComputeChecksums_DifferentValuesProduceDifferentAggregates(t *testing.T) {
	r1 := ComputeChecksums("secret/app", 1, map[string]string{"KEY": "value1"})
	r2 := ComputeChecksums("secret/app", 1, map[string]string{"KEY": "value2"})
	if r1.Aggregate == r2.Aggregate {
		t.Error("expected different aggregates for different values")
	}
}

func TestComputeChecksums_EmptySecrets(t *testing.T) {
	r := ComputeChecksums("secret/empty", 1, map[string]string{})
	if len(r.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(r.Entries))
	}
	if r.Aggregate == "" {
		t.Error("expected non-empty aggregate even for empty secrets")
	}
}

func TestFormatChecksums_ContainsPath(t *testing.T) {
	r := ComputeChecksums("secret/myapp", 2, sampleChecksumSecrets())
	out := FormatChecksums(r)
	if !strings.Contains(out, "secret/myapp") {
		t.Errorf("expected output to contain path, got:\n%s", out)
	}
}

func TestFormatChecksums_ContainsAggregate(t *testing.T) {
	r := ComputeChecksums("secret/myapp", 2, sampleChecksumSecrets())
	out := FormatChecksums(r)
	if !strings.Contains(out, r.Aggregate) {
		t.Errorf("expected output to contain aggregate checksum, got:\n%s", out)
	}
}

func TestFormatChecksums_ContainsKeys(t *testing.T) {
	r := ComputeChecksums("secret/myapp", 2, sampleChecksumSecrets())
	out := FormatChecksums(r)
	for _, e := range r.Entries {
		if !strings.Contains(out, e.Key) {
			t.Errorf("expected output to contain key %q", e.Key)
		}
	}
}
