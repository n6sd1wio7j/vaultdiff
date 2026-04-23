package diff

import (
	"strings"
	"testing"
)

var sampleGraphSnapshots = map[int][]Entry{
	1: {
		{Key: "DB_HOST", Status: StatusAdded, NewValue: "localhost"},
		{Key: "DB_PASS", Status: StatusAdded, NewValue: "secret"},
	},
	2: {
		{Key: "DB_HOST", Status: StatusUnchanged, NewValue: "localhost"},
		{Key: "DB_PASS", Status: StatusModified, OldValue: "secret", NewValue: "newsecret"},
		{Key: "API_KEY", Status: StatusAdded, NewValue: "abc123"},
	},
	3: {
		{Key: "DB_HOST", Status: StatusRemoved, OldValue: "localhost"},
		{Key: "DB_PASS", Status: StatusUnchanged, NewValue: "newsecret"},
	},
}

func TestBuildGraph_KeysPresent(t *testing.T) {
	nodes := BuildGraph(sampleGraphSnapshots)
	keys := map[string]bool{}
	for _, n := range nodes {
		keys[n.Key] = true
	}
	for _, k := range []string{"DB_HOST", "DB_PASS", "API_KEY"} {
		if !keys[k] {
			t.Errorf("expected key %q in graph", k)
		}
	}
}

func TestBuildGraph_VersionOrder(t *testing.T) {
	nodes := BuildGraph(sampleGraphSnapshots)
	for _, n := range nodes {
		if n.Key == "DB_HOST" {
			if len(n.Versions) != 3 {
				t.Fatalf("expected 3 versions for DB_HOST, got %d", len(n.Versions))
			}
			if n.Versions[0].Version != 1 || n.Versions[1].Version != 2 || n.Versions[2].Version != 3 {
				t.Errorf("versions out of order: %+v", n.Versions)
			}
		}
	}
}

func TestBuildGraph_StatusCorrect(t *testing.T) {
	nodes := BuildGraph(sampleGraphSnapshots)
	for _, n := range nodes {
		if n.Key == "DB_PASS" {
			if n.Versions[1].Status != StatusModified {
				t.Errorf("expected modified at v2, got %s", n.Versions[1].Status)
			}
		}
	}
}

func TestBuildGraph_EmptySnapshots(t *testing.T) {
	nodes := BuildGraph(map[int][]Entry{})
	if len(nodes) != 0 {
		t.Errorf("expected empty graph, got %d nodes", len(nodes))
	}
}

func TestFormatGraph_ContainsKeys(t *testing.T) {
	nodes := BuildGraph(sampleGraphSnapshots)
	out := FormatGraph(nodes)
	for _, k := range []string{"DB_HOST", "DB_PASS", "API_KEY"} {
		if !strings.Contains(out, k) {
			t.Errorf("expected key %q in formatted graph", k)
		}
	}
}

func TestFormatGraph_ContainsSymbols(t *testing.T) {
	nodes := BuildGraph(sampleGraphSnapshots)
	out := FormatGraph(nodes)
	for _, sym := range []string{"+", "-", "~", "="} {
		if !strings.Contains(out, sym) {
			t.Errorf("expected symbol %q in formatted graph", sym)
		}
	}
}

func TestFormatGraph_Empty(t *testing.T) {
	out := FormatGraph([]GraphNode{})
	if !strings.Contains(out, "no graph data") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
