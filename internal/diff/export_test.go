package diff

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"strings"
	"testing"
)

func sampleEntries() []Entry {
	return []Entry{
		{Key: "DB_PASS", Status: StatusAdded, OldValue: "", NewValue: "secret123"},
		{Key: "API_KEY", Status: StatusModified, OldValue: "old", NewValue: "new"},
		{Key: "HOST", Status: StatusUnchanged, OldValue: "localhost", NewValue: "localhost"},
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := Export(sampleEntries(), ExportOptions{Format: "xml"}, &buf)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestExport_JSONContainsKeys(t *testing.T) {
	var buf bytes.Buffer
	err := Export(sampleEntries(), ExportOptions{Format: FormatJSON, MaskSecrets: false}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 entries, got %d", len(out))
	}
	if out[0]["key"] != "DB_PASS" {
		t.Errorf("expected DB_PASS, got %v", out[0]["key"])
	}
}

func TestExport_JSONMasksSecrets(t *testing.T) {
	var buf bytes.Buffer
	_ = Export(sampleEntries(), ExportOptions{Format: FormatJSON, MaskSecrets: true}, &buf)
	if strings.Contains(buf.String(), "secret123") {
		t.Error("masked export should not contain raw secret value")
	}
}

func TestExport_CSVHasHeader(t *testing.T) {
	var buf bytes.Buffer
	err := Export(sampleEntries(), ExportOptions{Format: FormatCSV}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := csv.NewReader(&buf)
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("invalid CSV: %v", err)
	}
	if records[0][0] != "key" {
		t.Errorf("expected CSV header 'key', got %q", records[0][0])
	}
	if len(records) != 4 { // header + 3 entries
		t.Errorf("expected 4 rows, got %d", len(records))
	}
}

func TestExport_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	err := Export(sampleEntries(), ExportOptions{Format: FormatText}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}
