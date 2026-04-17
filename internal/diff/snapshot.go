package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot captures diff entries at a point in time for later comparison.
type Snapshot struct {
	Path      string       `json:"path"`
	Version   int          `json:"version"`
	CapturedAt time.Time   `json:"captured_at"`
	Entries   []Entry      `json:"entries"`
}

// CaptureSnapshot creates a new Snapshot from the given entries.
func CaptureSnapshot(path string, version int, entries []Entry) Snapshot {
	return Snapshot{
		Path:       path,
		Version:    version,
		CapturedAt: time.Now().UTC(),
		Entries:    entries,
	}
}

// SaveSnapshot writes a Snapshot to a JSON file.
func SaveSnapshot(filePath string, s Snapshot) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("snapshot: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// LoadSnapshot reads a Snapshot from a JSON file.
func LoadSnapshot(filePath string) (Snapshot, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: open file: %w", err)
	}
	defer f.Close()
	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: decode: %w", err)
	}
	return s, nil
}

// DiffSnapshot compares two snapshots and returns change entries.
func DiffSnapshot(old, new Snapshot) []Entry {
	oldMap := make(map[string]string, len(old.Entries))
	for _, e := range old.Entries {
		oldMap[e.Key] = e.OldValue
		if e.Status == StatusAdded {
			oldMap[e.Key] = ""
		} else {
			oldMap[e.Key] = e.OldValue
		}
	}
	newMap := make(map[string]string, len(new.Entries))
	for _, e := range new.Entries {
		newMap[e.Key] = e.NewValue
	}
	return Compare(oldMap, newMap)
}
