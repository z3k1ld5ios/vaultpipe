package filter_test

import (
	"testing"

	"github.com/your-org/vaultpipe/internal/filter"
)

var base = map[string]string{
	"DB_PASSWORD": "secret",
	"DB_USER":     "admin",
	"API_KEY":     "key123",
	"DEBUG":       "true",
}

func TestFilter_NoRules_ReturnsAll(t *testing.T) {
	out := filter.Filter(base, filter.Config{})
	if len(out) != len(base) {
		t.Fatalf("expected %d keys, got %d", len(base), len(out))
	}
}

func TestFilter_AllowKeys_LimitsOutput(t *testing.T) {
	out := filter.Filter(base, filter.Config{AllowKeys: []string{"DB_PASSWORD", "API_KEY"}})
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["DB_USER"]; ok {
		t.Error("DB_USER should have been filtered out")
	}
}

func TestFilter_DenyKeys_ExcludesMatches(t *testing.T) {
	out := filter.Filter(base, filter.Config{DenyKeys: []string{"DB_"}})
	if _, ok := out["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should be denied")
	}
	if _, ok := out["DB_USER"]; ok {
		t.Error("DB_USER should be denied")
	}
	if _, ok := out["API_KEY"]; !ok {
		t.Error("API_KEY should remain")
	}
}

func TestFilter_AllowAndDeny_DenyWins(t *testing.T) {
	out := filter.Filter(base, filter.Config{
		AllowKeys: []string{"DB_PASSWORD", "API_KEY"},
		DenyKeys:  []string{"DB_PASSWORD"},
	})
	if _, ok := out["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should be denied even if allowed")
	}
	if _, ok := out["API_KEY"]; !ok {
		t.Error("API_KEY should be present")
	}
}

func TestFilter_EmptySecrets_ReturnsEmpty(t *testing.T) {
	out := filter.Filter(map[string]string{}, filter.Config{AllowKeys: []string{"X"}})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d entries", len(out))
	}
}
