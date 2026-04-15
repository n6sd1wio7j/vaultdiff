package diff

import (
	"bytes"
	"strings"
	"testing"
)

func TestRender_AddedEntry(t *testing.T) {
	var buf bytes.Buffer
	entries := []DiffEntry{
		{Key: "new_key", Change: Added, NewValue: "new_val"},
	}
	Render(&buf, entries, RenderOptions{})
	out := buf.String()
	if !strings.Contains(out, "+ new_key") {
		t.Errorf("expected '+ new_key' in output, got: %s", out)
	}
	if !strings.Contains(out, "new_val") {
		t.Errorf("expected value 'new_val' in output, got: %s", out)
	}
}

func TestRender_RemovedEntry(t *testing.T) {
	var buf bytes.Buffer
	entries := []DiffEntry{
		{Key: "old_key", Change: Removed, OldValue: "old_val"},
	}
	Render(&buf, entries, RenderOptions{})
	out := buf.String()
	if !strings.Contains(out, "- old_key") {
		t.Errorf("expected '- old_key' in output, got: %s", out)
	}
}

func TestRender_ModifiedEntry(t *testing.T) {
	var buf bytes.Buffer
	entries := []DiffEntry{
		{Key: "db_pass", Change: Modified, OldValue: "old", NewValue: "new"},
	}
	Render(&buf, entries, RenderOptions{})
	out := buf.String()
	if !strings.Contains(out, "old -> new") {
		t.Errorf("expected 'old -> new' in output, got: %s", out)
	}
}

func TestRender_UnchangedHiddenByDefault(t *testing.T) {
	var buf bytes.Buffer
	entries := []DiffEntry{
		{Key: "stable", Change: Unchanged, NewValue: "val"},
	}
	Render(&buf, entries, RenderOptions{ShowUnchanged: false})
	if buf.Len() != 0 {
		t.Errorf("expected no output for unchanged entry, got: %s", buf.String())
	}
}

func TestRender_UnchangedShownWhenEnabled(t *testing.T) {
	var buf bytes.Buffer
	entries := []DiffEntry{
		{Key: "stable", Change: Unchanged, NewValue: "val"},
	}
	Render(&buf, entries, RenderOptions{ShowUnchanged: true})
	if buf.Len() == 0 {
		t.Error("expected output for unchanged entry when ShowUnchanged is true")
	}
}

func TestRender_MaskValues(t *testing.T) {
	var buf bytes.Buffer
	entries := []DiffEntry{
		{Key: "secret", Change: Added, NewValue: "super_secret_value"},
	}
	Render(&buf, entries, RenderOptions{MaskValues: true})
	out := buf.String()
	if strings.Contains(out, "super_secret_value") {
		t.Error("expected value to be masked, but found plaintext")
	}
	if !strings.Contains(out, "***") {
		t.Errorf("expected '***' mask in output, got: %s", out)
	}
}
