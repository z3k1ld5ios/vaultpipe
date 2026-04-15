package render

import (
	"testing"
)

func TestSanitizeKey_Valid(t *testing.T) {
	valid := []string{"FOO", "foo", "FOO_BAR", "_PRIVATE", "A1B2", "MY_VAR_123"}
	for _, k := range valid {
		if err := SanitizeKey(k); err != nil {
			t.Errorf("expected %q to be valid, got error: %v", k, err)
		}
	}
}

func TestSanitizeKey_Invalid(t *testing.T) {
	invalid := []string{"", "1STARTS_DIGIT", "HAS-DASH", "HAS SPACE", "DOT.KEY"}
	for _, k := range invalid {
		if err := SanitizeKey(k); err == nil {
			t.Errorf("expected %q to be invalid, got nil error", k)
		}
	}
}

func TestSanitizeMap_AllValid(t *testing.T) {
	m := map[string]string{"FOO": "bar", "BAZ_QUX": "val"}
	if err := SanitizeMap(m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSanitizeMap_ContainsInvalid(t *testing.T) {
	m := map[string]string{"VALID": "ok", "bad-key": "oops"}
	if err := SanitizeMap(m); err == nil {
		t.Fatal("expected error for invalid key, got nil")
	}
}

func TestMaskValue_Short(t *testing.T) {
	if got := MaskValue("ab"); got != "**" {
		t.Errorf("expected '**', got %q", got)
	}
}

func TestMaskValue_Long(t *testing.T) {
	got := MaskValue("s3cr3t")
	if got != "s3****" {
		t.Errorf("expected 's3****', got %q", got)
	}
}

func TestMaskValue_Empty(t *testing.T) {
	if got := MaskValue(""); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}
