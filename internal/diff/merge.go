package diff

import "fmt"

// MergeStrategy defines how conflicts are resolved during a merge.
type MergeStrategy string

const (
	StrategyOurs   MergeStrategy = "ours"
	StrategyTheirs MergeStrategy = "theirs"
	StrategyUnion  MergeStrategy = "union"
)

// MergeOptions controls merge behaviour.
type MergeOptions struct {
	Strategy MergeStrategy
}

// MergeResult holds the merged key-value map and any conflicts.
type MergeResult struct {
	Merged    map[string]string
	Conflicts []string
}

// Merge combines two secret maps (base and incoming) according to the given strategy.
// Conflicts arise when both sides have different non-empty values for the same key.
func Merge(base, incoming map[string]string, opts MergeOptions) MergeResult {
	if opts.Strategy == "" {
		opts.Strategy = StrategyTheirs
	}

	merged := make(map[string]string)
	var conflicts []string

	for k, v := range base {
		merged[k] = v
	}

	for k, inVal := range incoming {
		baseVal, exists := merged[k]
		if !exists || baseVal == "" {
			merged[k] = inVal
			continue
		}
		if baseVal == inVal {
			continue
		}
		// conflict
		conflicts = append(conflicts, k)
		switch opts.Strategy {
		case StrategyOurs:
			// keep base value — already set
		case StrategyTheirs:
			merged[k] = inVal
		case StrategyUnion:
			merged[k] = baseVal + "|" + inVal
		}
	}

	return MergeResult{Merged: merged, Conflicts: conflicts}
}

// FormatMergeResult returns a human-readable summary of the merge.
func FormatMergeResult(r MergeResult) string {
	if len(r.Conflicts) == 0 {
		return fmt.Sprintf("merge complete: %d keys, no conflicts", len(r.Merged))
	}
	out := fmt.Sprintf("merge complete: %d keys, %d conflict(s)\n", len(r.Merged), len(r.Conflicts))
	for _, k := range r.Conflicts {
		out += fmt.Sprintf("  conflict: %s\n", k)
	}
	return out
}
