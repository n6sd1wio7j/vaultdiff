package diff

import (
	"strings"
	"testing"
)

var v1 = map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
var v2 = map[string]string{"DB_HOST": "prod-db", "DB_PORT": "5432"}
var v3 = map[string]string{"DB_HOST": "prod-db", "DB_PORT": "5432", "DB_NAME": "app"}

func TestBuildChain_ReturnsEmptyForSingleVersion(t *testing.T) {
	c := BuildChain("secret/db", []map[string]string{v1})
	if len(c.Steps) != 0 {
		t.Errorf("expected 0 steps, got %d", len(c.Steps))
	}
}

func TestBuildChain_StepCountMatchesVersionGaps(t *testing.T) {
	c := BuildChain("secret/db", []map[string]string{v1, v2, v3})
	if len(c.Steps) != 2 {
		t.Errorf("expected 2 steps, got %d", len(c.Steps))
	}
}

func TestBuildChain_DetectsModifiedInStep(t *testing.T) {
	c := BuildChain("secret/db", []map[string]string{v1, v2})
	if len(c.Steps) == 0 {
		t.Fatal("expected at least one step")
	}
	var found bool
	for _, e := range c.Steps[0].Entries {
		if e.Key == "DB_HOST" && e.Status == StatusModified {
			found = true
		}
	}
	if !found {
		t.Error("expected DB_HOST to be modified in step 1")
	}
}

func TestBuildChain_DetectsAddedInStep(t *testing.T) {
	c := BuildChain("secret/db", []map[string]string{v2, v3})
	var found bool
	for _, e := range c.Steps[0].Entries {
		if e.Key == "DB_NAME" && e.Status == StatusAdded {
			found = true
		}
	}
	if !found {
		t.Error("expected DB_NAME to be added")
	}
}

func TestHasAnyDrift_TrueWhenChangesExist(t *testing.T) {
	c := BuildChain("secret/db", []map[string]string{v1, v2})
	if !c.HasAnyDrift() {
		t.Error("expected HasAnyDrift to be true")
	}
}

func TestHasAnyDrift_FalseWhenNoChanges(t *testing.T) {
	c := BuildChain("secret/db", []map[string]string{v1, v1})
	if c.HasAnyDrift() {
		t.Error("expected HasAnyDrift to be false")
	}
}

func TestFormatChain_ContainsPath(t *testing.T) {
	c := BuildChain("secret/db", []map[string]string{v1, v2})
	out := FormatChain(c)
	if !strings.Contains(out, "secret/db") {
		t.Errorf("expected path in output, got: %s", out)
	}
}

func TestFormatChain_EmptySteps(t *testing.T) {
	c := Chain{Path: "secret/empty"}
	out := FormatChain(c)
	if !strings.Contains(out, "no steps") {
		t.Errorf("expected 'no steps' in output, got: %s", out)
	}
}

func TestFormatChain_NoChangesLabel(t *testing.T) {
	c := BuildChain("secret/db", []map[string]string{v1, v1})
	out := FormatChain(c)
	if !strings.Contains(out, "no changes") {
		t.Errorf("expected 'no changes' label, got: %s", out)
	}
}
