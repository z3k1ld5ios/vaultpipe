package env

import (
	"strings"
	"testing"
)

func TestSet_AddsLabel(t *testing.T) {
	l := NewLabeler()
	l.Set("DB_PASSWORD", "source", "vault")
	v, ok := l.Get("DB_PASSWORD", "source")
	if !ok {
		t.Fatal("expected label to be found")
	}
	if v != "vault" {
		t.Fatalf("expected 'vault', got %q", v)
	}
}

func TestSet_OverwritesDuplicateLabelKey(t *testing.T) {
	l := NewLabeler()
	l.Set("TOKEN", "tier", "low")
	l.Set("TOKEN", "tier", "high")
	v, _ := l.Get("TOKEN", "tier")
	if v != "high" {
		t.Fatalf("expected 'high', got %q", v)
	}
}

func TestGet_CaseInsensitiveEnvKey(t *testing.T) {
	l := NewLabeler()
	l.Set("api_key", "source", "vault")
	v, ok := l.Get("API_KEY", "source")
	if !ok || v != "vault" {
		t.Fatalf("expected case-insensitive match, got ok=%v v=%q", ok, v)
	}
}

func TestGet_MissingLabel_ReturnsFalse(t *testing.T) {
	l := NewLabeler()
	_, ok := l.Get("MISSING", "source")
	if ok {
		t.Fatal("expected false for missing label")
	}
}

func TestAll_ReturnsAllLabels(t *testing.T) {
	l := NewLabeler()
	l.Set("SECRET", "source", "vault")
	l.Set("SECRET", "tier", "critical")
	labels := l.All("SECRET")
	if len(labels) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(labels))
	}
}

func TestDelete_RemovesAllLabels(t *testing.T) {
	l := NewLabeler()
	l.Set("KEY", "source", "vault")
	l.Delete("KEY")
	if labels := l.All("KEY"); len(labels) != 0 {
		t.Fatalf("expected empty after delete, got %d labels", len(labels))
	}
}

func TestAnnotate_AppendsLabelSuffix(t *testing.T) {
	l := NewLabeler()
	l.Set("DB_PASS", "source", "vault")
	m := map[string]string{"DB_PASS": "s3cr3t"}
	annotated := l.Annotate(m)
	v := annotated["DB_PASS"]
	if !strings.Contains(v, "source=vault") {
		t.Fatalf("expected annotation in value, got %q", v)
	}
	if !strings.HasPrefix(v, "s3cr3t") {
		t.Fatalf("expected original value prefix, got %q", v)
	}
}

func TestAnnotate_NoLabels_ReturnsOriginalValue(t *testing.T) {
	l := NewLabeler()
	m := map[string]string{"PLAIN": "value"}
	out := l.Annotate(m)
	if out["PLAIN"] != "value" {
		t.Fatalf("expected unchanged value, got %q", out["PLAIN"])
	}
}

func TestAnnotate_DoesNotMutateInput(t *testing.T) {
	l := NewLabeler()
	l.Set("K", "x", "y")
	m := map[string]string{"K": "v"}
	l.Annotate(m)
	if m["K"] != "v" {
		t.Fatal("input map was mutated")
	}
}
