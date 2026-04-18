package diff

import "fmt"

// RollbackEntry represents a single key reversion to a prior value.
type RollbackEntry struct {
	Key      string `json:"key"`
	From     string `json:"from"`
	To       string `json:"to"`
	Reverted bool   `json:"reverted"`
}

// RollbackPlan holds the set of changes needed to revert to a baseline.
type RollbackPlan struct {
	Path    string          `json:"path"`
	Entries []RollbackEntry `json:"entries"`
}

// BuildRollbackPlan computes a rollback plan from current diff entries,
// producing reversions for any added, removed, or modified keys.
func BuildRollbackPlan(path string, entries []Entry) RollbackPlan {
	plan := RollbackPlan{Path: path}
	for _, e := range entries {
		switch e.Status {
		case Added:
			plan.Entries = append(plan.Entries, RollbackEntry{
				Key:  e.Key,
				From: e.NewValue,
				To:   "",
			})
		case Removed:
			plan.Entries = append(plan.Entries, RollbackEntry{
				Key:  e.Key,
				From: "",
				To:   e.OldValue,
			})
		case Modified:
			plan.Entries = append(plan.Entries, RollbackEntry{
				Key:  e.Key,
				From: e.NewValue,
				To:   e.OldValue,
			})
		}
	}
	return plan
}

// FormatRollbackPlan returns a human-readable summary of the rollback plan.
func FormatRollbackPlan(plan RollbackPlan) string {
	if len(plan.Entries) == 0 {
		return fmt.Sprintf("rollback plan for %s: no changes to revert\n", plan.Path)
	}
	out := fmt.Sprintf("rollback plan for %s (%d reversion(s)):\n", plan.Path, len(plan.Entries))
	for _, e := range plan.Entries {
		switch {
		case e.From != "" && e.To == "":
			out += fmt.Sprintf("  - DELETE %s (was %q)\n", e.Key, e.From)
		case e.From == "" && e.To != "":
			out += fmt.Sprintf("  + RESTORE %s => %q\n", e.Key, e.To)
		default:
			out += fmt.Sprintf("  ~ REVERT %s: %q => %q\n", e.Key, e.From, e.To)
		}
	}
	return out
}
