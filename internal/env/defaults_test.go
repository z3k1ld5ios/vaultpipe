package env

import (
	"testing"
)

func TestApplyDefaults_FillsMissingKeys(t *testing.T) {
	a := NewDefaultsApplier(map[string]string{"FOO": "default_foo", "BAR": "default_bar"})
	out := a.Apply(map[string]string{})
	if out["FOO"] != "default_foo" || out["BAR"] != "default_bar" {
		t.Fatalf("expected defaults to be applied, got %v", out)
	}
}

func TestApplyDefaults_DoesNotOverrideExisting(t *testing.T) {
	a := NewDefaultsApplier(map[string]string{"FOO": "default_foo"})
	out := a.Apply(map[string]string{"FOO": "real_foo"})
	if out["FOO"] != "real_foo" {
		t.Fatalf("expected existing key to be preserved, got %q", out["FOO"])
	}
}

func TestApplyDefaults_MergesBoth(t *testing.T) {
	a := NewDefaultsApplier(map[string]string{"FOO": "df", "BAR": "db"})
	out := a.Apply(map[string]string{"FOO": "real", "BAZ": "baz_val"})
	if out["FOO"] != "real" {
		t.Errorf("FOO should be real, got %q", out["FOO"])
	}
	if out["BAR"] != "db" {
		t.Errorf("BAR should be default db, got %q", out["BAR"])
	}
	if out["BAZ"] != "baz_val" {
		t.Errorf("BAZ should be baz_val, got %q", out["BAZ"])
	}
}

func TestApplyDefaults_EmptyDefaults_ReturnsBase(t *testing.T) {
	a := NewDefaultsApplier(map[string]string{})
	base := map[string]string{"X": "1"}
	out := a.Apply(base)
	if out["X"] != "1" || len(out) != 1 {
		t.Fatalf("expected base unchanged, got %v", out)
	}
}

func TestApplyDefaults_DoesNotMutateBase(t *testing.T) {
	a := NewDefaultsApplier(map[string]string{"NEW": "val"})
	base := map[string]string{"OLD": "old"}
	a.Apply(base)
	if _, ok := base["NEW"]; ok {
		t.Fatal("Apply must not mutate the base map")
	}
}

func TestKeys_ReturnsAllDefaultKeys(t *testing.T) {
	a := NewDefaultsApplier(map[string]string{"A": "1", "B": "2"})
	keys := a.Keys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
}
