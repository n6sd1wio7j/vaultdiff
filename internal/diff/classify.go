// Package diff provides utilities for comparing, rendering, and analyzing
// Vault secret versions. classify.go assigns risk classifications to diff
// entries based on key patterns and change status.
package diff

import (
	"fmt"
	"strings"
)

// RiskLevel represents the severity of a secret change.
type RiskLevel int

const (
	RiskNone     RiskLevel = iota // No risk (unchanged)
	RiskLow                       // Low risk (non-sensitive key changed)
	RiskMedium                    // Medium risk (potentially sensitive key)
	RiskHigh                      // High risk (known sensitive key pattern)
	RiskCritical                  // Critical risk (credential or token removed/modified)
)

// String returns a human-readable label for the risk level.
func (r RiskLevel) String() string {
	switch r {
	case RiskNone:
		return "none"
	case RiskLow:
		return "low"
	case RiskMedium:
		return "medium"
	case RiskHigh:
		return "high"
	case RiskCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// ClassifiedEntry pairs a DiffEntry with its computed risk level.
type ClassifiedEntry struct {
	Entry DiffEntry
	Risk  RiskLevel
	Reason string
}

// ClassifyOptions controls how entries are classified.
type ClassifyOptions struct {
	// ExtraHighRiskSuffixes are additional key suffixes treated as high risk.
	ExtraHighRiskSuffixes []string
	// SkipUnchanged omits entries with no change from the result.
	SkipUnchanged bool
}

// criticalSuffixes are key suffixes that indicate credential material.
var criticalSuffixes = []string{
	"_password", "_passwd", "_secret", "_token",
	"_private_key", "_api_key", "_access_key",
}

// highRiskSuffixes are key suffixes that indicate sensitive configuration.
var highRiskSuffixes = []string{
	"_key", "_cert", "_credential", "_credentials",
	"_dsn", "_url", "_connection_string",
}

// mediumRiskSuffixes are key suffixes that may carry sensitive data.
var mediumRiskSuffixes = []string{
	"_user", "_username", "_host", "_endpoint",
	"_id", "_name",
}

// Classify assigns a RiskLevel to each DiffEntry and returns the results.
// Unchanged entries are included unless SkipUnchanged is set.
func Classify(entries []DiffEntry, opts ClassifyOptions) []ClassifiedEntry {
	allCritical := criticalSuffixes
	allHigh := append(highRiskSuffixes, opts.ExtraHighRiskSuffixes...)

	var result []ClassifiedEntry
	for _, e := range entries {
		if opts.SkipUnchanged && e.Status == StatusUnchanged {
			continue
		}
		risk, reason := classifyEntry(e, allCritical, allHigh)
		result = append(result, ClassifiedEntry{
			Entry:  e,
			Risk:   risk,
			Reason: reason,
		})
	}
	return result
}

// classifyEntry determines the risk level for a single entry.
func classifyEntry(e DiffEntry, critical, high []string) (RiskLevel, string) {
	if e.Status == StatusUnchanged {
		return RiskNone, "no change"
	}

	lower := strings.ToLower(e.Key)

	for _, suf := range critical {
		if strings.HasSuffix(lower, suf) {
			return RiskCritical, fmt.Sprintf("key matches critical pattern %q", suf)
		}
	}

	for _, suf := range high {
		if strings.HasSuffix(lower, suf) {
			return RiskHigh, fmt.Sprintf("key matches high-risk pattern %q", suf)
		}
	}

	for _, suf := range mediumRiskSuffixes {
		if strings.HasSuffix(lower, suf) {
			return RiskMedium, fmt.Sprintf("key matches medium-risk pattern %q", suf)
		}
	}

	return RiskLow, "non-sensitive key changed"
}

// FormatClassified returns a human-readable summary of classified entries.
func FormatClassified(entries []ClassifiedEntry) string {
	if len(entries) == 0 {
		return "no classified entries\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s %-10s %s\n", "KEY", "RISK", "REASON"))
	sb.WriteString(strings.Repeat("-", 70) + "\n")
	for _, c := range entries {
		sb.WriteString(fmt.Sprintf("%-30s %-10s %s\n", c.Entry.Key, c.Risk.String(), c.Reason))
	}
	return sb.String()
}
