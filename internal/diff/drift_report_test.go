package diff

import (
	"strings"
	"testing"
)

var sampleDriftEntries = []Entry{
	{Key: "DB_HOST", OldValue: "localhost", NewValue: "prod-db", Status: StatusModified},
	{Key: "API_KEY", OldValue: "", NewValue: "secret", Status: StatusAdded},
	{Key: "LOG_LEVEL", OldValue: "debug", NewValue: "", Status: StatusRemoved},
	{Key: "PORT", OldValue: "8080", NewValue: "8080", Status: StatusUnchanged},
}

func TestBuildDriftReport_FieldsSet(t *testing.T) {
	opts := DriftReportOptions{Env: "staging", Path: "secret/app", VersionA: 1, VersionB: 2}
	r := BuildDriftReport(sampleDriftEntries, opts)

	if r.Env != "staging" {
		t.Errorf("expected env staging, got %s", r.Env)
	}
	if r.Path != "secret/app" {
		t.Errorf("expected path secret/app, got %s", r.Path)
	}
	if r.VersionA != 1 || r.VersionB != 2 {
		t.Errorf("unexpected versions: %d %d", r.VersionA, r.VersionB)
	}
	if r.GeneratedAt.IsZero() {
		t.Error("expected GeneratedAt to be set")
	}
}

func TestBuildDriftReport_SummaryCounts(t *testing.T) {
	opts := DriftReportOptions{}
	r := BuildDriftReport(sampleDriftEntries, opts)

	if r.Summary.Added != 1 || r.Summary.Removed != 1 || r.Summary.Modified != 1 {
		t.Errorf("unexpected summary: %+v", r.Summary)
	}
}

func TestBuildDriftReport_ScoreNonZero(t *testing.T) {
	opts := DriftReportOptions{}
	r := BuildDriftReport(sampleDriftEntries, opts)

	if r.Score.Score == 0 {
		t.Error("expected non-zero drift score")
	}
}

func TestFormatDriftReport_ContainsPath(t *testing.T) {
	opts := DriftReportOptions{Path: "secret/myapp", VersionA: 3, VersionB: 4, IncludeSummary: true}
	r := BuildDriftReport(sampleDriftEntries, opts)
	out := FormatDriftReport(r, opts)

	if !strings.Contains(out, "secret/myapp") {
		t.Error("expected output to contain path")
	}
	if !strings.Contains(out, "v3 → v4") {
		t.Error("expected output to contain version range")
	}
}

func TestFormatDriftReport_SummaryLine(t *testing.T) {
	opts := DriftReportOptions{IncludeSummary: true}
	r := BuildDriftReport(sampleDriftEntries, opts)
	out := FormatDriftReport(r, opts)

	if !strings.Contains(out, "Added:") {
		t.Error("expected summary line in output")
	}
}

func TestFormatDriftReport_ScoreLine(t *testing.T) {
	opts := DriftReportOptions{IncludeScore: true}
	r := BuildDriftReport(sampleDriftEntries, opts)
	out := FormatDriftReport(r, opts)

	if !strings.Contains(out, "Score") {
		t.Error("expected score in output")
	}
}

func TestFormatDriftReport_EnvOptional(t *testing.T) {
	opts := DriftReportOptions{Env: ""}
	r := BuildDriftReport(sampleDriftEntries, opts)
	out := FormatDriftReport(r, opts)

	if strings.Contains(out, "Env:") {
		t.Error("expected no Env line when env is empty")
	}
}
