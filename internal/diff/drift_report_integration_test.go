package diff

import (
	"strings"
	"testing"
)

func TestDriftReportIntegration_FullPipeline(t *testing.T) {
	secA := map[string]string{
		"DB_HOST":   "localhost",
		"API_KEY":   "old-key",
		"LOG_LEVEL": "debug",
	}
	secB := map[string]string{
		"DB_HOST":  "prod-db",
		"API_KEY":  "new-key",
		"NEW_FLAG": "true",
	}

	entries := Compare(secA, secB)

	opts := DriftReportOptions{
		Env:            "production",
		Path:           "secret/myapp",
		VersionA:       5,
		VersionB:       6,
		IncludeScore:   true,
		IncludeSummary: true,
	}

	report := BuildDriftReport(entries, opts)

	if report.Summary.Modified != 2 {
		t.Errorf("expected 2 modified, got %d", report.Summary.Modified)
	}
	if report.Summary.Added != 1 {
		t.Errorf("expected 1 added, got %d", report.Summary.Added)
	}
	if report.Summary.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", report.Summary.Removed)
	}
	if report.Score.Score == 0 {
		t.Error("expected non-zero score")
	}

	out := FormatDriftReport(report, opts)

	for _, want := range []string{"production", "secret/myapp", "v5 → v6", "Added:", "Score"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}
