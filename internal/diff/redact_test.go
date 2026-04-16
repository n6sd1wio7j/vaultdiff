package diff

import (
	"testing"
)

func TestShouldRedact_ExactMatch(t *testing.T) {
	opts := DefaultRedactOptions()
	if !ShouldRedact("password", opts) {
		t.Error("expected 'password' to be redacted")
	}
}

func TestShouldRedact_SuffixMatch(t *testing.T) {
	opts := DefaultRedactOptions()
	for _, key := range []string{"db_password", "aws_secret", "auth_token", "private_key"} {
		if !ShouldRedact(key, opts) {
			t.Errorf("expected key %q to be redacted", key)
		}
	}
}

func TestShouldRedact_SafeKey(t *testing.T) {
	opts := DefaultRedactOptions()
	if ShouldRedact("database_host", opts) {
		t.Error("expected 'database_host' not to be redacted")
	}
}

func TestShouldRedact_CaseInsensitive(t *testing.T) {
	opts := DefaultRedactOptions()
	if !ShouldRedact("DB_PASSWORD", opts) {
		t.Error("expected 'DB_PASSWORD' to be redacted")
	}
}

func TestRedact_MasksValues(t *testing.T) {
	entries := []DiffEntry{
		{Key: "db_password", OldValue: "hunter2", NewValue: "secret123", Status: StatusModified},
		{Key: "host", OldValue: "localhost", NewValue: "prod.db", Status: StatusModified},
	}
	opts := DefaultRedactOptions()
	result := Redact(entries, opts)

	if result[0].OldValue != "***" || result[0].NewValue != "***" {
		t.Errorf("expected db_password values to be masked, got old=%q new=%q", result[0].OldValue, result[0].NewValue)
	}
	if result[1].OldValue != "localhost" || result[1].NewValue != "prod.db" {
		t.Error("expected host values to remain unmasked")
	}
}

func TestRedact_EmptyValueNotOverwritten(t *testing.T) {
	entries := []DiffEntry{
		{Key: "api_token", OldValue: "", NewValue: "newtoken", Status: StatusAdded},
	}
	opts := DefaultRedactOptions()
	result := Redact(entries, opts)
	if result[0].OldValue != "" {
		t.Error("expected empty OldValue to remain empty")
	}
	if result[0].NewValue != "***" {
		t.Error("expected NewValue to be masked")
	}
}

func TestRedact_DoesNotMutateOriginal(t *testing.T) {
	original := []DiffEntry{
		{Key: "secret", OldValue: "original", NewValue: "updated", Status: StatusModified},
	}
	opts := DefaultRedactOptions()
	Redact(original, opts)
	if original[0].OldValue != "original" {
		t.Error("expected original slice to be unmodified")
	}
}
