package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultpipe/internal/output"
)

func TestNewFormatter_DefaultsToText(t *testing.T) {
	f := output.NewFormatter(nil, "")
	if f == nil {
		t.Fatal("expected non-nil formatter")
	}
}

func TestPrintSecretKeys_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatText)

	err := f.PrintSecretKeys("secret/myapp", []string{"DB_PASS", "API_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "secret/myapp") {
		t.Errorf("expected path in output, got: %s", out)
	}
	if !strings.Contains(out, "DB_PASS") {
		t.Errorf("expected DB_PASS in output, got: %s", out)
	}
}

func TestPrintSecretKeys_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatJSON)

	err := f.PrintSecretKeys("secret/myapp", []string{"TOKEN"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `"path"`) {
		t.Errorf("expected JSON path field, got: %s", out)
	}
	if !strings.Contains(out, `"TOKEN"`) {
		t.Errorf("expected TOKEN in JSON output, got: %s", out)
	}
}

func TestPrintEnvPreview_MasksValues(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatText)

	entries := map[string]string{
		"SECRET_KEY": "super-secret-value",
		"DB_PASS":    "hunter2",
	}

	err := f.PrintEnvPreview(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "super-secret-value") {
		t.Error("secret value must not appear in preview output")
	}
	if strings.Contains(out, "hunter2") {
		t.Error("secret value must not appear in preview output")
	}
	if !strings.Contains(out, "***") {
		t.Error("expected masked placeholder in output")
	}
}

func TestPrintEnvPreview_JSONMasksValues(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatJSON)

	entries := map[string]string{"API_TOKEN": "plaintext"}

	err := f.PrintEnvPreview(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "plaintext") {
		t.Error("raw value must not appear in JSON preview")
	}
	if !strings.Contains(out, "***") {
		t.Error("expected masked value in JSON preview")
	}
}
