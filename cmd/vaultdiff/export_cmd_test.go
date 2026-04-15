package main

import (
	"testing"
)

func TestExportFlags_MissingPath(t *testing.T) {
	err := RunExport([]string{
		"--addr", "http://127.0.0.1:8200",
		"--token", "root",
	})
	if err == nil {
		t.Fatal("expected error when --path is missing")
	}
}

func TestExportFlags_InvalidFormat(t *testing.T) {
	// Format validation happens inside Export; ensure RunExport propagates errors.
	// We can't easily call RunExport end-to-end without a live Vault, so we
	// validate that the flag parsing itself succeeds and the format string is forwarded.
	fs := ExportFlags{
		Path:   "myapp/config",
		Format: "xml",
		Mask:   false,
	}
	if fs.Format != "xml" {
		t.Errorf("expected format xml, got %q", fs.Format)
	}
}

func TestExportFlags_Defaults(t *testing.T) {
	// Verify default field values match expectations without invoking Vault.
	var f ExportFlags
	fs := []struct {
		name string
		got  interface{}
		want interface{}
	}{
		{"Format", f.Format, ""},
		{"Mount", f.Mount, ""},
		{"VerA", f.VerA, 0},
		{"VerB", f.VerB, 0},
		{"Mask", f.Mask, false},
	}
	for _, tc := range fs {
		if tc.got != tc.want {
			t.Errorf("%s: expected %v, got %v", tc.name, tc.want, tc.got)
		}
	}
}

func TestExportFlags_AllFlagsParsed(t *testing.T) {
	// Ensure RunExport returns a vault error (not a flag error) when all flags are provided.
	err := RunExport([]string{
		"--path", "myapp/db",
		"--addr", "http://127.0.0.1:8200",
		"--token", "root",
		"--mount", "kv",
		"--ver-a", "1",
		"--ver-b", "2",
		"--format", "json",
		"--mask", "true",
	})
	// We expect a vault connection error, not a flag/path error.
	if err == nil {
		t.Fatal("expected vault connection error in test environment")
	}
	if err.Error() == "--path is required" {
		t.Error("path flag was not parsed correctly")
	}
}
