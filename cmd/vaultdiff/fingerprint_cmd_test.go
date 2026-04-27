package main

import (
	"testing"
)

func TestFingerprintFlags_MissingPath(t *testing.T) {
	err := RunFingerprint([]string{
		"--addr", "http://127.0.0.1:8200",
		"--token", "root",
	})
	if err == nil || err.Error() != "--path is required" {
		t.Errorf("expected missing path error, got %v", err)
	}
}

func TestFingerprintFlags_MissingAddr(t *testing.T) {
	err := RunFingerprint([]string{
		"--path", "secret/app",
		"--token", "root",
	})
	if err == nil || err.Error() != "--addr is required" {
		t.Errorf("expected missing addr error, got %v", err)
	}
}

func TestFingerprintFlags_MissingToken(t *testing.T) {
	err := RunFingerprint([]string{
		"--path", "secret/app",
		"--addr", "http://127.0.0.1:8200",
	})
	if err == nil || err.Error() != "--token is required" {
		t.Errorf("expected missing token error, got %v", err)
	}
}
