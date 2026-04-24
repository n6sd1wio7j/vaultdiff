package diff

import (
	"fmt"
	"regexp"
	"strings"
)

// SchemaRule defines a validation rule applied to secret keys and values.
type SchemaRule struct {
	Key     string // exact key match; empty means apply to all
	Pattern string // regex the value must match
	Required bool   // key must be present with a non-empty value
}

// SchemaViolation describes a rule that was not satisfied.
type SchemaViolation struct {
	Key     string
	Rule    string
	Message string
}

// DefaultSchemaRules returns a baseline set of schema rules.
func DefaultSchemaRules() []SchemaRule {
	return []SchemaRule{
		{Key: "DATABASE_URL", Pattern: `^postgres://`, Required: true},
		{Key: "PORT", Pattern: `^[0-9]+$`},
		{Key: "LOG_LEVEL", Pattern: `^(debug|info|warn|error)$`},
	}
}

// ValidateSchema checks diff entries against the provided rules.
func ValidateSchema(entries []Entry, rules []SchemaRule) []SchemaViolation {
	var violations []SchemaViolation

	// Build a map of key -> current value from entries
	present := make(map[string]string)
	for _, e := range entries {
		if e.Status != StatusRemoved {
			present[e.Key] = e.NewValue
		}
	}

	for _, rule := range rules {
		if rule.Key != "" {
			val, exists := present[rule.Key]
			if rule.Required && (!exists || val == "") {
				violations = append(violations, SchemaViolation{
					Key:     rule.Key,
					Rule:    "required",
					Message: fmt.Sprintf("key %q is required but missing or empty", rule.Key),
				})
				continue
			}
			if exists && rule.Pattern != "" {
				if matched, _ := regexp.MatchString(rule.Pattern, val); !matched {
					violations = append(violations, SchemaViolation{
						Key:     rule.Key,
						Rule:    "pattern",
						Message: fmt.Sprintf("key %q value does not match pattern %q", rule.Key, rule.Pattern),
					})
				}
			}
		}
	}
	return violations
}

// FormatSchemaViolations renders violations as a human-readable string.
func FormatSchemaViolations(violations []SchemaViolation) string {
	if len(violations) == 0 {
		return "schema: no violations found\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "schema violations (%d):\n", len(violations))
	for _, v := range violations {
		fmt.Fprintf(&sb, "  [%s] %s: %s\n", v.Rule, v.Key, v.Message)
	}
	return sb.String()
}
