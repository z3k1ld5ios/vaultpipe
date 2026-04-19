package env

import (
	"testing"
)

// TestTransform_WithOverride_Composition verifies that Transformer can be
// composed with Override to normalize values before merging overrides.
func TestTransform_WithOverride_Composition(t *testing.T) {
	tr := NewTransformer()
	base := map[string]string{"API_KEY": "  MySecret  ", "HOST": "localhost"}

	trimmed, err := tr.ApplyMap(base, "trim", []string{"API_KEY"})
	if err != nil {
		t.Fatalf("trim failed: %v", err)
	}

	ov := NewOverride(map[string]string{"HOST": "prod.example.com"})
	result := ov.Apply(trimmed)

	if result["API_KEY"] != "MySecret" {
		t.Errorf("expected trimmed API_KEY, got %q", result["API_KEY"])
	}
	if result["HOST"] != "prod.example.com" {
		t.Errorf("expected overridden HOST, got %q", result["HOST"])
	}
}

// TestTransform_WithDefaults_Composition verifies Transformer + DefaultsApplier.
func TestTransform_WithDefaults_Composition(t *testing.T) {
	tr := NewTransformer()
	secrets := map[string]string{"REGION": "US-EAST-1"}

	lowered, err := tr.ApplyMap(secrets, "lower", []string{"REGION"})
	if err != nil {
		t.Fatalf("lower failed: %v", err)
	}

	da := NewDefaultsApplier(map[string]string{"REGION": "us-west-2", "ENV": "production"})
	result := da.Apply(lowered)

	if result["REGION"] != "us-east-1" {
		t.Errorf("existing key should not be overridden, got %q", result["REGION"])
	}
	if result["ENV"] != "production" {
		t.Errorf("default ENV should be set, got %q", result["ENV"])
	}
}
