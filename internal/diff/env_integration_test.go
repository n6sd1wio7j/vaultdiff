package diff

import (
	"strings"
	"testing"
)

func TestEnvIntegration_FullPipeline(t *testing.T) {
	staging := map[string]string{
		"APP_ENV":      "staging",
		"DB_URL":       "postgres://db-staging/app",
		"SECRET_TOKEN": "stg-token",
		"LOG_LEVEL":    "debug",
		"DEPRECATED":   "old-value",
	}
	production := map[string]string{
		"APP_ENV":   "production",
		"DB_URL":    "postgres://db-prod/app",
		"LOG_LEVEL": "info",
		"NEW_FLAG":  "true",
	}

	opts := EnvCompareOptions{
		EnvA:   "staging",
		EnvB:   "production",
		Ignore: []string{"SECRET_"},
	}

	d := CompareEnvs(staging, production, opts)

	// SECRET_TOKEN must be excluded
	for _, e := range d.Entries {
		if strings.HasPrefix(e.Key, "SECRET_") {
			t.Errorf("secret key leaked into diff: %s", e.Key)
		}
	}

	// Summarize the env diff entries
	summary := Summarize(d.Entries)
	if summary.Modified == 0 {
		t.Error("expected at least one modified key")
	}
	if summary.Added == 0 {
		t.Error("expected at least one added key (NEW_FLAG)")
	}
	if summary.Removed == 0 {
		t.Error("expected at least one removed key (DEPRECATED)")
	}

	// Format and verify output
	out := FormatEnvDiff(d)
	if !strings.Contains(out, "staging") {
		t.Error("output missing env A name")
	}
	if !strings.Contains(out, "production") {
		t.Error("output missing env B name")
	}
}
