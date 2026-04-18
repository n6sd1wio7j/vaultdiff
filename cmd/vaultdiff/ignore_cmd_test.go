package main

import (
	"testing"

	"github.com/yourusername/vaultdiff/internal/diff"
)

func TestIgnoreFlags_MissingPath(t *testing.T) {
	err := RunIgnoreDiff([]string{"--addr", "http://localhost:8200", "--token", "tok"})
	if err == nil || err.Error() != "--path is required" {
		t.Fatalf("expected path error, got %v", err)
	}
}

func TestIgnoreFlags_MissingAddr(t *testing.T) {
	err := RunIgnoreDiff([]string{"--path", "secret/data/app", "--token", "tok"})
	if err == nil || err.Error() != "--addr is required" {
		t.Fatalf("expected addr error, got %v", err)
	}
}

func TestIgnoreFlags_MissingToken(t *testing.T) {
	err := RunIgnoreDiff([]string{"--path", "secret/data/app", "--addr", "http://localhost:8200"})
	if err == nil || err.Error() != "--token is required" {
		t.Fatalf("expected token error, got %v", err)
	}
}

func TestApplyIgnore_Integration_KeyAndStatus(t *testing.T) {
	entries := []diff.Entry{
		{Key: "API_KEY", Status: diff.StatusAdded},
		{Key: "LOG_LEVEL", Status: diff.StatusUnchanged},
		{Key: "DB_PASS", Status: diff.StatusModified},
	}
	result := diff.ApplyIgnore(entries, diff.IgnoreOptions{
		Keys:     []string{"API_KEY"},
		Statuses: []string{"unchanged"},
	})
	if len(result) != 1 || result[0].Key != "DB_PASS" {
		t.Fatalf("expected only DB_PASS, got %+v", result)
	}
}
