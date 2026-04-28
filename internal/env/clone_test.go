package env

import (
	"errors"
	"strings"
	"testing"
)

func TestClone_NilInput_ReturnsEmpty(t *testing.T) {
	c := NewCloner()
	out, err := c.Clone(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestClone_CopiesAllKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	c := NewCloner()
	out, err := c.Clone(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 || out["A"] != "1" || out["B"] != "2" {
		t.Errorf("unexpected clone result: %v", out)
	}
}

func TestClone_DoesNotMutateSource(t *testing.T) {
	src := map[string]string{"X": "original"}
	c := NewCloner()
	out, _ := c.Clone(src)
	out["X"] = "mutated"
	if src["X"] != "original" {
		t.Error("source map was mutated")
	}
}

func TestClone_WithFilter_ExcludesKeys(t *testing.T) {
	src := map[string]string{"KEEP": "yes", "SKIP": "no"}
	c := NewCloner(WithCloneFilter(func(k string) bool {
		return k == "KEEP"
	}))
	out, err := c.Clone(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["SKIP"]; ok {
		t.Error("expected SKIP to be excluded")
	}
	if out["KEEP"] != "yes" {
		t.Errorf("expected KEEP=yes, got %q", out["KEEP"])
	}
}

func TestClone_WithTransform_ModifiesOutput(t *testing.T) {
	src := map[string]string{"key": "value"}
	c := NewCloner(WithCloneTransform(func(k, v string) (string, string, error) {
		return strings.ToUpper(k), strings.ToUpper(v), nil
	}))
	out, err := c.Clone(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "VALUE" {
		t.Errorf("expected KEY=VALUE, got %v", out)
	}
}

func TestClone_TransformError_ReturnsError(t *testing.T) {
	src := map[string]string{"bad": "val"}
	c := NewCloner(WithCloneTransform(func(k, v string) (string, string, error) {
		return "", "", errors.New("transform failed")
	}))
	_, err := c.Clone(src)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMustClone_PanicsOnError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic, got none")
		}
	}()
	src := map[string]string{"k": "v"}
	c := NewCloner(WithCloneTransform(func(k, v string) (string, string, error) {
		return "", "", errors.New("boom")
	}))
	c.MustClone(src)
}
