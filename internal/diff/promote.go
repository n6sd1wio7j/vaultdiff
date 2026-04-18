package diff

import (
	"fmt"
	"strings"
)

// PromoteAction describes what should happen to a key during promotion.
type PromoteAction struct {
	Key    string
	Value  string
	Action string // "set" or "delete"
}

// PromotePlan holds the list of actions to promote changes from source to target.
type PromotePlan struct {
	SourcePath string
	TargetPath string
	Actions    []PromoteAction
}

// BuildPromotePlan creates a promotion plan from diff entries.
// Only added, removed, and modified entries are included.
func BuildPromotePlan(sourcePath, targetPath string, entries []Entry) PromotePlan {
	plan := PromotePlan{
		SourcePath: sourcePath,
		TargetPath: targetPath,
	}
	for _, e := range entries {
		switch e.Status {
		case StatusAdded, StatusModified:
			plan.Actions = append(plan.Actions, PromoteAction{
				Key:    e.Key,
				Value:  e.NewValue,
				Action: "set",
			})
		case StatusRemoved:
			plan.Actions = append(plan.Actions, PromoteAction{
				Key:    e.Key,
				Action: "delete",
			})
		}
	}
	return plan
}

// FormatPromotePlan returns a human-readable summary of the promotion plan.
func FormatPromotePlan(plan PromotePlan) string {
	if len(plan.Actions) == 0 {
		return fmt.Sprintf("No changes to promote from %s to %s\n", plan.SourcePath, plan.TargetPath)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Promote: %s → %s\n", plan.SourcePath, plan.TargetPath)
	fmt.Fprintf(&sb, "Actions (%d):\n", len(plan.Actions))
	for _, a := range plan.Actions {
		if a.Action == "delete" {
			fmt.Fprintf(&sb, "  DELETE  %s\n", a.Key)
		} else {
			fmt.Fprintf(&sb, "  SET     %s\n", a.Key)
		}
	}
	return sb.String()
}
