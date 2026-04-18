package env

import (
	"testing"
)

func TestWhitelist_EmptyAllowsAll(t *testing.T) {
	w := NewWhitelist(nil, nil)
	if !w.Allow("ANY_KEY") {
		t.Fatal("empty whitelist should allow everything")
	}
}

func TestWhitelist_ExactKeyMatch(t *testing.T) {
	w := NewWhitelist([]string{"DB_PASSWORD"}, nil)
	if !w.Allow("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD to be allowed")
	}
	if w.Allow("DB_USER") {
		t.Error("expected DB_USER to be rejected")
	}
}

func TestWhitelist_CaseInsensitiveKey(t *testing.T) {
	w := NewWhitelist([]string{"db_password"}, nil)
	if !w.Allow("DB_PASSWORD") {
		t.Error("key match should be case-insensitive")
	}
}

func TestWhitelist_PrefixMatch(t *testing.T) {
	w := NewWhitelist(nil, []string{"APP_"})
	if !w.Allow("APP_SECRET") {
		t.Error("expected APP_SECRET to match prefix APP_")
	}
	if w.Allow("DB_SECRET") {
		t.Error("expected DB_SECRET to be rejected")
	}
}

func TestWhitelist_Filter(t *testing.T) {
	w := NewWhitelist([]string{"TOKEN"}, []string{"APP_"})
	input := []string{"TOKEN=abc", "APP_KEY=x", "OTHER=y"}
	out := w.Filter(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestWhitelist_Validate_AllAllowed(t *testing.T) {
	w := NewWhitelist([]string{"FOO", "BAR"}, nil)
	err := w.Validate(map[string]string{"FOO": "1", "BAR": "2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWhitelist_Validate_RejectedKey(t *testing.T) {
	w := NewWhitelist([]string{"FOO"}, nil)
	err := w.Validate(map[string]string{"FOO": "1", "SECRET": "x"})
	if err == nil {
		t.Fatal("expected validation error for rejected key")
	}
}
