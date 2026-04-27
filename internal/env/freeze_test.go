package env

import (
	"testing"
)

func TestFreeze_RetainsAllKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	fm := Freeze(src)
	if fm.Len() != 3 {
		t.Fatalf("expected 3 keys, got %d", fm.Len())
	}
}

func TestFreeze_IsolatesOriginal(t *testing.T) {
	src := map[string]string{"KEY": "original"}
	fm := Freeze(src)
	src["KEY"] = "mutated"
	v, _ := fm.Get("KEY")
	if v != "original" {
		t.Fatalf("expected frozen value %q, got %q", "original", v)
	}
}

func TestGet_ExistingKey_ReturnsValue(t *testing.T) {
	fm := Freeze(map[string]string{"TOKEN": "secret"})
	v, ok := fm.Get("TOKEN")
	if !ok || v != "secret" {
		t.Fatalf("expected (secret, true), got (%q, %v)", v, ok)
	}
}

func TestGet_MissingKey_ReturnsFalse(t *testing.T) {
	fm := Freeze(map[string]string{})
	_, ok := fm.Get("MISSING")
	if ok {
		t.Fatal("expected ok=false for missing key")
	}
}

func TestSet_ReturnsFrozenError(t *testing.T) {
	fm := Freeze(map[string]string{"A": "1"})
	err := fm.Set("A", "2")
	if err == nil {
		t.Fatal("expected error on Set, got nil")
	}
	if !isErrFrozen(err) {
		t.Fatalf("expected ErrFrozen, got %v", err)
	}
}

func TestDelete_ReturnsFrozenError(t *testing.T) {
	fm := Freeze(map[string]string{"A": "1"})
	err := fm.Delete("A")
	if err == nil {
		t.Fatal("expected error on Delete, got nil")
	}
	if !isErrFrozen(err) {
		t.Fatalf("expected ErrFrozen, got %v", err)
	}
}

func TestKeys_ReturnsSorted(t *testing.T) {
	fm := Freeze(map[string]string{"Z": "z", "A": "a", "M": "m"})
	keys := fm.Keys()
	expected := []string{"A", "M", "Z"}
	for i, k := range keys {
		if k != expected[i] {
			t.Fatalf("keys[%d]: expected %q, got %q", i, expected[i], k)
		}
	}
}

func TestToMap_ReturnsMutableCopy(t *testing.T) {
	fm := Freeze(map[string]string{"X": "10"})
	copy := fm.ToMap()
	copy["X"] = "mutated"
	v, _ := fm.Get("X")
	if v != "10" {
		t.Fatal("ToMap copy mutated the frozen map")
	}
}

func TestFreeze_EmptyMap_ZeroLen(t *testing.T) {
	fm := Freeze(map[string]string{})
	if fm.Len() != 0 {
		t.Fatalf("expected Len=0, got %d", fm.Len())
	}
}

// isErrFrozen checks whether the error wraps ErrFrozen.
func isErrFrozen(err error) bool {
	return err != nil && errors.Is(err, ErrFrozen)
}
