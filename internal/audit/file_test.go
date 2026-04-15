package audit_test

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/internal/audit"
)

func TestOpenFileLogger_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "logs", "audit.jsonl")

	fl, err := audit.OpenFileLogger(path)
	if err != nil {
		t.Fatalf("OpenFileLogger: %v", err)
	}
	defer fl.Close()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected log file to be created")
	}
}

func TestOpenFileLogger_AppendsEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.jsonl")

	for i := 0; i < 3; i++ {
		fl, err := audit.OpenFileLogger(path)
		if err != nil {
			t.Fatalf("open iteration %d: %v", i, err)
		}
		_ = fl.Record(audit.Entry{Path: "secret/data/app", ToVersion: i + 1})
		fl.Close()
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open file: %v", err)
	}
	defer f.Close()

	var count int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var e audit.Entry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			t.Errorf("invalid JSON on line %d: %v", count+1, err)
		}
		count++
	}
	if count != 3 {
		t.Errorf("expected 3 log entries, got %d", count)
	}
}

func TestDefaultLogPath_ContainsDate(t *testing.T) {
	p := audit.DefaultLogPath("/var/log/vaultdiff")
	if !strings.Contains(p, "vaultdiff-") {
		t.Errorf("unexpected log path format: %s", p)
	}
	if filepath.Ext(p) != ".jsonl" {
		t.Errorf("expected .jsonl extension, got: %s", filepath.Ext(p))
	}
}
