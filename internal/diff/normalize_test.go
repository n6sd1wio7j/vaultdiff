package diff

import (
	"testing"
)

func TestNormalizeKey_TrimSpace(t *testing.T) {
	opts := NormalizeOptions{TrimSpace: true}
	got := NormalizeKey("  MY_KEY  ", opts)
	if got != "MY_KEY" {
		t.Errorf("expected 'MY_KEY', got %q", got)
	}
}

func TestNormalizeKey_Lowercase(t *testing.T) {
	opts := NormalizeOptions{LowercaseKeys: true}
	got := NormalizeKey("MY_KEY", opts)
	if got != "my_key" {
		t.Errorf("expected 'my_key', got %q", got)
	}
}

func TestNormalizeKey_NoOp(t *testing.T) {
	opts := NormalizeOptions{}
	got := NormalizeKey("  MY_KEY  ", opts)
	if got != "  MY_KEY  " {
		t.Errorf("expected unchanged key, got %q", got)
	}
}

func TestNormalizeValue_TrimSpace(t *testing.T) {
	opts := NormalizeOptions{TrimSpace: true}
	got := NormalizeValue("  secret  ", opts)
	if got != "secret" {
		t.Errorf("expected 'secret', got %q", got)
	}
}

func TestNormalizeValue_StripNonPrint(t *testing.T) {
	opts := NormalizeOptions{StripNonPrint: true}
	got := NormalizeValue("val\x00ue", opts)
	if got != "value" {
		t.Errorf("expected 'value', got %q", got)
	}
}

func TestNormalizeValue_PreservesTab(t *testing.T) {
	opts := NormalizeOptions{StripNonPrint: true}
	got := NormalizeValue("val\tue", opts)
	if got != "val\tue" {
		t.Errorf("expected tab preserved, got %q", got)
	}
}

func TestNormalizeSecrets_AppliesBothMaps(t *testing.T) {
	a := map[string]string{"  KEY  ": "  val1  "}
	b := map[string]string{"  KEY  ": "  val2  "}
	opts := DefaultNormalizeOptions()
	na, nb := NormalizeSecrets(a, b, opts)
	if _, ok := na["KEY"]; !ok {
		t.Error("expected trimmed key 'KEY' in normalized map a")
	}
	if nb["KEY"] != "val2" {
		t.Errorf("expected 'val2', got %q", nb["KEY"])
	}
}

func TestDefaultNormalizeOptions_Defaults(t *testing.T) {
	opts := DefaultNormalizeOptions()
	if !opts.TrimSpace {
		t.Error("expected TrimSpace to be true")
	}
	if opts.LowercaseKeys {
		t.Error("expected LowercaseKeys to be false")
	}
	if !opts.StripNonPrint {
		t.Error("expected StripNonPrint to be true")
	}
}
