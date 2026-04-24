package main

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/internal/diff"
)

func TestEnvDiffFlags_MissingPath(t *testing.T) {
	err := RunEnvDiff([]string{
		"--addr-a", "http://vault-a:8200",
		"--token-a", "tok-a",
		"--addr-b", "http://vault-b:8200",
		"--token-b", "tok-b",
	})
	if err == nil || !strings.Contains(err.Error(), "--path") {
		t.Errorf("expected --path error, got: %v", err)
	}
}

func TestEnvDiffFlags_MissingAddrA(t *testing.T) {
	err := RunEnvDiff([]string{
		"--path", "myapp/config",
		"--token-a", "tok-a",
		"--addr-b", "http://vault-b:8200",
		"--token-b", "tok-b",
	})
	if err == nil || !strings.Contains(err.Error(), "addr-a") {
		t.Errorf("expected addr-a error, got: %v", err)
	}
}

func TestEnvDiffFlags_MissingAddrB(t *testing.T) {
	err := RunEnvDiff([]string{
		"--path", "myapp/config",
		"--addr-a", "http://vault-a:8200",
		"--token-a", "tok-a",
		"--token-b", "tok-b",
	})
	if err == nil || !strings.Contains(err.Error(), "addr-b") {
		t.Errorf("expected addr-b error, got: %v", err)
	}
}

func TestCompareEnvs_Unit_Staging_Prod(t *testing.T) {
	a := map[string]string{"KEY": "val-staging", "SHARED": "same"}
	b := map[string]string{"KEY": "val-prod", "SHARED": "same"}

	opts := diff.EnvCompareOptions{EnvA: "staging", EnvB: "production"}
	d := diff.CompareEnvs(a, b, opts)

	if d.EnvA != "staging" {
		t.Errorf("unexpected EnvA: %s", d.EnvA)
	}
	if len(d.Entries) == 0 {
		t.Error("expected entries")
	}
}

func TestCompareEnvs_Unit_IgnorePrefix(t *testing.T) {
	a := map[string]string{"SECRET_KEY": "old", "APP": "v1"}
	b := map[string]string{"SECRET_KEY": "new", "APP": "v2"}

	opts := diff.EnvCompareOptions{
		EnvA:   "a",
		EnvB:   "b",
		Ignore: []string{"SECRET_"},
	}
	d := diff.CompareEnvs(a, b, opts)

	for _, e := range d.Entries {
		if strings.HasPrefix(e.Key, "SECRET_") {
			t.Errorf("SECRET_ key should have been ignored: %s", e.Key)
		}
	}
}
