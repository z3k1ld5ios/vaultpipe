package env

import (
	"testing"
)

func secrets() map[string]string {
	return map[string]string{
		"db_password": "s3cr3t",
		"api_key":     "abc123",
	}
}

func TestResolve_DirectKey(t *testing.T) {
	r := NewResolver(secrets())
	out, err := r.Resolve(map[string]string{"DB_PASS": "db_password"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_PASS"] != "s3cr3t" {
		t.Errorf("expected s3cr3t, got %q", out["DB_PASS"])
	}
}

func TestResolve_DefaultValue(t *testing.T) {
	r := NewResolver(secrets())
	out, err := r.Resolve(map[string]string{"MISSING": "not_there:fallback"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["MISSING"] != "fallback" {
		t.Errorf("expected fallback, got %q", out["MISSING"])
	}
}

func TestResolve_MissingNoDefault(t *testing.T) {
	r := NewResolver(secrets())
	_, err := r.Resolve(map[string]string{"X": "not_there"})
	if err == nil {
		t.Fatal("expected error for missing key without default")
	}
}

func TestResolve_MultipleKeys(t *testing.T) {
	r := NewResolver(secrets())
	out, err := r.Resolve(map[string]string{
		"DB_PASS": "db_password",
		"API_KEY": "api_key",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
}

func TestResolve_EmptyMappings(t *testing.T) {
	r := NewResolver(secrets())
	out, err := r.Resolve(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty result, got %d entries", len(out))
	}
}

func TestResolve_DefaultColonInValue(t *testing.T) {
	// default value itself may contain colons — only first colon is separator
	r := NewResolver(map[string]string{})
	out, err := r.Resolve(map[string]string{"ADDR": "host:localhost:8080"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["ADDR"] != "localhost:8080" {
		t.Errorf("expected localhost:8080, got %q", out["ADDR"])
	}
}
