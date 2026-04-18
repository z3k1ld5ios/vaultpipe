package preflight_test

import (
	"errors"
	"testing"

	"github.com/your-org/vaultpipe/internal/preflight"
)

func TestRun_AllPass(t *testing.T) {
	r := preflight.NewRunner(
		preflight.Check{Name: "always-ok", Run: func() error { return nil }},
	)
	results, err := r.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Passed {
		t.Fatalf("expected 1 passing result, got %+v", results)
	}
}

func TestRun_OneFails(t *testing.T) {
	r := preflight.NewRunner(
		preflight.Check{Name: "ok", Run: func() error { return nil }},
		preflight.Check{Name: "bad", Run: func() error { return errors.New("boom") }},
	)
	results, err := r.Run()
	if err == nil {
		t.Fatal("expected error")
	}
	if results[1].Passed {
		t.Error("second check should have failed")
	}
	if results[1].Message != "boom" {
		t.Errorf("unexpected message: %s", results[1].Message)
	}
}

func TestRun_EmptyChecks(t *testing.T) {
	r := preflight.NewRunner()
	_, err := r.Run()
	if err != nil {
		t.Fatalf("empty runner should not error: %v", err)
	}
}

func TestRequireEnv_Set(t *testing.T) {
	env := func(k string) string { return "vault.example.com" }
	c := preflight.RequireEnv("VAULT_ADDR", env)
	if err := c.Run(); err != nil {
		t.Fatalf("expected no error: %v", err)
	}
}

func TestRequireEnv_Missing(t *testing.T) {
	env := func(k string) string { return "" }
	c := preflight.RequireEnv("VAULT_ADDR", env)
	if err := c.Run(); err == nil {
		t.Fatal("expected error for missing env var")
	}
}

func TestRequireNonEmptySecrets_NonEmpty(t *testing.T) {
	c := preflight.RequireNonEmptySecrets(map[string]string{"key": "val"})
	if err := c.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRequireNonEmptySecrets_Empty(t *testing.T) {
	c := preflight.RequireNonEmptySecrets(map[string]string{})
	if err := c.Run(); err == nil {
		t.Fatal("expected error for empty secrets")
	}
}
