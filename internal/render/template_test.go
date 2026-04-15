package render

import (
	"testing"
)

func TestRenderValue_NoTemplate(t *testing.T) {
	r := NewRenderer(map[string]string{"FOO": "bar"})
	got, err := r.RenderValue("plain-value")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "plain-value" {
		t.Errorf("expected 'plain-value', got %q", got)
	}
}

func TestRenderValue_ResolvesSecret(t *testing.T) {
	r := NewRenderer(map[string]string{"DB_PASS": "s3cr3t"})
	got, err := r.RenderValue(`{{ vault "DB_PASS" }}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "s3cr3t" {
		t.Errorf("expected 's3cr3t', got %q", got)
	}
}

func TestRenderValue_MissingKey(t *testing.T) {
	r := NewRenderer(map[string]string{})
	_, err := r.RenderValue(`{{ vault "MISSING" }}`)
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestRenderValue_Interpolated(t *testing.T) {
	r := NewRenderer(map[string]string{"HOST": "localhost"})
	got, err := r.RenderValue(`postgres://{{ vault "HOST" }}:5432/db`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "postgres://localhost:5432/db"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestRenderMap_AllResolved(t *testing.T) {
	r := NewRenderer(map[string]string{"TOKEN": "abc123", "USER": "admin"})
	input := map[string]string{
		"APP_TOKEN": `{{ vault "TOKEN" }}`,
		"APP_USER":  `{{ vault "USER" }}`,
		"STATIC":    "unchanged",
	}
	out, err := r.RenderMap(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_TOKEN"] != "abc123" {
		t.Errorf("APP_TOKEN: expected 'abc123', got %q", out["APP_TOKEN"])
	}
	if out["APP_USER"] != "admin" {
		t.Errorf("APP_USER: expected 'admin', got %q", out["APP_USER"])
	}
	if out["STATIC"] != "unchanged" {
		t.Errorf("STATIC: expected 'unchanged', got %q", out["STATIC"])
	}
}

func TestRenderMap_PropagatesError(t *testing.T) {
	r := NewRenderer(map[string]string{})
	input := map[string]string{
		"KEY": `{{ vault "NO_SUCH" }}`,
	}
	_, err := r.RenderMap(input)
	if err == nil {
		t.Fatal("expected error from RenderMap, got nil")
	}
}
