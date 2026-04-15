package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultdiff/internal/audit"
	"github.com/yourusername/vaultdiff/internal/diff"
)

func TestRecord_WritesJSONLine(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	e := audit.Entry{
		Timestamp:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Environment: "production",
		Path:        "secret/data/app",
		FromVersion: 1,
		ToVersion:   2,
		HasChanges:  true,
		Changes: []diff.DiffEntry{
			{Key: "DB_PASS", Status: diff.Modified},
		},
	}

	if err := l.Record(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	if !strings.HasPrefix(line, "{") {
		t.Fatalf("expected JSON line, got: %s", line)
	}

	var got audit.Entry
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if got.Environment != "production" {
		t.Errorf("environment: want production, got %s", got.Environment)
	}
	if got.ToVersion != 2 {
		t.Errorf("to_version: want 2, got %d", got.ToVersion)
	}
	if len(got.Changes) != 1 {
		t.Errorf("changes: want 1, got %d", len(got.Changes))
	}
}

func TestRecord_SetsTimestampWhenZero(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	e := audit.Entry{Path: "secret/data/svc"}
	if err := l.Record(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got audit.Entry
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if got.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp to be set automatically")
	}
}

func TestNewLogger_DefaultsToStdout(t *testing.T) {
	// Ensure NewLogger(nil) does not panic.
	l := audit.NewLogger(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}
