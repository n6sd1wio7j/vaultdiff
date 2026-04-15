package diff

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

var fixedTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func TestReport_NoChanges(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", OldValue: "bar", NewValue: "bar", Status: StatusUnchanged},
	}
	var buf bytes.Buffer
	err := Report(&buf, entries, ReportOptions{Timestamp: fixedTime})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No changes detected.") {
		t.Errorf("expected no-changes message, got:\n%s", buf.String())
	}
}

func TestReport_SummaryLine(t *testing.T) {
	entries := []Entry{
		{Key: "NEW_KEY", OldValue: "", NewValue: "val", Status: StatusAdded},
		{Key: "OLD_KEY", OldValue: "val", NewValue: "", Status: StatusRemoved},
		{Key: "MOD_KEY", OldValue: "a", NewValue: "b", Status: StatusModified},
	}
	var buf bytes.Buffer
	err := Report(&buf, entries, ReportOptions{Timestamp: fixedTime})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "+1 added") {
		t.Errorf("expected '+1 added' in output:\n%s", out)
	}
	if !strings.Contains(out, "-1 removed") {
		t.Errorf("expected '-1 removed' in output:\n%s", out)
	}
	if !strings.Contains(out, "~1 modified") {
		t.Errorf("expected '~1 modified' in output:\n%s", out)
	}
}

func TestReport_HeaderContainsPath(t *testing.T) {
	entries := []Entry{
		{Key: "X", OldValue: "1", NewValue: "2", Status: StatusModified},
	}
	var buf bytes.Buffer
	opts := ReportOptions{
		SecretPath: "secret/myapp/prod",
		SourceEnv:  "staging",
		TargetEnv:  "production",
		Timestamp:  fixedTime,
	}
	_ = Report(&buf, entries, opts)
	out := buf.String()
	if !strings.Contains(out, "secret/myapp/prod") {
		t.Errorf("expected secret path in header:\n%s", out)
	}
	if !strings.Contains(out, "staging") || !strings.Contains(out, "production") {
		t.Errorf("expected env names in header:\n%s", out)
	}
}

func TestReport_TimestampInHeader(t *testing.T) {
	var buf bytes.Buffer
	_ = Report(&buf, []Entry{}, ReportOptions{Timestamp: fixedTime})
	if !strings.Contains(buf.String(), "2024-06-01T12:00:00Z") {
		t.Errorf("expected timestamp in header:\n%s", buf.String())
	}
}
