package diff

import (
	"strings"
	"testing"
)

var sampleEnvA = map[string]string{
	"APP_ENV":      "staging",
	"DB_HOST":      "db-staging.internal",
	"SECRET_TOKEN": "abc123",
	"LOG_LEVEL":    "debug",
}

var sampleEnvB = map[string]string{
	"APP_ENV":      "production",
	"DB_HOST":      "db-prod.internal",
	"SECRET_TOKEN": "abc123",
	"NEW_FEATURE":  "enabled",
}

func TestCompareEnvs_DetectsModified(t *testing.T) {
	opts := DefaultEnvCompareOptions()
	d := CompareEnvs(sampleEnvA, sampleEnvB, opts)

	modified := 0
	for _, e := range d.Entries {
		if e.Status == StatusModified {
			modified++
		}
	}
	if modified == 0 {
		t.Error("expected at least one modified entry")
	}
}

func TestCompareEnvs_DetectsAdded(t *testing.T) {
	opts := DefaultEnvCompareOptions()
	d := CompareEnvs(sampleEnvA, sampleEnvB, opts)

	for _, e := range d.Entries {
		if e.Key == "NEW_FEATURE" && e.Status != StatusAdded {
			t.Errorf("expected NEW_FEATURE to be added, got %s", e.Status)
		}
	}
}

func TestCompareEnvs_DetectsRemoved(t *testing.T) {
	opts := DefaultEnvCompareOptions()
	d := CompareEnvs(sampleEnvA, sampleEnvB, opts)

	for _, e := range d.Entries {
		if e.Key == "LOG_LEVEL" && e.Status != StatusRemoved {
			t.Errorf("expected LOG_LEVEL to be removed, got %s", e.Status)
		}
	}
}

func TestCompareEnvs_IgnorePrefix(t *testing.T) {
	opts := DefaultEnvCompareOptions()
	opts.Ignore = []string{"SECRET_"}
	d := CompareEnvs(sampleEnvA, sampleEnvB, opts)

	for _, e := range d.Entries {
		if strings.HasPrefix(e.Key, "SECRET_") {
			t.Errorf("expected SECRET_ keys to be ignored, got %s", e.Key)
		}
	}
}

func TestCompareEnvs_EnvNames(t *testing.T) {
	opts := EnvCompareOptions{EnvA: "dev", EnvB: "prod"}
	d := CompareEnvs(sampleEnvA, sampleEnvB, opts)
	if d.EnvA != "dev" || d.EnvB != "prod" {
		t.Errorf("unexpected env names: %s, %s", d.EnvA, d.EnvB)
	}
}

func TestFormatEnvDiff_ContainsHeader(t *testing.T) {
	opts := DefaultEnvCompareOptions()
	d := CompareEnvs(sampleEnvA, sampleEnvB, opts)
	out := FormatEnvDiff(d)
	if !strings.Contains(out, "staging") || !strings.Contains(out, "production") {
		t.Errorf("expected env names in output, got: %s", out)
	}
}

func TestFormatEnvDiff_NoDifferences(t *testing.T) {
	opts := DefaultEnvCompareOptions()
	d := CompareEnvs(sampleEnvA, sampleEnvA, opts)
	out := FormatEnvDiff(d)
	if !strings.Contains(out, "no differences") {
		t.Errorf("expected no differences message, got: %s", out)
	}
}
