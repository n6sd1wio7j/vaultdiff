package diff

import "fmt"

// Annotation holds a human-readable note attached to a diff entry key.
type Annotation struct {
	Key  string `json:"key"`
	Note string `json:"note"`
}

// AnnotateOptions configures how annotations are applied.
type AnnotateOptions struct {
	Annotations []Annotation
}

// Annotate returns a map of key -> note for quick lookup.
func Annotate(entries []Entry, opts AnnotateOptions) map[string]string {
	noteMap := make(map[string]string, len(opts.Annotations))
	for _, a := range opts.Annotations {
		noteMap[a.Key] = a.Note
	}

	result := make(map[string]string)
	for _, e := range entries {
		if note, ok := noteMap[e.Key]; ok {
			result[e.Key] = note
		}
	}
	return result
}

// FormatAnnotated renders an entry line with an optional annotation note.
func FormatAnnotated(e Entry, note string, mask bool) string {
	base := formatLine(symbolFor(e.Change), e.Key, maskOrValue(e.NewValue, mask))
	if note != "" {
		return fmt.Sprintf("%s  # %s", base, note)
	}
	return base
}
