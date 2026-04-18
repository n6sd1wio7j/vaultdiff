package diff

import (
	"strings"
	"testing"
)

var sampleRollbackEntries = []Entry{
	{Key: "DB_HOST", Status: Modified, OldValue: "old-host", NewValue: "new-host"},
	{Key: "API_KEY", Status: Added, OldValue: "", NewValue: "abc123"},
	{Key: "LEGACY", Status: Removed, OldValue: "legacy-val", NewValue: ""},
	{Key: "PORT", Status: Unchanged, OldValue: "5432", NewValue: "5432"},
}

func TestBuildRollbackPlan_ExcludesUnchanged(t *testing.T) {
	plan := BuildRollbackPlan("secret/app", sampleRollbackEntries)
	if len(plan.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(plan.Entries))
	}
}

func TestBuildRollbackPlan_ModifiedReverts(t *testing.T) {
	plan := BuildRollbackPlan("secret/app", sampleRollbackEntries)
	var found *RollbackEntry
	for i := range plan.Entries {
		if plan.Entries[i].Key == "DB_HOST" {
			found = &plan.Entries[i]
		}
	}
	if found == nil {
		t.Fatal("expected DB_HOST in plan")
	}
	if found.From != "new-host" || found.To != "old-host" {
		t.Errorf("unexpected revert values: %+v", found)
	}
}

func TestBuildRollbackPlan_AddedBecomesDelete(t *testing.T) {
	plan := BuildRollbackPlan("secret/app", sampleRollbackEntries)
	for _, e := range plan.Entries {
		if e.Key == "API_KEY" {
			if e.To != "" {
				t.Errorf("expected To to be empty for added key, got %q", e.To)
			}
			return
		}
	}
	t.Fatal("API_KEY not found in plan")
}

func TestBuildRollbackPlan_RemovedBecomesRestore(t *testing.T) {
	plan := BuildRollbackPlan("secret/app", sampleRollbackEntries)
	for _, e := range plan.Entries {
		if e.Key == "LEGACY" {
			if e.To != "legacy-val" {
				t.Errorf("expected To=legacy-val, got %q", e.To)
			}
			return
		}
	}
	t.Fatal("LEGACY not found in plan")
}

func TestFormatRollbackPlan_NoChanges(t *testing.T) {
	plan := BuildRollbackPlan("secret/app", []Entry{})
	out := FormatRollbackPlan(plan)
	if !strings.Contains(out, "no changes") {
		t.Errorf("expected 'no changes' in output, got: %s", out)
	}
}

func TestFormatRollbackPlan_ContainsKeys(t *testing.T) {
	plan := BuildRollbackPlan("secret/app", sampleRollbackEntries)
	out := FormatRollbackPlan(plan)
	for _, key := range []string{"DB_HOST", "API_KEY", "LEGACY"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected key %q in output", key)
		}
	}
}
