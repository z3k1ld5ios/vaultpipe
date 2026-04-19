package env

import (
	"os"
	"testing"
)

func TestExpand_LiteralString(t *testing.T) {
	e := NewExpander(map[string]string{}, false)
	got, err := e.Expand("hello world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "hello world" {
		t.Errorf("expected %q, got %q", "hello world", got)
	}
}

func TestExpand_ResolvesFromSecrets(t *testing.T) {
	e := NewExpander(map[string]string{"DB_PASS": "s3cr3t"}, false)
	got, err := e.Expand("pass=${DB_PASS}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "pass=s3cr3t" {
		t.Errorf("got %q", got)
	}
}

func TestExpand_FallbackToOS(t *testing.T) {
	os.Setenv("_TEST_EXPAND_VAR", "fromOS")
	defer os.Unsetenv("_TEST_EXPAND_VAR")

	e := NewExpander(map[string]string{}, true)
	got, err := e.Expand("${_TEST_EXPAND_VAR}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "fromOS" {
		t.Errorf("got %q", got)
	}
}

func TestExpand_UnresolvedReturnsError(t *testing.T) {
	e := NewExpander(map[string]string{}, false)
	_, err := e.Expand("${MISSING_VAR}")
	if err == nil {
		t.Fatal("expected error for unresolved variable")
	}
}

func TestExpandMap_AllResolved(t *testing.T) {
	secrets := map[string]string{"HOST": "localhost", "PORT": "5432"}
	e := NewExpander(secrets, false)
	input := map[string]string{
		"DSN": "postgres://${HOST}:${PORT}/db",
		"RAW": "static",
	}
	out, err := e.ExpandMap(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DSN"] != "postgres://localhost:5432/db" {
		t.Errorf("DSN: got %q", out["DSN"])
	}
	if out["RAW"] != "static" {
		t.Errorf("RAW: got %q", out["RAW"])
	}
}

func TestExpandMap_ErrorOnMissing(t *testing.T) {
	e := NewExpander(map[string]string{}, false)
	_, err := e.ExpandMap(map[string]string{"KEY": "${NOPE}"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestContainsReference(t *testing.T) {
	if !ContainsReference("${FOO}") {
		t.Error("expected true for ${FOO}")
	}
	if ContainsReference("plain") {
		t.Error("expected false for plain string")
	}
}
