package env

import (
	"testing"
)

func TestApply_FirstWins_DefaultBehaviour(t *testing.T) {
	d := NewDeduplicator()
	a := map[string]string{"FOO": "first", "BAR": "a"}
	b := map[string]string{"FOO": "second", "BAZ": "b"}

	result := d.Apply(a, b)

	if result["FOO"] != "first" {
		t.Errorf("expected first-wins: got %q", result["FOO"])
	}
	if result["BAR"] != "a" {
		t.Errorf("expected BAR=a: got %q", result["BAR"])
	}
	if result["BAZ"] != "b" {
		t.Errorf("expected BAZ=b: got %q", result["BAZ"])
	}
}

func TestApply_LastWins_OverridesEarlier(t *testing.T) {
	d := NewDeduplicator(WithLastWins())
	a := map[string]string{"FOO": "first"}
	b := map[string]string{"FOO": "second"}

	result := d.Apply(a, b)

	if result["FOO"] != "second" {
		t.Errorf("expected last-wins: got %q", result["FOO"])
	}
}

func TestApply_SingleSource_ReturnsCopy(t *testing.T) {
	d := NewDeduplicator()
	src := map[string]string{"A": "1", "B": "2"}

	result := d.Apply(src)

	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	src["A"] = "mutated"
	if result["A"] != "1" {
		t.Error("Apply should not share underlying map with source")
	}
}

func TestApply_EmptySources_ReturnsEmpty(t *testing.T) {
	d := NewDeduplicator()
	result := d.Apply()
	if len(result) != 0 {
		t.Errorf("expected empty map, got %d keys", len(result))
	}
}

func TestApply_ThreeSources_FirstWins(t *testing.T) {
	d := NewDeduplicator()
	a := map[string]string{"X": "alpha"}
	b := map[string]string{"X": "beta"}
	c := map[string]string{"X": "gamma", "Y": "delta"}

	result := d.Apply(a, b, c)

	if result["X"] != "alpha" {
		t.Errorf("expected alpha, got %q", result["X"])
	}
	if result["Y"] != "delta" {
		t.Errorf("expected delta, got %q", result["Y"])
	}
}

func TestKeys_ReturnsSortedUniqueKeys(t *testing.T) {
	m := map[string]string{"ZEBRA": "z", "APPLE": "a", "MANGO": "m"}
	keys := Keys(m)

	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	expected := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("keys[%d]: expected %q, got %q", i, expected[i], k)
		}
	}
}

func TestKeys_EmptyMap_ReturnsEmpty(t *testing.T) {
	keys := Keys(map[string]string{})
	if len(keys) != 0 {
		t.Errorf("expected empty slice, got %d elements", len(keys))
	}
}
