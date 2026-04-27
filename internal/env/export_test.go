package env

import (
	"strings"
	"testing"
)

func TestExport_ShellFormat_Quoted(t *testing.T) {
	e := NewExporter(FormatShell, true)
	out, err := e.Export(map[string]string{"FOO": "bar", "BAZ": "qux"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `export BAZ="qux"`) {
		t.Errorf("expected shell export for BAZ, got:\n%s", out)
	}
	if !strings.Contains(out, `export FOO="bar"`) {
		t.Errorf("expected shell export for FOO, got:\n%s", out)
	}
}

func TestExport_ShellFormat_Unquoted(t *testing.T) {
	e := NewExporter(FormatShell, false)
	out, err := e.Export(map[string]string{"KEY": "val"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export KEY=val") {
		t.Errorf("expected unquoted shell export, got: %s", out)
	}
}

func TestExport_DotenvFormat_Quoted(t *testing.T) {
	e := NewExporter(FormatDotenv, true)
	out, err := e.Export(map[string]string{"DB_PASS": "s3cr3t"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `DB_PASS="s3cr3t"`) {
		t.Errorf("expected dotenv line, got: %s", out)
	}
}

func TestExport_DotenvFormat_Unquoted(t *testing.T) {
	e := NewExporter(FormatDotenv, false)
	out, err := e.Export(map[string]string{"A": "1", "B": "2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "A=1") || !strings.Contains(out, "B=2") {
		t.Errorf("expected unquoted dotenv lines, got: %s", out)
	}
}

func TestExport_JSONFormat(t *testing.T) {
	e := NewExporter(FormatJSON, false)
	out, err := e.Export(map[string]string{"X": "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != `{"X":"hello"}` {
		t.Errorf("expected JSON object, got: %s", out)
	}
}

func TestExport_JSONFormat_MultipleKeys_Sorted(t *testing.T) {
	e := NewExporter(FormatJSON, false)
	out, err := e.Export(map[string]string{"Z": "last", "A": "first"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := `{"A":"first","Z":"last"}`
	if out != expected {
		t.Errorf("expected %s, got %s", expected, out)
	}
}

func TestExport_NilMap_ReturnsError(t *testing.T) {
	e := NewExporter(FormatDotenv, true)
	_, err := e.Export(nil)
	if err == nil {
		t.Error("expected error for nil map, got nil")
	}
}

func TestExport_EmptyMap_ReturnsEmptyString(t *testing.T) {
	e := NewExporter(FormatShell, true)
	out, err := e.Export(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty string, got: %q", out)
	}
}

func TestExport_UnknownFormat_ReturnsError(t *testing.T) {
	e := NewExporter(ExportFormat(99), false)
	_, err := e.Export(map[string]string{"K": "v"})
	if err == nil {
		t.Error("expected error for unknown format")
	}
}
