package main

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/internal/diff"
)

func TestThresholdFlags_MissingPath(t *testing.T) {
	err := RunThreshold([]string{"--addr", "http://localhost:8200", "--token", "tok"})
	if err == nil || !strings.Contains(err.Error(), "--path") {
		t.Errorf("expected --path error, got: %v", err)
	}
}

func TestThresholdFlags_MissingAddr(t *testing.T) {
	err := RunThreshold([]string{"--path", "secret/app", "--token", "tok"})
	if err == nil || !strings.Contains(err.Error(), "--addr") {
		t.Errorf("expected --addr error, got: %v", err)
	}
}

func TestThresholdFlags_MissingToken(t *testing.T) {
	err := RunThreshold([]string{"--path", "secret/app", "--addr", "http://localhost:8200"})
	if err == nil || !strings.Contains(err.Error(), "--token") {
		t.Errorf("expected --token error, got: %v", err)
	}
}

func TestCheckThreshold_Integration_NoBreachOnLowDrift(t *testing.T) {
	entries := []diff.Entry{
		{Key: "DB_HOST", OldValue: "a", NewValue: "b", Status: diff.StatusModified},
	}
	score := diff.ScoreDrift(entries)
	opts := diff.DefaultThresholdOptions()
	res := diff.CheckThreshold(score, opts)
	if res.Breached {
		t.Errorf("expected no breach for minimal drift, got: %v", res.Violations)
	}
}

func TestCheckThreshold_Integration_BreachOnHighScore(t *testing.T) {
	entries := make([]diff.Entry, 20)
	for i := range entries {
		entries[i] = diff.Entry{Key: "K", Status: diff.StatusAdded, NewValue: "v"}
	}
	score := diff.ScoreDrift(entries)
	opts := diff.DefaultThresholdOptions()
	opts.MaxAdded = 5
	res := diff.CheckThreshold(score, opts)
	if !res.Breached {
		t.Error("expected breach when added keys exceed max")
	}
}
