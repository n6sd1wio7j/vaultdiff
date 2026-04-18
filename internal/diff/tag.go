package diff

import "fmt"

// Tag represents a label attached to a diff entry.
type Tag struct {
	Key   string
	Value string
}

// TaggedEntry pairs a DiffEntry with a set of tags.
type TaggedEntry struct {
	Entry DiffEntry
	Tags  []Tag
}

// TagOptions controls which tags are applied.
type TagOptions struct {
	Env     string
	Version string
	Custom  map[string]string
}

// TagEntries attaches tags to each entry based on the provided options.
func TagEntries(entries []DiffEntry, opts TagOptions) []TaggedEntry {
	result := make([]TaggedEntry, 0, len(entries))
	for _, e := range entries {
		tags := []Tag{}
		if opts.Env != "" {
			tags = append(tags, Tag{Key: "env", Value: opts.Env})
		}
		if opts.Version != "" {
			tags = append(tags, Tag{Key: "version", Value: opts.Version})
		}
		for k, v := range opts.Custom {
			tags = append(tags, Tag{Key: k, Value: v})
		}
		result = append(result, TaggedEntry{Entry: e, Tags: tags})
	}
	return result
}

// FormatTagged returns a human-readable string for a TaggedEntry.
func FormatTagged(te TaggedEntry) string {
	line := fmt.Sprintf("[%s] %s", te.Entry.Status, te.Entry.Key)
	for _, t := range te.Tags {
		line += fmt.Sprintf(" %s=%s", t.Key, t.Value)
	}
	return line
}
