package diff

import (
	"strings"
	"testing"
)

func TestLintIntegration_FullPipeline(t *testing.T) {
	old := map[string]string{
		"DB_PASSWORD": "hunter2",
		"API_KEY":     "abc",
	}
	new := map[string]string{
		"API_KEY":  "",
		"api_host": "localhost",
	}

	entries := Compare(old, new)
	violations := Lint(entries, DefaultLintRules())

	rulesSeen := map[string]bool{}
	for _, v := range violations {
		rulesSeen[v.Rule] = true
	}

	if !rulesSeen["removed-required"] {
		t.Error("expected removed-required violation")
	}
	if !rulesSeen["empty-value"] {
		t.Error("expected empty-value violation")
	}
	if !rulesSeen["key-uppercase"] {
		t.Error("expected key-uppercase violation")
	}

	out := FormatLint(violations)
	if !strings.Contains(out, "violation") {
		t.Errorf("expected violations in output, got: %s", out)
	}
}
