package env

import (
	"errors"
	"strings"
	"testing"
)

func upperStep(m map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = strings.ToUpper(v)
	}
	return out, nil
}

func prefixStep(prefix string) StepFunc {
	return func(m map[string]string) (map[string]string, error) {
		out := make(map[string]string, len(m))
		for k, v := range m {
			out[prefix+k] = v
		}
		return out, nil
	}
}

func failStep(m map[string]string) (map[string]string, error) {
	return m, errors.New("step failed")
}

func TestPipeline_NoSteps_ReturnsCopy(t *testing.T) {
	src := map[string]string{"key": "value"}
	p := NewPipeline()
	out, err := p.Run(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["key"] != "value" {
		t.Errorf("expected value, got %q", out["key"])
	}
}

func TestPipeline_SingleStep_Transforms(t *testing.T) {
	src := map[string]string{"greeting": "hello"}
	p := NewPipeline().Append(upperStep)
	out, err := p.Run(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["greeting"] != "HELLO" {
		t.Errorf("expected HELLO, got %q", out["greeting"])
	}
}

func TestPipeline_MultipleSteps_AppliedInOrder(t *testing.T) {
	src := map[string]string{"x": "hello"}
	p := NewPipeline().Append(upperStep, prefixStep("VP_"))
	out, err := p.Run(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["VP_x"]; !ok {
		t.Errorf("expected key VP_x, got %v", out)
	}
	if out["VP_x"] != "HELLO" {
		t.Errorf("expected HELLO, got %q", out["VP_x"])
	}
}

func TestPipeline_StepError_HaltsExecution(t *testing.T) {
	called := false
	neverCalled := func(m map[string]string) (map[string]string, error) {
		called = true
		return m, nil
	}
	p := NewPipeline().Append(failStep, neverCalled)
	_, err := p.Run(map[string]string{"k": "v"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if called {
		t.Error("step after failure should not have been called")
	}
}

func TestPipeline_DoesNotMutateSource(t *testing.T) {
	src := map[string]string{"a": "original"}
	p := NewPipeline().Append(upperStep)
	_, _ = p.Run(src)
	if src["a"] != "original" {
		t.Errorf("source map was mutated")
	}
}

func TestPipeline_Len_ReturnsStepCount(t *testing.T) {
	p := NewPipeline().Append(upperStep, prefixStep("X_"))
	if p.Len() != 2 {
		t.Errorf("expected 2 steps, got %d", p.Len())
	}
}
