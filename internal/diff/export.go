package diff

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ExportFormat defines the output format for exporting diff results.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatCSV  ExportFormat = "csv"
	FormatText ExportFormat = "text"
)

// ExportOptions controls how the diff is exported.
type ExportOptions struct {
	Format      ExportFormat
	MaskSecrets bool
}

// Export writes the diff entries to w in the specified format.
func Export(entries []Entry, opts ExportOptions, w io.Writer) error {
	switch opts.Format {
	case FormatJSON:
		return exportJSON(entries, opts.MaskSecrets, w)
	case FormatCSV:
		return exportCSV(entries, opts.MaskSecrets, w)
	case FormatText:
		return exportText(entries, opts.MaskSecrets, w)
	default:
		return fmt.Errorf("unsupported export format: %q", opts.Format)
	}
}

func exportJSON(entries []Entry, mask bool, w io.Writer) error {
	type jsonEntry struct {
		Key    string `json:"key"`
		Status string `json:"status"`
		OldVal string `json:"old_value,omitempty"`
		NewVal string `json:"new_value,omitempty"`
	}
	out := make([]jsonEntry, 0, len(entries))
	for _, e := range entries {
		out = append(out, jsonEntry{
			Key:    e.Key,
			Status: string(e.Status),
			OldVal: maskOrValue(e.OldValue, mask),
			NewVal: maskOrValue(e.NewValue, mask),
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func exportCSV(entries []Entry, mask bool, w io.Writer) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"key", "status", "old_value", "new_value"}); err != nil {
		return err
	}
	for _, e := range entries {
		row := []string{e.Key, string(e.Status), maskOrValue(e.OldValue, mask), maskOrValue(e.NewValue, mask)}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func exportText(entries []Entry, mask bool, w io.Writer) error {
	var sb strings.Builder
	for _, e := range entries {
		sb.WriteString(formatLine(e, mask))
		sb.WriteByte('\n')
	}
	_, err := fmt.Fprint(w, sb.String())
	return err
}
