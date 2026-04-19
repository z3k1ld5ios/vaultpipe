package env

import (
	"testing"
)

func TestCoerce_NoRules_ReturnsCopy(t *testing.T) {
	c := NewCoercer(nil)
	in := map[string]string{"KEY": "value"}
	out, err := c.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "value" {
		t.Errorf("expected value, got %q", out["KEY"])
	}
}

func TestCoerce_Upper(t *testing.T) {
	c := NewCoercer(map[string]CoerceType{"MODE": CoerceUpper})
	out, err := c.Apply(map[string]string{"MODE": "production"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["MODE"] != "PRODUCTION" {
		t.Errorf("expected PRODUCTION, got %q", out["MODE"])
	}
}

func TestCoerce_Lower(t *testing.T) {
	c := NewCoercer(map[string]CoerceType{"ENV": CoerceLower})
	out, err := c.Apply(map[string]string{"ENV": "STAGING"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["ENV"] != "staging" {
		t.Errorf("expected staging, got %q", out["ENV"])
	}
}

func TestCoerce_Bool_Valid(t *testing.T) {
	c := NewCoercer(map[string]CoerceType{"FLAG": CoerceBool})
	out, err := c.Apply(map[string]string{"FLAG": "1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FLAG"] != "true" {
		t.Errorf("expected true, got %q", out["FLAG"])
	}
}

func TestCoerce_Bool_Invalid(t *testing.T) {
	c := NewCoercer(map[string]CoerceType{"FLAG": CoerceBool})
	_, err := c.Apply(map[string]string{"FLAG": "notabool"})
	if err == nil {
		t.Fatal("expected error for invalid bool")
	}
}

func TestCoerce_Int_Valid(t *testing.T) {
	c := NewCoercer(map[string]CoerceType{"PORT": CoerceInt})
	out, err := c.Apply(map[string]string{"PORT": "8080"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected 8080, got %q", out["PORT"])
	}
}

func TestCoerce_Int_Invalid(t *testing.T) {
	c := NewCoercer(map[string]CoerceType{"PORT": CoerceInt})
	_, err := c.Apply(map[string]string{"PORT": "abc"})
	if err == nil {
		t.Fatal("expected error for invalid int")
	}
}

func TestCoerce_MissingKey_Skipped(t *testing.T) {
	c := NewCoercer(map[string]CoerceType{"MISSING": CoerceUpper})
	out, err := c.Apply(map[string]string{"OTHER": "val"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["MISSING"]; ok {
		t.Error("missing key should not be added")
	}
}

func TestCoerce_DoesNotMutateInput(t *testing.T) {
	c := NewCoercer(map[string]CoerceType{"K": CoerceUpper})
	in := map[string]string{"K": "hello"}
	_, _ = c.Apply(in)
	if in["K"] != "hello" {
		t.Error("input map was mutated")
	}
}
