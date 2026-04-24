package diff

import "fmt"

// ThresholdOptions controls when a drift score triggers an alert.
type ThresholdOptions struct {
	MaxScore     float64 // maximum allowed weighted score before breach
	MaxAdded     int     // maximum allowed added keys
	MaxRemoved   int     // maximum allowed removed keys
	MaxModified  int     // maximum allowed modified keys
}

// ThresholdResult captures whether any threshold was breached.
type ThresholdResult struct {
	Breached   bool
	Violations []string
}

// DefaultThresholdOptions returns permissive defaults.
func DefaultThresholdOptions() ThresholdOptions {
	return ThresholdOptions{
		MaxScore:    100.0,
		MaxAdded:    50,
		MaxRemoved:  50,
		MaxModified: 50,
	}
}

// CheckThreshold evaluates a DriftScore against the given options.
func CheckThreshold(score DriftScore, opts ThresholdOptions) ThresholdResult {
	var violations []string

	if score.WeightedScore > opts.MaxScore {
		violations = append(violations,
			fmt.Sprintf("weighted score %.2f exceeds max %.2f", score.WeightedScore, opts.MaxScore))
	}
	if score.Added > opts.MaxAdded {
		violations = append(violations,
			fmt.Sprintf("added keys %d exceeds max %d", score.Added, opts.MaxAdded))
	}
	if score.Removed > opts.MaxRemoved {
		violations = append(violations,
			fmt.Sprintf("removed keys %d exceeds max %d", score.Removed, opts.MaxRemoved))
	}
	if score.Modified > opts.MaxModified {
		violations = append(violations,
			fmt.Sprintf("modified keys %d exceeds max %d", score.Modified, opts.MaxModified))
	}

	return ThresholdResult{
		Breached:   len(violations) > 0,
		Violations: violations,
	}
}

// FormatThresholdResult returns a human-readable summary of the result.
func FormatThresholdResult(r ThresholdResult) string {
	if !r.Breached {
		return "threshold check passed: no violations\n"
	}
	out := fmt.Sprintf("threshold check FAILED (%d violation(s)):\n", len(r.Violations))
	for _, v := range r.Violations {
		out += fmt.Sprintf("  ! %s\n", v)
	}
	return out
}
