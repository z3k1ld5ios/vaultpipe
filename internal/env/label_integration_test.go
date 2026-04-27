package env

import (
	"strings"
	"testing"
)

// TestLabeler_WithOverride_LabelsPreservedAfterMerge verifies that labels
// attached before an override step are still accessible after the env map
// has been transformed by Override.
func TestLabeler_WithOverride_LabelsPreservedAfterMerge(t *testing.T) {
	base := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	overrides := map[string]string{
		"DB_HOST": "prod.db.internal",
	}

	l := NewLabeler()
	l.Set("DB_HOST", "source", "vault")
	l.Set("DB_PORT", "source", "static")

	o := NewOverride(base, overrides)
	result, err := o.Apply()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Labels should still be queryable after the env map has been merged.
	src, ok := l.Get("DB_HOST", "source")
	if !ok || src != "vault" {
		t.Fatalf("expected label 'vault', got ok=%v src=%q", ok, src)
	}

	// The merged value should reflect the override.
	if result["DB_HOST"] != "prod.db.internal" {
		t.Fatalf("expected override value, got %q", result["DB_HOST"])
	}
}

// TestLabeler_Annotate_WithPrefixFilter_OnlyAnnotatesFilteredKeys verifies
// that Annotate only produces annotations for keys that survive a prefix
// filter, and that unannotated keys remain unchanged.
func TestLabeler_Annotate_WithPrefixFilter_OnlyAnnotatesFilteredKeys(t *testing.T) {
	l := NewLabeler()
	l.Set("VAULT_TOKEN", "tier", "secret")
	l.Set("APP_NAME", "tier", "config")

	full := map[string]string{
		"VAULT_TOKEN": "hvs.abc",
		"APP_NAME":    "vaultpipe",
	}

	pf := NewPrefixFilter("VAULT_")
	filtered := pf.Apply(full)

	annotated := l.Annotate(filtered)

	if _, exists := annotated["APP_NAME"]; exists {
		t.Fatal("APP_NAME should have been filtered out before annotation")
	}

	v, ok := annotated["VAULT_TOKEN"]
	if !ok {
		t.Fatal("expected VAULT_TOKEN in annotated output")
	}
	if !strings.Contains(v, "tier=secret") {
		t.Fatalf("expected tier annotation, got %q", v)
	}
}
