package diff

import (
	"fmt"
	"sort"
	"strings"
)

// GroupOptions controls how entries are grouped.
type GroupOptions struct {
	ByStatus bool
	ByPrefix bool
	PrefixSep string
}

// Group holds a named collection of diff entries.
type Group struct {
	Name    string
	Entries []Entry
}

// GroupEntries partitions entries according to GroupOptions.
func GroupEntries(entries []Entry, opts GroupOptions) []Group {
	if opts.ByStatus {
		return groupByStatus(entries)
	}
	if opts.ByPrefix {
		sep := opts.PrefixSep
		if sep == "" {
			sep = "/"
		}
		return groupByPrefix(entries, sep)
	}
	return []Group{{Name: "all", Entries: entries}}
}

func groupByStatus(entries []Entry) []Group {
	m := map[string][]Entry{}
	order := []string{}
	for _, e := range entries {
		s := string(e.Status)
		if _, ok := m[s]; !ok {
			order = append(order, s)
		}
		m[s] = append(m[s], e)
	}
	groups := make([]Group, 0, len(order))
	for _, k := range order {
		groups = append(groups, Group{Name: k, Entries: m[k]})
	}
	return groups
}

func groupByPrefix(entries []Entry, sep string) []Group {
	m := map[string][]Entry{}
	for _, e := range entries {
		parts := strings.SplitN(e.Key, sep, 2)
		prefix := parts[0]
		if len(parts) == 1 {
			prefix = "(root)"
		}
		m[prefix] = append(m[prefix], e)
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	groups := make([]Group, 0, len(keys))
	for _, k := range keys {
		groups = append(groups, Group{Name: k, Entries: m[k]})
	}
	return groups
}

// FormatGroups renders grouped entries as a string.
func FormatGroups(groups []Group) string {
	var sb strings.Builder
	for _, g := range groups {
		sb.WriteString(fmt.Sprintf("[%s] (%d entries)\n", g.Name, len(g.Entries)))
		for _, e := range g.Entries {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", e.Key, e.Status))
		}
	}
	return sb.String()
}
