package diff

import (
	"fmt"
	"strings"
)

// LintRule defines a rule applied to diff entries.
type LintRule struct {
	Name    string
	Message string
	Check   func(e Entry) bool
}

// LintViolation is a rule violation for a specific entry.
type LintViolation struct {
	Key     string
	Rule    string
	Message string
}

// DefaultLintRules returns a standard set of lint rules.
func DefaultLintRules() []LintRule {
	return []LintRule{
		{
			Name:    "empty-value",
			Message: "secret value is empty",
			Check: func(e Entry) bool {
				return (e.Status == Added || e.Status == Modified) && strings.TrimSpace(e.NewValue) == ""
			},
		},
		{
			Name:    "key-uppercase",
			Message: "key should be uppercase",
			Check: func(e Entry) bool {
				return e.Key != strings.ToUpper(e.Key)
			},
		},
		{
			Name:    "removed-required",
			Message: "required key was removed",
			Check: func(e Entry) bool {
				required := []string{"TOKEN", "PASSWORD", "SECRET", "KEY"}
				if e.Status != Removed {
					return false
				}
				for _, r := range required {
					if strings.Contains(strings.ToUpper(e.Key), r) {
						return true
					}
				}
				return false
			},
		},
	}
}

// Lint runs lint rules against entries and returns violations.
func Lint(entries []Entry, rules []LintRule) []LintViolation {
	var violations []LintViolation
	for _, e := range entries {
		for _, r := range rules {
			if r.Check(e) {
				violations = append(violations, LintViolation{
					Key:     e.Key,
					Rule:    r.Name,
					Message: r.Message,
				})
			}
		}
	}
	return violations
}

// FormatLint returns a human-readable lint report.
func FormatLint(violations []LintViolation) string {
	if len(violations) == 0 {
		return "lint: no violations found\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("lint: %d violation(s)\n", len(violations)))
	for _, v := range violations {
		sb.WriteString(fmt.Sprintf("  [%s] %s — %s\n", v.Rule, v.Key, v.Message))
	}
	return sb.String()
}

// ViolationsByRule groups lint violations by rule name, returning a map
// of rule name to the list of violations that matched that rule.
func ViolationsByRule(violations []LintViolation) map[string][]LintViolation {
	result := make(map[string][]LintViolation)
	for _, v := range violations {
		result[v.Rule] = append(result[v.Rule], v)
	}
	return result
}
