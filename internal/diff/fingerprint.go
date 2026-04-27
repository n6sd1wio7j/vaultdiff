package diff

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// FingerprintOptions controls how fingerprints are computed.
type FingerprintOptions struct {
	// IncludeUnchanged includes unchanged entries in the fingerprint.
	IncludeUnchanged bool
	// MaskValues replaces values with their SHA256 hash before fingerprinting.
	MaskValues bool
}

// DefaultFingerprintOptions returns sensible defaults.
func DefaultFingerprintOptions() FingerprintOptions {
	return FingerprintOptions{
		IncludeUnchanged: true,
		MaskValues:       false,
	}
}

// FingerprintResult holds the computed fingerprint for a secret path.
type FingerprintResult struct {
	Path        string
	Version     int
	Fingerprint string
	EntryCount  int
}

// ComputeFingerprint produces a stable SHA256 fingerprint over the given entries.
func ComputeFingerprint(path string, version int, entries []Entry, opts FingerprintOptions) FingerprintResult {
	filtered := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if !opts.IncludeUnchanged && e.Status == StatusUnchanged {
			continue
		}
		filtered = append(filtered, e)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Key < filtered[j].Key
	})

	h := sha256.New()
	for _, e := range filtered {
		v := e.NewValue
		if opts.MaskValues {
			sum := sha256.Sum256([]byte(v))
			v = hex.EncodeToString(sum[:])
		}
		fmt.Fprintf(h, "%s=%s:%s\n", e.Key, v, string(e.Status))
	}

	return FingerprintResult{
		Path:        path,
		Version:     version,
		Fingerprint: hex.EncodeToString(h.Sum(nil)),
		EntryCount:  len(filtered),
	}
}

// FormatFingerprint returns a human-readable representation of a FingerprintResult.
func FormatFingerprint(r FingerprintResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Path:        %s\n", r.Path))
	sb.WriteString(fmt.Sprintf("Version:     %d\n", r.Version))
	sb.WriteString(fmt.Sprintf("Entries:     %d\n", r.EntryCount))
	sb.WriteString(fmt.Sprintf("Fingerprint: %s\n", r.Fingerprint))
	return sb.String()
}
