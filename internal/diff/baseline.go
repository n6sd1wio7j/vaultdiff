package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Baseline represents a saved snapshot of diff entries for future comparison.
type Baseline struct {
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
	Entries   []Entry   `json:"entries"`
}

// SaveBaseline writes a baseline snapshot to the given file path.
func SaveBaseline(filePath string, secretPath string, entries []Entry) error {
	b := Baseline{
		Path:      secretPath,
		CreatedAt: time.Now().UTC(),
		Entries:   entries,
	}
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("baseline: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(b)
}

// LoadBaseline reads a baseline snapshot from the given file path.
func LoadBaseline(filePath string) (*Baseline, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("baseline: open file: %w", err)
	}
	defer f.Close()
	var b Baseline
	if err := json.NewDecoder(f).Decode(&b); err != nil {
		return nil, fmt.Errorf("baseline: decode: %w", err)
	}
	return &b, nil
}

// CompareToBaseline diffs current entries against a saved baseline.
func CompareToBaseline(baseline *Baseline, current map[string]string) []Entry {
	previous := make(map[string]string, len(baseline.Entries))
	for _, e := range baseline.Entries {
		previous[e.Key] = e.NewValue
	}
	return Compare(previous, current)
}
