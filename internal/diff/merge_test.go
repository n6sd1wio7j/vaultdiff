package diff

import (
	"strings"
	"testing"
)

func TestMerge_NoConflicts(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	incoming := map[string]string{"C": "3"}
	r := Merge(base, incoming, MergeOptions{})
	if len(r.Conflicts) != 0 {
		t.Fatalf("expected no conflicts, got %v", r.Conflicts)
	}
	if r.Merged["C"] != "3" {
		t.Errorf("expected C=3, got %s", r.Merged["C"])
	}
}

func TestMerge_StrategyTheirs(t *testing.T) {
	base := map[string]string{"KEY": "old"}
	incoming := map[string]string{"KEY": "new"}
	r := Merge(base, incoming, MergeOptions{Strategy: StrategyTheirs})
	if r.Merged["KEY"] != "new" {
		t.Errorf("expected new, got %s", r.Merged["KEY"])
	}
	if len(r.Conflicts) != 1 || r.Conflicts[0] != "KEY" {
		t.Errorf("expected conflict on KEY")
	}
}

func TestMerge_StrategyOurs(t *testing.T) {
	base := map[string]string{"KEY": "old"}
	incoming := map[string]string{"KEY": "new"}
	r := Merge(base, incoming, MergeOptions{Strategy: StrategyOurs})
	if r.Merged["KEY"] != "old" {
		t.Errorf("expected old, got %s", r.Merged["KEY"])
	}
}

func TestMerge_StrategyUnion(t *testing.T) {
	base := map[string]string{"KEY": "a"}
	incoming := map[string]string{"KEY": "b"}
	r := Merge(base, incoming, MergeOptions{Strategy: StrategyUnion})
	if r.Merged["KEY"] != "a|b" {
		t.Errorf("expected a|b, got %s", r.Merged["KEY"])
	}
}

func TestMerge_IncomingFillsEmpty(t *testing.T) {
	base := map[string]string{"KEY": ""}
	incoming := map[string]string{"KEY": "filled"}
	r := Merge(base, incoming, MergeOptions{})
	if len(r.Conflicts) != 0 {
		t.Errorf("empty base value should not conflict")
	}
	if r.Merged["KEY"] != "filled" {
		t.Errorf("expected filled")
	}
}

func TestFormatMergeResult_NoConflicts(t *testing.T) {
	r := MergeResult{Merged: map[string]string{"A": "1"}}
	out := FormatMergeResult(r)
	if !strings.Contains(out, "no conflicts") {
		t.Errorf("expected no conflicts in output")
	}
}

func TestFormatMergeResult_WithConflicts(t *testing.T) {
	r := MergeResult{
		Merged:    map[string]string{"X": "v"},
		Conflicts: []string{"X"},
	}
	out := FormatMergeResult(r)
	if !strings.Contains(out, "conflict: X") {
		t.Errorf("expected conflict line for X")
	}
}
