package diff

import (
	"strings"
	"testing"
)

var sampleGroupEntries = []Entry{
	{Key: "db/password", Status: StatusModified, OldValue: "old", NewValue: "new"},
	{Key: "db/user", Status: StatusUnchanged, OldValue: "admin", NewValue: "admin"},
	{Key: "app/token", Status: StatusAdded, NewValue: "tok"},
	{Key: "app/debug", Status: StatusRemoved, OldValue: "true"},
	{Key: "timeout", Status: StatusAdded, NewValue: "30s"},
}

func TestGroupEntries_ByStatus(t *testing.T) {
	groups := GroupEntries(sampleGroupEntries, GroupOptions{ByStatus: true})
	statuses := map[string]bool{}
	for _, g := range groups {
		statuses[g.Name] = true
		for _, e := range g.Entries {
			if string(e.Status) != g.Name {
				t.Errorf("entry %s in wrong group %s", e.Key, g.Name)
			}
		}
	}
	if !statuses[string(StatusModified)] {
		t.Error("expected modified group")
	}
}

func TestGroupEntries_ByPrefix(t *testing.T) {
	groups := GroupEntries(sampleGroupEntries, GroupOptions{ByPrefix: true, PrefixSep: "/"})
	names := map[string]int{}
	for _, g := range groups {
		names[g.Name] = len(g.Entries)
	}
	if names["db"] != 2 {
		t.Errorf("expected 2 db entries, got %d", names["db"])
	}
	if names["app"] != 2 {
		t.Errorf("expected 2 app entries, got %d", names["app"])
	}
	if names["(root)"] != 1 {
		t.Errorf("expected 1 root entry, got %d", names["(root)"])
	}
}

func TestGroupEntries_DefaultSep(t *testing.T) {
	groups := GroupEntries(sampleGroupEntries, GroupOptions{ByPrefix: true})
	if len(groups) == 0 {
		t.Fatal("expected groups")
	}
}

func TestGroupEntries_NoOption(t *testing.T) {
	groups := GroupEntries(sampleGroupEntries, GroupOptions{})
	if len(groups) != 1 || groups[0].Name != "all" {
		t.Errorf("expected single 'all' group, got %v", groups)
	}
}

func TestFormatGroups_ContainsName(t *testing.T) {
	groups := GroupEntries(sampleGroupEntries, GroupOptions{ByPrefix: true, PrefixSep: "/"})
	out := FormatGroups(groups)
	if !strings.Contains(out, "[db]") {
		t.Error("expected [db] in output")
	}
	if !strings.Contains(out, "db/password") {
		t.Error("expected db/password in output")
	}
}
