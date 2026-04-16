package main

import (
	"testing"
)

func TestWatchFlags_MissingPath(t *testing.T) {
	err := RunWatch([]string{"--addr", "http://localhost:8200", "--token", "tok"})
	if err == nil || err.Error() != "--path is required" {
		t.Errorf("expected path error, got %v", err)
	}
}

func TestWatchFlags_MissingAddr(t *testing.T) {
	err := RunWatch([]string{"--token", "tok", "--path", "secret/foo"})
	if err == nil {
		t.Error("expected error for missing addr")
	}
}

func TestWatchFlags_MissingToken(t *testing.T) {
	err := RunWatch([]string{"--addr", "http://localhost:8200", "--path", "secret/foo"})
	if err == nil {
		t.Error("expected error for missing token")
	}
}
