package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/yourusername/vaultdiff/internal/diff"
)

// Entry represents a single audit log record for a diff operation.
type Entry struct {
	Timestamp   time.Time        `json:"timestamp"`
	Environment string           `json:"environment"`
	Path        string           `json:"path"`
	FromVersion int              `json:"from_version"`
	ToVersion   int              `json:"to_version"`
	Changes     []diff.DiffEntry `json:"changes"`
	HasChanges  bool             `json:"has_changes"`
	User        string           `json:"user,omitempty"`
}

// Logger writes audit entries to a destination writer.
type Logger struct {
	w io.Writer
}

// NewLogger creates a Logger writing to w. Pass os.Stdout for console output
// or an *os.File for file-based audit trails.
func NewLogger(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{w: w}
}

// Record serialises an audit Entry as a JSON line to the underlying writer.
func (l *Logger) Record(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", b)
	return err
}
