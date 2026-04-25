package diff

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// ChecksumEntry holds a key and its SHA-256 digest.
type ChecksumEntry struct {
	Key      string `json:"key"`
	Checksum string `json:"checksum"`
}

// ChecksumResult holds the per-key digests and an aggregate digest for the
// entire secret version.
type ChecksumResult struct {
	Path      string          `json:"path"`
	Version   int             `json:"version"`
	Aggregate string          `json:"aggregate"`
	Entries   []ChecksumEntry `json:"entries"`
}

// ComputeChecksums returns a ChecksumResult for the provided key-value map.
// The aggregate digest is derived from all key=checksum pairs sorted by key so
// that the result is deterministic regardless of map iteration order.
func ComputeChecksums(path string, version int, secrets map[string]string) ChecksumResult {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]ChecksumEntry, 0, len(keys))
	for _, k := range keys {
		h := sha256.Sum256([]byte(secrets[k]))
		entries = append(entries, ChecksumEntry{
			Key:      k,
			Checksum: hex.EncodeToString(h[:]),
		})
	}

	// Build aggregate from sorted "key=checksum" pairs.
	parts := make([]string, len(entries))
	for i, e := range entries {
		parts[i] = e.Key + "=" + e.Checksum
	}
	agg := sha256.Sum256([]byte(strings.Join(parts, "\n")))

	return ChecksumResult{
		Path:      path,
		Version:   version,
		Aggregate: hex.EncodeToString(agg[:]),
		Entries:   entries,
	}
}

// FormatChecksums returns a human-readable representation of a ChecksumResult.
func FormatChecksums(r ChecksumResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("path: %s  version: %d\n", r.Path, r.Version))
	sb.WriteString(fmt.Sprintf("aggregate: %s\n", r.Aggregate))
	sb.WriteString(strings.Repeat("-", 72) + "\n")
	for _, e := range r.Entries {
		sb.WriteString(fmt.Sprintf("  %-40s %s\n", e.Key, e.Checksum))
	}
	return sb.String()
}
