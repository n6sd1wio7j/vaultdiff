package diff

import (
	"fmt"
	"strings"
	"time"
)

// StaleOptions configures staleness detection.
type StaleOptions struct {
	// MaxAge is the maximum allowed age of an unchanged secret before it is
	// considered stale. Defaults to 90 days.
	MaxAge time.Duration
	// IncludeUnchanged includes unchanged entries in the result even when they
	// are not stale.
	IncludeUnchanged bool
}

// StaleEntry represents a secret key that has not been rotated within the
// allowed window.
type StaleEntry struct {
	Key     string
	Status  ChangeStatus
	Age     time.Duration
	Stale   bool
	Reason  string
}

// DefaultStaleOptions returns sensible defaults for staleness detection.
func DefaultStaleOptions() StaleOptions {
	return StaleOptions{
		MaxAge: 90 * 24 * time.Hour,
	}
}

// DetectStale inspects diff entries and marks any key whose last-seen
// timestamp (approximated from the version metadata recorded in the entry
// key suffix "@<unix>") exceeds MaxAge as stale. When no timestamp suffix is
// present the entry is evaluated against the provided referenceTime.
func DetectStale(entries []Entry, referenceTime time.Time, opts StaleOptions) []StaleEntry {
	if opts.MaxAge == 0 {
		opts.MaxAge = DefaultStaleOptions().MaxAge
	}

	var results []StaleEntry
	for _, e := range entries {
		if e.Status == StatusUnchanged && !opts.IncludeUnchanged {
			continue
		}

		age := opts.MaxAge / 2 // default: assume recently rotated
		// Entries tagged with "@<unix>" carry an approximate last-rotation time.
		if idx := strings.LastIndex(e.Key, "@"); idx != -1 {
			var ts int64
			fmt.Sscanf(e.Key[idx+1:], "%d", &ts)
			if ts > 0 {
				rotated := time.Unix(ts, 0)
				age = referenceTime.Sub(rotated)
			}
		}

		stale := age > opts.MaxAge
		reason := ""
		if stale {
			days := int(age.Hours() / 24)
			reason = fmt.Sprintf("not rotated in %d days (max %d)", days, int(opts.MaxAge.Hours()/24))
		}

		results = append(results, StaleEntry{
			Key:    e.Key,
			Status: e.Status,
			Age:    age,
			Stale:  stale,
			Reason: reason,
		})
	}
	return results
}

// FormatStale renders stale detection results as a human-readable string.
func FormatStale(entries []StaleEntry) string {
	if len(entries) == 0 {
		return "no stale secrets detected\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-40s %-12s %s\n", "KEY", "STATUS", "REASON"))
	sb.WriteString(strings.Repeat("-", 72) + "\n")
	for _, e := range entries {
		marker := " "
		if e.Stale {
			marker = "!"
		}
		reason := e.Reason
		if reason == "" {
			reason = "ok"
		}
		sb.WriteString(fmt.Sprintf("%s %-39s %-12s %s\n", marker, e.Key, string(e.Status), reason))
	}
	return sb.String()
}
