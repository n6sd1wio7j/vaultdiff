package diff

import (
	"fmt"
	"sort"
	"strings"
)

// GraphNode represents a single key's change history across versions.
type GraphNode struct {
	Key      string
	Versions []GraphPoint
}

// GraphPoint represents the status of a key at a specific version.
type GraphPoint struct {
	Version int
	Status  ChangeStatus
	Value   string
}

// ChangeStatus mirrors the entry status type used across the diff package.
type ChangeStatus = string

const (
	StatusAdded     ChangeStatus = "added"
	StatusRemoved   ChangeStatus = "removed"
	StatusModified  ChangeStatus = "modified"
	StatusUnchanged ChangeStatus = "unchanged"
)

// BuildGraph constructs a graph of key changes across multiple diff snapshots.
// Each snapshot is associated with a version number.
func BuildGraph(snapshots map[int][]Entry) []GraphNode {
	keyMap := map[string]*GraphNode{}

	versions := make([]int, 0, len(snapshots))
	for v := range snapshots {
		versions = append(versions, v)
	}
	sort.Ints(versions)

	for _, v := range versions {
		entries := snapshots[v]
		for _, e := range entries {
			node, ok := keyMap[e.Key]
			if !ok {
				node = &GraphNode{Key: e.Key}
				keyMap[e.Key] = node
			}
			node.Versions = append(node.Versions, GraphPoint{
				Version: v,
				Status:  e.Status,
				Value:   e.NewValue,
			})
		}
	}

	nodes := make([]GraphNode, 0, len(keyMap))
	for _, n := range keyMap {
		nodes = append(nodes, *n)
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Key < nodes[j].Key
	})
	return nodes
}

// FormatGraph renders a text-based graph of key changes across versions.
func FormatGraph(nodes []GraphNode) string {
	if len(nodes) == 0 {
		return "(no graph data)\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s  %s\n", "KEY", "VERSION HISTORY"))
	sb.WriteString(strings.Repeat("-", 60) + "\n")
	for _, node := range nodes {
		var points []string
		for _, p := range node.Versions {
			sym := symbolForStatus(p.Status)
			points = append(points, fmt.Sprintf("v%d:%s", p.Version, sym))
		}
		sb.WriteString(fmt.Sprintf("%-30s  %s\n", node.Key, strings.Join(points, " → ")))
	}
	return sb.String()
}

func symbolForStatus(status ChangeStatus) string {
	switch status {
	case StatusAdded:
		return "+"
	case StatusRemoved:
		return "-"
	case StatusModified:
		return "~"
	default:
		return "="
	}
}
