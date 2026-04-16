package diff

import (
	"strings"
	"testing"
)

func samplePatchEntries() []Entry {
	return []Entry{
		{Key: "HOST", Status: StatusUnchanged, OldValue: "localhost", NewValue: "localhost"},
		{Key: "PORT", Status: StatusModified, OldValue: "8080", NewValue: "9090"},
		{Key: "TOKEN", Status: StatusAdded, OldValue: "", NewValue: "abc123"},
		{Key: "LEGACY", Status: StatusRemoved, OldValue: "old", NewValue: ""},
	}
}

func TestPatch_UnchangedOp(t *testing.T) {
	p := Patch(samplePatchEntries())
	if p[0].Op != " " || p[0].Key != "HOST" {
		t.Errorf("expected unchanged HOST, got %+v", p[0])
	}
}

func TestPatch_ModifiedProducesTwoLines(t *testing.T) {
	p := Patch(samplePatchEntries())
	// HOST(0), PORT-(1), PORT+(2), TOKEN+(3), LEGACY-(4)
	if p[1].Op != "-" || p[1].Key != "PORT" {
		t.Errorf("expected removal line for PORT, got %+v", p[1])
	}
	if p[2].Op != "+" || p[2].Key != "PORT" {
		t.Errorf("expected addition line for PORT, got %+v", p[2])
	}
}

func TestPatch_AddedOp(t *testing.T) {
	p := Patch(samplePatchEntries())
	var found *PatchEntry
	for i := range p {
		if p[i].Key == "TOKEN" {
			found = &p[i]
			break
		}
	}
	if found == nil || found.Op != "+" {
		t.Error("expected + op for TOKEN")
	}
}

func TestPatch_RemovedOp(t *testing.T) {
	p := Patch(samplePatchEntries())
	var found *PatchEntry
	for i := range p {
		if p[i].Key == "LEGACY" {
			found = &p[i]
			break
		}
	}
	if found == nil || found.Op != "-" {
		t.Error("expected - op for LEGACY")
	}
}

func TestFormatPatch_ContainsKeys(t *testing.T) {
	p := Patch(samplePatchEntries())
	out := FormatPatch(p, false)
	for _, key := range []string{"HOST", "PORT", "TOKEN", "LEGACY"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected key %s in patch output", key)
		}
	}
}

func TestFormatPatch_MasksValues(t *testing.T) {
	p := Patch(samplePatchEntries())
	out := FormatPatch(p, true)
	if strings.Contains(out, "abc123") {
		t.Error("expected secret value to be masked")
	}
	if !strings.Contains(out, "***") {
		t.Error("expected masked placeholder in output")
	}
}
