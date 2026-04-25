package env

import (
	"errors"
	"strings"
	"testing"
)

func upperAllStep(m map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = strings.ToUpper(v)
	}
	return out, nil
}

func prefixKeyStep(prefix string) ChainStep {
	return func(m map[string]string) (map[string]string, error) {
		out := make(map[string]string, len(m))
		for k, v := range m {
			out[prefix+k] = v
		}
		return out, nil
	}
}

func failingStep(m map[string]string) (map[string]string, error) {
	return nil, errors.New("step failed")
}

func TestChain_NoSteps_ReturnsSameMap(t *testing.T) {
	input := map[string]string{"KEY": "value"}
	chain := NewChain()
	out, err := chain.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "value" {
		t.Errorf("expected value to be preserved, got %q", out["KEY"])
	}
}

func TestChain_SingleStep_TransformsMap(t *testing.T) {
	input := map[string]string{"key": "hello"}
	chain := NewChain(upperAllStep)
	out, err := chain.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["key"] != "HELLO" {
		t.Errorf("expected HELLO, got %q", out["key"])
	}
}

func TestChain_MultipleSteps_Composed(t *testing.T) {
	input := map[string]string{"token": "abc"}
	chain := NewChain(upperAllStep, prefixKeyStep("APP_"))
	out, err := chain.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_token"] != "ABC" {
		t.Errorf("expected APP_token=ABC, got %v", out)
	}
}

func TestChain_StopsOnFirstError(t *testing.T) {
	input := map[string]string{"k": "v"}
	chain := NewChain(failingStep, upperAllStep)
	_, err := chain.Apply(input)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "step failed" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestChain_Append_AddsSteps(t *testing.T) {
	base := NewChain(upperAllStep)
	extended := base.Append(prefixKeyStep("X_"))
	if base.Len() != 1 {
		t.Errorf("base chain should be unchanged, got len %d", base.Len())
	}
	if extended.Len() != 2 {
		t.Errorf("expected extended chain len 2, got %d", extended.Len())
	}
	input := map[string]string{"key": "val"}
	out, err := extended.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["X_key"] != "VAL" {
		t.Errorf("expected X_key=VAL, got %v", out)
	}
}

func TestChain_Len_ReturnsStepCount(t *testing.T) {
	chain := NewChain(upperAllStep, upperAllStep, upperAllStep)
	if chain.Len() != 3 {
		t.Errorf("expected 3, got %d", chain.Len())
	}
}
