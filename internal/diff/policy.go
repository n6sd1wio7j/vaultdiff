package diff

import (
	"fmt"
	"strings"
)

// PolicyRule defines a rule that flags certain diff entries.
type PolicyRule struct {
	Name     string
	KeySuffix string // e.g. "_key", "_secret"
	Status   string // "added", "removed", "modified", or "" for any
	Severity string // "warn" or "error"
}

// PolicyViolation represents a rule match against a diff entry.
type PolicyViolation struct {
	Rule    PolicyRule
	Entry   Entry
	Message string
}

// DefaultPolicyRules returns a baseline set of policy rules.
func DefaultPolicyRules() []PolicyRule {
	return []PolicyRule{
		{Name: "no-removed-secrets", KeySuffix: "_secret", Status: "removed", Severity: "error"},
		{Name: "no-removed-keys", KeySuffix: "_key", Status: "removed", Severity: "error"},
		{Name: "warn-modified-token", KeySuffix: "_token", Status: "modified", Severity: "warn"},
		{Name: "warn-added-password", KeySuffix: "_password", Status: "added", Severity: "warn"},
	}
}

// EnforcePolicy checks entries against rules and returns violations.
func EnforcePolicy(entries []Entry, rules []PolicyRule) []PolicyViolation {
	var violations []PolicyViolation
	for _, e := range entries {
		for _, r := range rules {
			if !strings.HasSuffix(strings.ToLower(e.Key), r.KeySuffix) {
				continue
			}
			if r.Status != "" && string(e.Status) != r.Status {
				continue
			}
			violations = append(violations, PolicyViolation{
				Rule:  r,
				Entry: e,
				Message: fmt.Sprintf("[%s] rule %q triggered on key %q", r.Severity, r.Name, e.Key),
			})
		}
	}
	return violations
}

// FormatViolations returns a human-readable summary of policy violations.
func FormatViolations(violations []PolicyViolation) string {
	if len(violations) == 0 {
		return "policy check passed: no violations found\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("policy violations (%d):\n", len(violations)))
	for _, v := range violations {
		sb.WriteString("  " + v.Message + "\n")
	}
	return sb.String()
}
