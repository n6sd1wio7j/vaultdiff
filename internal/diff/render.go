package diff

import (
	"fmt"
	"io"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorGray   = "\033[90m"
)

// RenderOptions controls how the diff output is formatted.
type RenderOptions struct {
	Color     bool
	ShowUnchanged bool
	MaskValues   bool
}

// Render writes a human-readable diff of the entries to w.
func Render(w io.Writer, entries []DiffEntry, opts RenderOptions) {
	for _, e := range entries {
		if e.Change == Unchanged && !opts.ShowUnchanged {
			continue
		}

		prefix, color := symbolFor(e.Change, opts.Color)
		line := formatLine(prefix, color, e, opts.MaskValues)
		fmt.Fprintln(w, line)
	}
}

func symbolFor(c ChangeType, useColor bool) (string, string) {
	switch c {
	case Added:
		return "+", colorGreen
	case Removed:
		return "-", colorRed
	case Modified:
		return "~", colorYellow
	default:
		return " ", colorGray
	}
}

func formatLine(prefix, color string, e DiffEntry, mask bool) string {
	var sb strings.Builder

	if color != "" {
		sb.WriteString(color)
	}
	sb.WriteString(prefix)
	sb.WriteString(" ")
	sb.WriteString(e.Key)

	switch e.Change {
	case Added:
		sb.WriteString(": ")
		sb.WriteString(maskOrValue(e.NewValue, mask))
	case Removed:
		sb.WriteString(": ")
		sb.WriteString(maskOrValue(e.OldValue, mask))
	case Modified:
		sb.WriteString(": ")
		sb.WriteString(maskOrValue(e.OldValue, mask))
		sb.WriteString(" -> ")
		sb.WriteString(maskOrValue(e.NewValue, mask))
	case Unchanged:
		sb.WriteString(": ")
		sb.WriteString(maskOrValue(e.NewValue, mask))
	}

	if color != "" {
		sb.WriteString(colorReset)
	}
	return sb.String()
}

func maskOrValue(v interface{}, mask bool) string {
	if mask {
		return "***"
	}
	if v == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", v)
}
