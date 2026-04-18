package diff

import (
	"strings"
	"testing"
)

var samplePromoteEntries = []Entry{
	{Key: "DB_HOST", OldValue: "localhost", NewValue: "prod-db", Status: StatusModified},
	{Key: "NEW_KEY", OldValue: "", NewValue: "value1", Status: StatusAdded},
	{Key: "OLD_KEY", OldValue: "old", NewValue: "", Status: StatusRemoved},
	{Key: "STABLE", OldValue: "same", NewValue: "same", Status: StatusUnchanged},
}

func TestBuildPromotePlan_ExcludesUnchanged(t *testing.T) {
	plan := BuildPromotePlan("secret/dev", "secret/prod", samplePromoteEntries)
	for _, a := range plan.Actions {
		if a.Key == "STABLE" {
			t.Error("unchanged key should not appear in promote plan")
		}
	}
}

func TestBuildPromotePlan_ModifiedIsSet(t *testing.T) {
	plan := BuildPromotePlan("secret/dev", "secret/prod", samplePromoteEntries)
	for _, a := range plan.Actions {
		if a.Key == "DB_HOST" && a.Action != "set" {
			t.Errorf("expected set for modified key, got %s", a.Action)
		}
	}
}

func TestBuildPromotePlan_AddedIsSet(t *testing.T) {
	plan := BuildPromotePlan("secret/dev", "secret/prod", samplePromoteEntries)
	for _, a := range plan.Actions {
		if a.Key == "NEW_KEY" {
			if a.Action != "set" {
				t.Errorf("expected set for added key, got %s", a.Action)
			}
			if a.Value != "value1" {
				t.Errorf("expected value1, got %s", a.Value)
			}
		}
	}
}

func TestBuildPromotePlan_RemovedIsDelete(t *testing.T) {
	plan := BuildPromotePlan("secret/dev", "secret/prod", samplePromoteEntries)
	for _, a := range plan.Actions {
		if a.Key == "OLD_KEY" && a.Action != "delete" {
			t.Errorf("expected delete for removed key, got %s", a.Action)
		}
	}
}

func TestFormatPromotePlan_NoChanges(t *testing.T) {
	plan := BuildPromotePlan("secret/dev", "secret/prod", []Entry{})
	out := FormatPromotePlan(plan)
	if !strings.Contains(out, "No changes") {
		t.Errorf("expected no-changes message, got: %s", out)
	}
}

func TestFormatPromotePlan_ContainsPaths(t *testing.T) {
	plan := BuildPromotePlan("secret/dev", "secret/prod", samplePromoteEntries)
	out := FormatPromotePlan(plan)
	if !strings.Contains(out, "secret/dev") || !strings.Contains(out, "secret/prod") {
		t.Errorf("expected paths in output, got: %s", out)
	}
}
