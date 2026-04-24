package env

import (
	"testing"
)

func TestSchema_AllValid(t *testing.T) {
	schema := NewSchema([]FieldSchema{
		{Key: "PORT", Type: TypeInt, Required: true},
		{Key: "DEBUG", Type: TypeBool},
		{Key: "NAME", Type: TypeString, Required: true},
	})
	env := map[string]string{"PORT": "8080", "DEBUG": "true", "NAME": "vaultpipe"}
	if errs := schema.Validate(env); errs != nil {
		t.Fatalf("expected no errors, got: %v", errs)
	}
}

func TestSchema_MissingRequiredKey(t *testing.T) {
	schema := NewSchema([]FieldSchema{
		{Key: "TOKEN", Type: TypeString, Required: true},
	})
	errs := schema.Validate(map[string]string{})
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestSchema_InvalidInt(t *testing.T) {
	schema := NewSchema([]FieldSchema{
		{Key: "PORT", Type: TypeInt},
	})
	errs := schema.Validate(map[string]string{"PORT": "not-a-number"})
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestSchema_InvalidBool(t *testing.T) {
	schema := NewSchema([]FieldSchema{
		{Key: "ENABLED", Type: TypeBool},
	})
	errs := schema.Validate(map[string]string{"ENABLED": "yes"})
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestSchema_PatternMatch(t *testing.T) {
	schema := NewSchema([]FieldSchema{
		{Key: "LOG_LEVEL", Type: TypeString, Pattern: `^(debug|info|warn|error)$`},
	})
	if errs := schema.Validate(map[string]string{"LOG_LEVEL": "info"}); errs != nil {
		t.Fatalf("expected no errors, got: %v", errs)
	}
}

func TestSchema_PatternMismatch(t *testing.T) {
	schema := NewSchema([]FieldSchema{
		{Key: "LOG_LEVEL", Type: TypeString, Pattern: `^(debug|info|warn|error)$`},
	})
	errs := schema.Validate(map[string]string{"LOG_LEVEL": "verbose"})
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestSchema_OptionalMissingKey_NoError(t *testing.T) {
	schema := NewSchema([]FieldSchema{
		{Key: "OPTIONAL", Type: TypeString, Required: false},
	})
	if errs := schema.Validate(map[string]string{}); errs != nil {
		t.Fatalf("expected no errors for missing optional key, got: %v", errs)
	}
}

func TestSchema_MultipleErrors(t *testing.T) {
	schema := NewSchema([]FieldSchema{
		{Key: "PORT", Type: TypeInt, Required: true},
		{Key: "MODE", Type: TypeString, Pattern: `^(prod|dev)$`, Required: true},
	})
	errs := schema.Validate(map[string]string{"PORT": "abc", "MODE": "staging"})
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d: %v", len(errs), errs)
	}
}
