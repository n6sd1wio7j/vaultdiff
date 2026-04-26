package main

import (
	"testing"
)

func TestDriftReportFlags_MissingPath(t *testing.T) {
	err := RunDriftReport([]string{
		"--addr", "http://localhost:8200",
		"--token", "root",
		"--version-a", "1",
		"--version-b", "2",
	})
	if err == nil || err.Error() != "--path is required" {
		t.Errorf("expected missing path error, got %v", err)
	}
}

func TestDriftReportFlags_MissingAddr(t *testing.T) {
	t.Setenv("VAULT_ADDR", "")
	err := RunDriftReport([]string{
		"--path", "secret/app",
		"--token", "root",
		"--version-a", "1",
		"--version-b", "2",
	})
	if err == nil || err.Error() != "--addr is required" {
		t.Errorf("expected missing addr error, got %v", err)
	}
}

func TestDriftReportFlags_MissingToken(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "")
	err := RunDriftReport([]string{
		"--path", "secret/app",
		"--addr", "http://localhost:8200",
		"--version-a", "1",
		"--version-b", "2",
	})
	if err == nil || err.Error() != "--token is required" {
		t.Errorf("expected missing token error, got %v", err)
	}
}

func TestDriftReportFlags_MissingVersions(t *testing.T) {
	err := RunDriftReport([]string{
		"--path", "secret/app",
		"--addr", "http://localhost:8200",
		"--token", "root",
	})
	if err == nil {
		t.Error("expected error for missing versions")
	}
}
