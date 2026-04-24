package diff

import (
	"strings"
	"testing"
)

var sampleSchemaEntries = []Entry{
	{Key: "DATABASE_URL", Status: StatusAdded, NewValue: "postgres://localhost/db"},
	{Key: "PORT", Status: StatusAdded, NewValue: "8080"},
	{Key: "LOG_LEVEL", Status: StatusAdded, NewValue: "info"},
	{Key: "APP_NAME", Status: StatusAdded, NewValue: "vaultdiff"},
}

func TestValidateSchema_NoViolations(t *testing.T) {
	rules := DefaultSchemaRules()
	violations := ValidateSchema(sampleSchemaEntries, rules)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d: %+v", len(violations), violations)
	}
}

func TestValidateSchema_BadDatabaseURL(t *testing.T) {
	entries := []Entry{
		{Key: "DATABASE_URL", Status: StatusAdded, NewValue: "mysql://localhost/db"},
	}
	rules := DefaultSchemaRules()
	violations := ValidateSchema(entries, rules)
	if len(violations) == 0 {
		t.Fatal("expected a pattern violation for DATABASE_URL")
	}
	if violations[0].Rule != "pattern" {
		t.Errorf("expected rule=pattern, got %q", violations[0].Rule)
	}
}

func TestValidateSchema_RequiredMissing(t *testing.T) {
	entries := []Entry{
		{Key: "PORT", Status: StatusAdded, NewValue: "9090"},
	}
	rules := DefaultSchemaRules()
	violations := ValidateSchema(entries, rules)
	found := false
	for _, v := range violations {
		if v.Key == "DATABASE_URL" && v.Rule == "required" {
			found = true
		}
	}
	if !found {
		t.Error("expected required violation for DATABASE_URL")
	}
}

func TestValidateSchema_RemovedKeySkipped(t *testing.T) {
	entries := []Entry{
		{Key: "DATABASE_URL", Status: StatusRemoved, OldValue: "postgres://localhost/db"},
	}
	rules := []SchemaRule{
		{Key: "DATABASE_URL", Pattern: `^postgres://`, Required: false},
	}
	violations := ValidateSchema(entries, rules)
	if len(violations) != 0 {
		t.Errorf("removed key should not trigger pattern violation, got %+v", violations)
	}
}

func TestValidateSchema_InvalidLogLevel(t *testing.T) {
	entries := []Entry{
		{Key: "DATABASE_URL", Status: StatusAdded, NewValue: "postgres://localhost/db"},
		{Key: "LOG_LEVEL", Status: StatusAdded, NewValue: "verbose"},
	}
	rules := DefaultSchemaRules()
	violations := ValidateSchema(entries, rules)
	found := false
	for _, v := range violations {
		if v.Key == "LOG_LEVEL" {
			found = true
		}
	}
	if !found {
		t.Error("expected violation for LOG_LEVEL with invalid value")
	}
}

func TestFormatSchemaViolations_NoViolations(t *testing.T) {
	out := FormatSchemaViolations(nil)
	if !strings.Contains(out, "no violations") {
		t.Errorf("expected 'no violations' in output, got %q", out)
	}
}

func TestFormatSchemaViolations_WithViolations(t *testing.T) {
	v := []SchemaViolation{
		{Key: "DATABASE_URL", Rule: "pattern", Message: "does not match"},
	}
	out := FormatSchemaViolations(v)
	if !strings.Contains(out, "DATABASE_URL") {
		t.Errorf("expected key in output, got %q", out)
	}
	if !strings.Contains(out, "pattern") {
		t.Errorf("expected rule in output, got %q", out)
	}
}
