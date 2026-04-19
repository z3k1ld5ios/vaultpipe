package env

import (
	"sort"
	"testing"
)

func TestApplyOverride_AddsNewKeys(t *testing.T) {
	base := map[string]string{"A": "1"}
	o := NewOverride(map[string]string{"B": "2"})
	out := o.Apply(base)
	if out["A"] != "1" || out["B"] != "2" {
		t.Fatalf("unexpected result: %v", out)
	}
}

func TestApplyOverride_OverridesExistingKeys(t *testing.T) {
	base := map[string]string{"A": "original"}
	o := NewOverride(map[string]string{"A": "overridden"})
	out := o.Apply(base)
	if out["A"] != "overridden" {
		t.Fatalf("expected overridden, got %q", out["A"])
	}
}

func TestApplyOverride_EmptyOverrides_ReturnsCopy(t *testing.T) {
	base := map[string]string{"X": "10"}
	o := NewOverride(nil)
	out := o.Apply(base)
	if out["X"] != "10" || len(out) != 1 {
		t.Fatalf("unexpected result: %v", out)
	}
}

func TestApplyOverride_EmptyBase_ReturnsOverrides(t *testing.T) {
	o := NewOverride(map[string]string{"K": "v"})
	out := o.Apply(map[string]string{})
	if out["K"] != "v" {
		t.Fatalf("expected K=v, got %v", out)
	}
}

func TestApplyOverride_DoesNotMutateBase(t *testing.T) {
	base := map[string]string{"A": "1"}
	o := NewOverride(map[string]string{"A": "99"})
	_ = o.Apply(base)
	if base["A"] != "1" {
		t.Fatal("base was mutated")
	}
}

func TestOverride_Keys_ReturnsAll(t *testing.T) {
	o := NewOverride(map[string]string{"X": "1", "Y": "2"})
	keys := o.Keys()
	sort.Strings(keys)
	if len(keys) != 2 || keys[0] != "X" || keys[1] != "Y" {
		t.Fatalf("unexpected keys: %v", keys)
	}
}

func TestOverride_Len(t *testing.T) {
	o := NewOverride(map[string]string{"A": "1", "B": "2", "C": "3"})
	if o.Len() != 3 {
		t.Fatalf("expected 3, got %d", o.Len())
	}
}
