package env

import (
	"strings"
	"testing"
)

// TestExport_WithPrefixFilter_OnlyExportsPrefixedKeys verifies that Exporter
// composes correctly with PrefixFilter: only keys surviving the filter are
// rendered in the output.
func TestExport_WithPrefixFilter_OnlyExportsPrefixedKeys(t *testing.T) {
	raw := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DEBUG":    "true",
	}

	pf := NewPrefixFilter("APP_")
	filtered := pf.Apply(raw)

	exporter := NewExporter(FormatDotenv, false)
	out, err := exporter.Export(filtered)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(out, "DEBUG") {
		t.Errorf("DEBUG should have been filtered out, got:\n%s", out)
	}
	if !strings.Contains(out, "APP_HOST=localhost") {
		t.Errorf("APP_HOST missing from output:\n%s", out)
	}
	if !strings.Contains(out, "APP_PORT=8080") {
		t.Errorf("APP_PORT missing from output:\n%s", out)
	}
}

// TestExport_WithOverride_ReflectsOverriddenValues ensures that when Override
// is applied before Export the final output contains the overridden values.
func TestExport_WithOverride_ReflectsOverriddenValues(t *testing.T) {
	base := map[string]string{
		"TOKEN": "old-token",
		"REGION": "us-east-1",
	}
	overrides := map[string]string{
		"TOKEN": "new-token",
	}

	ov := NewOverride(base)
	merged, err := ov.Apply(overrides)
	if err != nil {
		t.Fatalf("override failed: %v", err)
	}

	exporter := NewExporter(FormatShell, true)
	out, err := exporter.Export(merged)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	if !strings.Contains(out, `export TOKEN="new-token"`) {
		t.Errorf("expected overridden TOKEN, got:\n%s", out)
	}
	if !strings.Contains(out, `export REGION="us-east-1"`) {
		t.Errorf("expected REGION preserved, got:\n%s", out)
	}
}
