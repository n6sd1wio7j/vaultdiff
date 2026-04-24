package diff

import (
	"strings"
	"testing"
)

func TestSchemaIntegration_FullPipeline(t *testing.T) {
	// Simulate a real diff where some keys violate schema rules
	old := map[string]string{
		"DATABASE_URL": "postgres://prod-db/app",
		"LOG_LEVEL":    "info",
		"PORT":         "5432",
	}
	new := map[string]string{
		"DATABASE_URL": "mysql://staging-db/app", // violates pattern
		"LOG_LEVEL":    "trace",                  // violates pattern
		"PORT":         "not-a-port",             // violates pattern
	}

	entries := Compare(old, new)
	rules := DefaultSchemaRules()
	violations := ValidateSchema(entries, rules)

	if len(violations) == 0 {
		t.Fatal("expected schema violations from integration pipeline")
	}

	out := FormatSchemaViolations(violations)
	if !strings.Contains(out, "violations") {
		t.Errorf("expected formatted output to mention violations, got: %q", out)
	}

	// Ensure DATABASE_URL violation is present
	found := false
	for _, v := range violations {
		if v.Key == "DATABASE_URL" {
			found = true
		}
	}
	if !found {
		t.Error("expected DATABASE_URL violation in integration test")
	}
}
