package diff

import (
	"strings"
	"testing"
)

func TestCheckThreshold_NoViolations(t *testing.T) {
	score := DriftScore{Added: 1, Removed: 1, Modified: 1, WeightedScore: 5.0}
	opts := DefaultThresholdOptions()
	res := CheckThreshold(score, opts)
	if res.Breached {
		t.Fatalf("expected no breach, got violations: %v", res.Violations)
	}
}

func TestCheckThreshold_ScoreBreached(t *testing.T) {
	score := DriftScore{WeightedScore: 200.0}
	opts := DefaultThresholdOptions()
	res := CheckThreshold(score, opts)
	if !res.Breached {
		t.Fatal("expected breach on weighted score")
	}
	if len(res.Violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(res.Violations))
	}
}

func TestCheckThreshold_AddedBreached(t *testing.T) {
	score := DriftScore{Added: 10}
	opts := DefaultThresholdOptions()
	opts.MaxAdded = 5
	res := CheckThreshold(score, opts)
	if !res.Breached {
		t.Fatal("expected breach on added keys")
	}
	if !strings.Contains(res.Violations[0], "added keys") {
		t.Errorf("unexpected violation message: %s", res.Violations[0])
	}
}

func TestCheckThreshold_RemovedBreached(t *testing.T) {
	score := DriftScore{Removed: 20}
	opts := DefaultThresholdOptions()
	opts.MaxRemoved = 3
	res := CheckThreshold(score, opts)
	if !res.Breached {
		t.Fatal("expected breach on removed keys")
	}
}

func TestCheckThreshold_ModifiedBreached(t *testing.T) {
	score := DriftScore{Modified: 7}
	opts := DefaultThresholdOptions()
	opts.MaxModified = 2
	res := CheckThreshold(score, opts)
	if !res.Breached {
		t.Fatal("expected breach on modified keys")
	}
}

func TestCheckThreshold_MultipleViolations(t *testing.T) {
	score := DriftScore{Added: 10, Removed: 10, Modified: 10, WeightedScore: 999}
	opts := ThresholdOptions{MaxScore: 1, MaxAdded: 1, MaxRemoved: 1, MaxModified: 1}
	res := CheckThreshold(score, opts)
	if len(res.Violations) != 4 {
		t.Fatalf("expected 4 violations, got %d", len(res.Violations))
	}
}

func TestFormatThresholdResult_Pass(t *testing.T) {
	res := ThresholdResult{Breached: false}
	out := FormatThresholdResult(res)
	if !strings.Contains(out, "passed") {
		t.Errorf("expected 'passed' in output, got: %s", out)
	}
}

func TestFormatThresholdResult_Fail(t *testing.T) {
	res := ThresholdResult{Breached: true, Violations: []string{"score too high"}}
	out := FormatThresholdResult(res)
	if !strings.Contains(out, "FAILED") {
		t.Errorf("expected 'FAILED' in output, got: %s", out)
	}
	if !strings.Contains(out, "score too high") {
		t.Errorf("expected violation detail in output, got: %s", out)
	}
}
