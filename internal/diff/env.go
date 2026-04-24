package diff

import (
	"fmt"
	"sort"
	"strings"
)

// EnvCompareOptions controls environment comparison behaviour.
type EnvCompareOptions struct {
	EnvA   string
	EnvB   string
	Ignore []string // key prefixes to skip
}

// EnvDiff holds the result of comparing two environment snapshots.
type EnvDiff struct {
	EnvA    string
	EnvB    string
	Entries []Entry
}

// DefaultEnvCompareOptions returns sensible defaults.
func DefaultEnvCompareOptions() EnvCompareOptions {
	return EnvCompareOptions{
		EnvA: "staging",
		EnvB: "production",
	}
}

// CompareEnvs diffs secrets from two named environments.
// secretsA and secretsB are key→value maps for each environment.
func CompareEnvs(secretsA, secretsB map[string]string, opts EnvCompareOptions) EnvDiff {
	entries := Compare(secretsA, secretsB)

	filtered := entries[:0]
	for _, e := range entries {
		skip := false
		for _, pfx := range opts.Ignore {
			if strings.HasPrefix(e.Key, pfx) {
				skip = true
				break
			}
		}
		if !skip {
			filtered = append(filtered, e)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Key < filtered[j].Key
	})

	return EnvDiff{
		EnvA:    opts.EnvA,
		EnvB:    opts.EnvB,
		Entries: filtered,
	}
}

// FormatEnvDiff returns a human-readable summary of an EnvDiff.
func FormatEnvDiff(d EnvDiff) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("env diff: %s → %s\n", d.EnvA, d.EnvB))
	sb.WriteString(strings.Repeat("-", 40) + "\n")
	for _, e := range d.Entries {
		sb.WriteString(formatLine(e, RenderOptions{MaskSecrets: true}))
		sb.WriteString("\n")
	}
	if len(d.Entries) == 0 {
		sb.WriteString("no differences found\n")
	}
	return sb.String()
}
