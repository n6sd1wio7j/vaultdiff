package diff

import (
	"strings"
	"unicode"
)

// NormalizeOptions controls how keys and values are normalized before diffing.
type NormalizeOptions struct {
	TrimSpace      bool
	LowercaseKeys  bool
	StripNonPrint  bool
}

// DefaultNormalizeOptions returns sensible normalization defaults.
func DefaultNormalizeOptions() NormalizeOptions {
	return NormalizeOptions{
		TrimSpace:     true,
		LowercaseKeys: false,
		StripNonPrint: true,
	}
}

// NormalizeKey applies key normalization rules.
func NormalizeKey(key string, opts NormalizeOptions) string {
	if opts.TrimSpace {
		key = strings.TrimSpace(key)
	}
	if opts.LowercaseKeys {
		key = strings.ToLower(key)
	}
	return key
}

// NormalizeValue applies value normalization rules.
func NormalizeValue(val string, opts NormalizeOptions) string {
	if opts.TrimSpace {
		val = strings.TrimSpace(val)
	}
	if opts.StripNonPrint {
		val = strings.Map(func(r rune) rune {
			if unicode.IsPrint(r) || r == '\t' || r == '\n' {
				return r
			}
			return -1
		}, val)
	}
	return val
}

// NormalizeSecrets returns new maps with normalized keys and values.
func NormalizeSecrets(a, b map[string]string, opts NormalizeOptions) (map[string]string, map[string]string) {
	normalize := func(m map[string]string) map[string]string {
		out := make(map[string]string, len(m))
		for k, v := range m {
			out[NormalizeKey(k, opts)] = NormalizeValue(v, opts)
		}
		return out
	}
	return normalize(a), normalize(b)
}
