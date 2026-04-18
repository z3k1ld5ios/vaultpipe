package validate_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpipe/internal/validate"
)

func TestValidate_AllRulesPass(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "hunter2", "API_KEY": "abc123"}
	err := validate.Validate(secrets, validate.Rule{
		RequiredKeys:      []string{"DB_PASS", "API_KEY"},
		ForbidEmptyValues: true,
		MaxValueLength:    64,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidate_MissingRequiredKey(t *testing.T) {
	secrets := map[string]string{"API_KEY": "abc"}
	err := validate.Validate(secrets, validate.Rule{RequiredKeys: []string{"DB_PASS"}})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	if !strings.Contains(err.Error(), "DB_PASS") {
		t.Errorf("expected key name in error, got: %v", err)
	}
}

func TestValidate_EmptyRequiredKey(t *testing.T) {
	secrets := map[string]string{"TOKEN": ""}
	err := validate.Validate(secrets, validate.Rule{RequiredKeys: []string{"TOKEN"}})
	if err == nil {
		t.Fatal("expected error for empty required key")
	}
}

func TestValidate_ForbidEmptyValues(t *testing.T) {
	secrets := map[string]string{"A": "ok", "B": ""}
	err := validate.Validate(secrets, validate.Rule{ForbidEmptyValues: true})
	if err == nil {
		t.Fatal("expected error for empty value")
	}
	if !strings.Contains(err.Error(), "\"B\"") {
		t.Errorf("expected key B in error, got: %v", err)
	}
}

func TestValidate_MaxValueLength(t *testing.T) {
	secrets := map[string]string{"KEY": strings.Repeat("x", 200)}
	err := validate.Validate(secrets, validate.Rule{MaxValueLength: 100})
	if err == nil {
		t.Fatal("expected error for oversized value")
	}
}

func TestValidate_MultipleViolations(t *testing.T) {
	secrets := map[string]string{"A": ""}
	err := validate.Validate(secrets, validate.Rule{
		RequiredKeys:      []string{"MISSING"},
		ForbidEmptyValues: true,
	})
	ve, ok := err.(*validate.ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Violations) < 2 {
		t.Errorf("expected at least 2 violations, got %d", len(ve.Violations))
	}
}

func TestMustValidate_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	validate.MustValidate(map[string]string{}, validate.Rule{RequiredKeys: []string{"MUST_EXIST"}})
}
