package audit

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// FileLogger wraps Logger with a file-backed writer that is safe to close.
type FileLogger struct {
	*Logger
	f *os.File
}

// OpenFileLogger opens (or creates) a JSONL audit log at path.
// The caller must call Close when done.
func OpenFileLogger(path string) (*FileLogger, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("audit: create log dir: %w", err)
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o640)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	return &FileLogger{Logger: NewLogger(f), f: f}, nil
}

// Close flushes and closes the underlying file.
func (fl *FileLogger) Close() error {
	if err := fl.f.Sync(); err != nil {
		// Still attempt to close even if sync fails.
		_ = fl.f.Close()
		return fmt.Errorf("audit: sync log file: %w", err)
	}
	return fl.f.Close()
}

// DefaultLogPath returns a conventional log path using today's date.
func DefaultLogPath(baseDir string) string {
	date := time.Now().UTC().Format("2006-01-02")
	return filepath.Join(baseDir, fmt.Sprintf("vaultdiff-%s.jsonl", date))
}
