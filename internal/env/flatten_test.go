package env

import (
	"testing"
)

func TestFlatten_FlatInput(t *testing.T) {
	f := NewFlattener("__")
	out, err := f.Flatten(map[string]any{"host": "localhost", "port": "5432"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["host"] != "localhost" || out["port"] != "5432" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestFlatten_NestedMap(t *testing.T) {
	f := NewFlattener("__")
	out, err := f.Flatten(map[string]any{
		"db": map[string]any{"host": "db.local", "port": "5432"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["db__host"] != "db.local" {
		t.Errorf("expected db__host=db.local, got %q", out["db__host"])
	}
	if out["db__port"] != "5432" {
		t.Errorf("expected db__port=5432, got %q", out["db__port"])
	}
}

func TestFlatten_WithPrefix(t *testing.T) {
	f := NewFlattener("__").WithPrefix("APP")
	out, err := f.Flatten(map[string]any{"key": "val"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP__key"] != "val" {
		t.Errorf("expected APP__key=val, got %v", out)
	}
}

func TestFlatten_NonStringLeaf(t *testing.T) {
	f := NewFlattener("__")
	out, err := f.Flatten(map[string]any{"count": 42, "active": true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["count"] != "42" {
		t.Errorf("expected count=42, got %q", out["count"])
	}
	if out["active"] != "true" {
		t.Errorf("expected active=true, got %q", out["active"])
	}
}

func TestFlatten_NilLeaf(t *testing.T) {
	f := NewFlattener("__")
	out, err := f.Flatten(map[string]any{"empty": nil})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := out["empty"]; !ok || v != "" {
		t.Errorf("expected empty string for nil leaf, got %q", v)
	}
}

func TestFlatten_InvalidKey(t *testing.T) {
	f := NewFlattener("__")
	_, err := f.Flatten(map[string]any{"bad=key": "val"})
	if err == nil {
		t.Fatal("expected error for key containing '='")
	}
}

func TestFlatten_DefaultSeparator(t *testing.T) {
	f := NewFlattener("")
	out, err := f.Flatten(map[string]any{
		"a": map[string]any{"b": "c"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["a__b"] != "c" {
		t.Errorf("expected default separator __, got %v", out)
	}
}

func TestKeys_ReturnsSorted(t *testing.T) {
	f := NewFlattener("__")
	keys, err := f.Keys(map[string]any{"z": "1", "a": "2", "m": "3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if keys[0] != "a" || keys[1] != "m" || keys[2] != "z" {
		t.Errorf("expected sorted keys, got %v", keys)
	}
}
