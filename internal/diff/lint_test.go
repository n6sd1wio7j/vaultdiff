package diff

import (
	"strings"
	"testing"
)

var sampleLintEntries = []Entry{
	{Key: "DB_PASSWORD", Status: Removed, OldValue: "secret", NewValue: ""},
	{Key: "api_token", Status: Added, OldValue: "", NewValue: "abc123"},
	{Key: "EMPTY_KEY", Status: Modified, OldValue: "old", NewValue: ""},
	{Key: "NORMAL_KEY", Status: Unchanged, OldValue: "val", NewValue: "val"},
}

func TestLint_DetectsEmptyValue(t *testing.T) {
	rules := DefaultLintRules()
	vs := Lint(sampleLintEntries, rules)
	for _, v := range vs {
		if v.Key == "EMPTY_KEY" && v.Rule == "empty-value" {
			return
		}
	}
	t.Fatal("expected empty-value violation for EMPTY_KEY")
}

func TestLint_DetectsLowercaseKey(t *testing.T) {
	rules := DefaultLintRules()
	vs := Lint(sampleLintEntries, rules)
	for _, v := range vs {
		if v.Key == "api_token" && v.Rule == "key-uppercase" {
			return
		}
	}
	t.Fatal("expected key-uppercase violation for api_token")
}

func TestLint_DetectsRemovedRequired(t *testing.T) {
	rules := DefaultLintRules()
	vs := Lint(sampleLintEntries, rules)
	for _, v := range vs {
		if v.Key == "DB_PASSWORD" && v.Rule == "removed-required" {
			return
		}
	}
	t.Fatal("expected removed-required violation for DB_PASSWORD")
}

func TestLint_NoViolationsForClean(t *testing.T) {
	clean := []Entry{
		{Key: "MY_KEY", Status: Added, NewValue: "value"},
	}
	vs := Lint(clean, DefaultLintRules())
	if len(vs) != 0 {
		t.Fatalf("expected no violations, got %d", len(vs))
	}
}

func TestFormatLint_NoViolations(t *testing.T) {
	out := FormatLint(nil)
	if !strings.Contains(out, "no violations") {
		t.Fatalf("expected no violations message, got: %s", out)
	}
}

func TestFormatLint_ShowsViolations(t *testing.T) {
	vs := []LintViolation{
		{Key: "FOO", Rule: "empty-value", Message: "secret value is empty"},
	}
	out := FormatLint(vs)
	if !strings.Contains(out, "FOO") || !strings.Contains(out, "empty-value") {
		t.Fatalf("unexpected output: %s", out)
	}
}
