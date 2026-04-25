package diff

import (
	"fmt"
	"strings"
	"time"
)

// ChainEntry represents a single step in a diff chain across multiple versions.
type ChainEntry struct {
	Version int
	Timestamp time.Time
	Entries []Entry
}

// Chain holds an ordered sequence of diffs across versions.
type Chain struct {
	Path    string
	Steps   []ChainEntry
}

// BuildChain constructs a diff chain from a sequence of secret version maps.
// versions must be ordered oldest to newest; each element is a key→value map.
func BuildChain(path string, versions []map[string]string) Chain {
	chain := Chain{Path: path}
	if len(versions) < 2 {
		return chain
	}
	for i := 1; i < len(versions); i++ {
		entries := Compare(versions[i-1], versions[i])
		step := ChainEntry{
			Version:   i + 1,
			Timestamp: time.Now().UTC(),
			Entries:   entries,
		}
		chain.Steps = append(chain.Steps, step)
	}
	return chain
}

// HasAnyDrift reports whether any step in the chain contains changes.
func (c Chain) HasAnyDrift() bool {
	for _, step := range c.Steps {
		if HasChanges(step.Entries) {
			return true
		}
	}
	return false
}

// FormatChain returns a human-readable multi-step diff report.
func FormatChain(c Chain) string {
	if len(c.Steps) == 0 {
		return fmt.Sprintf("chain: %s — no steps recorded\n", c.Path)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("diff chain: %s (%d step(s))\n", c.Path, len(c.Steps)))
	for _, step := range c.Steps {
		sb.WriteString(fmt.Sprintf("  version %d (%s):\n", step.Version, step.Timestamp.Format(time.RFC3339)))
		for _, e := range step.Entries {
			if e.Status != StatusUnchanged {
				sb.WriteString(fmt.Sprintf("    [%s] %s\n", e.Status, e.Key))
			}
		}
		if !HasChanges(step.Entries) {
			sb.WriteString("    (no changes)\n")
		}
	}
	return sb.String()
}
