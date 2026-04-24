package env

import (
	"testing"
)

func TestNewScope_EmptyAllowsAll(t *testing.T) {
	s := NewScope(nil)
	if !s.Allows("ANY_KEY") {
		t.Error("empty scope should allow all keys")
	}
}

func TestNewScope_ExactMatch(t *testing.T) {
	s := NewScope([]string{"DB_PASSWORD", "API_KEY"})
	if !s.Allows("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD to be allowed")
	}
	if !s.Allows("API_KEY") {
		t.Error("expected API_KEY to be allowed")
	}
	if s.Allows("SECRET_TOKEN") {
		t.Error("expected SECRET_TOKEN to be denied")
	}
}

func TestNewScope_PrefixMatch(t *testing.T) {
	s := NewScope([]string{"APP_*"})
	if !s.Allows("APP_NAME") {
		t.Error("expected APP_NAME to be allowed by prefix")
	}
	if !s.Allows("APP_VERSION") {
		t.Error("expected APP_VERSION to be allowed by prefix")
	}
	if s.Allows("DB_HOST") {
		t.Error("expected DB_HOST to be denied")
	}
}

func TestNewScope_MixedExactAndPrefix(t *testing.T) {
	s := NewScope([]string{"DB_HOST", "APP_*"})
	if !s.Allows("DB_HOST") {
		t.Error("expected DB_HOST to be allowed")
	}
	if !s.Allows("APP_ENV") {
		t.Error("expected APP_ENV to be allowed by prefix")
	}
	if s.Allows("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD to be denied")
	}
}

func TestFilter_RetainsAllowedKeys(t *testing.T) {
	s := NewScope([]string{"APP_*", "LOG_LEVEL"})
	input := map[string]string{
		"APP_NAME":  "vaultpipe",
		"APP_ENV":   "prod",
		"LOG_LEVEL": "info",
		"DB_PASS":   "secret",
	}
	out := s.Filter(input)
	if len(out) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(out))
	}
	if _, ok := out["DB_PASS"]; ok {
		t.Error("DB_PASS should have been filtered out")
	}
}

func TestFilter_EmptyScope_ReturnsAll(t *testing.T) {
	s := NewScope([]string{})
	input := map[string]string{"A": "1", "B": "2"}
	out := s.Filter(input)
	if len(out) != len(input) {
		t.Errorf("expected all %d keys, got %d", len(input), len(out))
	}
}

func TestFilter_DoesNotMutateInput(t *testing.T) {
	s := NewScope([]string{"KEEP"})
	input := map[string]string{"KEEP": "yes", "DROP": "no"}
	s.Filter(input)
	if _, ok := input["DROP"]; !ok {
		t.Error("Filter must not mutate the input map")
	}
}
