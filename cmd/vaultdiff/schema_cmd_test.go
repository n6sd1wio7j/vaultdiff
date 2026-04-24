package main

import (
	"testing"
)

func TestSchemaFlags_MissingPath(t *testing.T) {
	err := RunSchema([]string{
		"--addr", "http://127.0.0.1:8200",
		"--token", "root",
	})
	if err == nil {
		t.Fatal("expected error for missing --path")
	}
	if err.Error() != "--path is required" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSchemaFlags_MissingAddr(t *testing.T) {
	err := RunSchema([]string{
		"--path", "secret/app",
		"--token", "root",
	})
	if err == nil {
		t.Fatal("expected error for missing --addr")
	}
	if err.Error() != "--addr or VAULT_ADDR is required" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSchemaFlags_MissingToken(t *testing.T) {
	err := RunSchema([]string{
		"--path", "secret/app",
		"--addr", "http://127.0.0.1:8200",
	})
	if err == nil {
		t.Fatal("expected error for missing --token")
	}
	if err.Error() != "--token or VAULT_TOKEN is required" {
		t.Errorf("unexpected error: %v", err)
	}
}
