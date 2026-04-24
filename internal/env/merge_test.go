package env

import (
	"errors"
	"testing"
)

func TestMerge_SecretWins_OverridesBase(t *testing.T) {
	m := NewMerger(StrategySecretWins)
	base := map[string]string{"FOO": "base", "BAR": "keep"}
	secrets := map[string]string{"FOO": "secret"}
	out, err := m.Merge(base, secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "secret" {
		t.Errorf("expected FOO=secret, got %q", out["FOO"])
	}
	if out["BAR"] != "keep" {
		t.Errorf("expected BAR=keep, got %q", out["BAR"])
	}
}

func TestMerge_BaseWins_PreservesExisting(t *testing.T) {
	m := NewMerger(StrategyBaseWins)
	base := map[string]string{"FOO": "original"}
	secrets := map[string]string{"FOO": "override", "NEW": "value"}
	out, err := m.Merge(base, secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "original" {
		t.Errorf("expected FOO=original, got %q", out["FOO"])
	}
	if out["NEW"] != "value" {
		t.Errorf("expected NEW=value, got %q", out["NEW"])
	}
}

func TestMerge_ErrorStrategy_ConflictReturnsError(t *testing.T) {
	m := NewMerger(StrategyError)
	base := map[string]string{"CONFLICT": "a"}
	secrets := map[string]string{"CONFLICT": "b"}
	_, err := m.Merge(base, secrets)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var ce *MergeConflictError
	if !errors.As(err, &ce) {
		t.Fatalf("expected MergeConflictError, got %T", err)
	}
	if ce.Key != "CONFLICT" {
		t.Errorf("expected key=CONFLICT, got %q", ce.Key)
	}
}

func TestMerge_ErrorStrategy_NoConflict_Succeeds(t *testing.T) {
	m := NewMerger(StrategyError)
	base := map[string]string{"A": "1"}
	secrets := map[string]string{"B": "2"}
	out, err := m.Merge(base, secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestMerge_DoesNotMutateInputs(t *testing.T) {
	m := NewMerger(StrategySecretWins)
	base := map[string]string{"X": "base"}
	secrets := map[string]string{"X": "secret", "Y": "new"}
	_, _ = m.Merge(base, secrets)
	if base["X"] != "base" {
		t.Error("base map was mutated")
	}
	if _, ok := base["Y"]; ok {
		t.Error("base map was mutated with new key")
	}
}

func TestMerge_EmptyInputs_ReturnsEmptyMap(t *testing.T) {
	m := NewMerger(StrategySecretWins)
	out, err := m.Merge(map[string]string{}, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d keys", len(out))
	}
}
