package main

import (
	"testing"
)

func TestParseFlags_MissingPath(t *testing.T) {
	_, err := ParseFlags([]string{"-addr", "http://localhost:8200", "-token", "root"})
	if err == nil {
		t.Fatal("expected error for missing -path, got nil")
	}
}

func TestParseFlags_MissingAddr(t *testing.T) {
	_, err := ParseFlags([]string{"-token", "root", "-path", "myapp/config"})
	if err == nil {
		t.Fatal("expected error for missing addr, got nil")
	}
}

func TestParseFlags_MissingToken(t *testing.T) {
	_, err := ParseFlags([]string{"-addr", "http://localhost:8200", "-path", "myapp/config"})
	if err == nil {
		t.Fatal("expected error for missing token, got nil")
	}
}

func TestParseFlags_Defaults(t *testing.T) {
	cfg, err := ParseFlags([]string{
		"-addr", "http://localhost:8200",
		"-token", "root",
		"-path", "myapp/config",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Mount != "secret" {
		t.Errorf("expected default mount 'secret', got %q", cfg.Mount)
	}
	if !cfg.MaskValues {
		t.Error("expected MaskValues to default to true")
	}
	if cfg.ShowAll {
		t.Error("expected ShowAll to default to false")
	}
	if cfg.OutputJSON {
		t.Error("expected OutputJSON to default to false")
	}
}

func TestParseFlags_AllFlags(t *testing.T) {
	cfg, err := ParseFlags([]string{
		"-addr", "http://vault:8200",
		"-token", "s.abc123",
		"-path", "prod/db",
		"-mount", "kv",
		"-version-a", "3",
		"-version-b", "4",
		"-show-all",
		"-mask=false",
		"-audit-log", "/tmp/audit.log",
		"-json",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Mount != "kv" {
		t.Errorf("expected mount 'kv', got %q", cfg.Mount)
	}
	if cfg.VersionA != 3 {
		t.Errorf("expected VersionA=3, got %d", cfg.VersionA)
	}
	if cfg.VersionB != 4 {
		t.Errorf("expected VersionB=4, got %d", cfg.VersionB)
	}
	if !cfg.ShowAll {
		t.Error("expected ShowAll=true")
	}
	if cfg.MaskValues {
		t.Error("expected MaskValues=false")
	}
	if cfg.AuditLog != "/tmp/audit.log" {
		t.Errorf("expected audit-log '/tmp/audit.log', got %q", cfg.AuditLog)
	}
	if !cfg.OutputJSON {
		t.Error("expected OutputJSON=true")
	}
}
