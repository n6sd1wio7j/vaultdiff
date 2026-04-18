package main

import (
	"testing"
)

func TestLintFlags_MissingPath(t *testing.T) {
	err := RunLint([]string{"--addr", "http://localhost:8200", "--token", "root"})
	if err == nil || err.Error() != "--path is required" {
		t.Fatalf("expected missing path error, got: %v", err)
	}
}

func TestLintFlags_MissingAddr(t *testing.T) {
	err := RunLint([]string{"--path", "secret/app", "--token", "root"})
	if err == nil || err.Error() != "--addr is required" {
		t.Fatalf("expected missing addr error, got: %v", err)
	}
}

func TestLintFlags_MissingToken(t *testing.T) {
	err := RunLint([]string{"--path", "secret/app", "--addr", "http://localhost:8200"})
	if err == nil || err.Error() != "--token is required" {
		t.Fatalf("expected missing token error, got: %v", err)
	}
}
