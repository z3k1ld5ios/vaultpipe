package env_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/env"
)

// TestOverride_WithPrefixFilter verifies that Override composes correctly with
// PrefixFilter: first strip prefixes, then apply overrides.
func TestOverride_WithPrefixFilter_Composition(t *testing.T) {
	prefixed := map[string]string{
		"APP_DB_HOST": "db.internal",
		"APP_DB_PORT": "5432",
	}

	pf := env.NewPrefixFilter("APP_")
	stripped := pf.Strip(prefixed)

	o := env.NewOverride(map[string]string{"DB_PORT": "5433"})
	out := o.Apply(stripped)

	if out["DB_HOST"] != "db.internal" {
		t.Errorf("expected db.internal, got %q", out["DB_HOST"])
	}
	if out["DB_PORT"] != "5433" {
		t.Errorf("expected override 5433, got %q", out["DB_PORT"])
	}
}

// TestOverride_WithMapper verifies that Override can layer values on top of a
// mapped secret environment.
func TestOverride_WithMapper_Composition(t *testing.T) {
	secrets := map[string]string{
		"password": "s3cr3t",
		"host":     "vault-host",
	}

	m := env.NewMapper(map[string]string{
		"password": "DB_PASSWORD",
		"host":     "DB_HOST",
	}, "")

	mapped, err := m.Apply(secrets)
	if err != nil {
		t.Fatalf("mapper error: %v", err)
	}

	o := env.NewOverride(map[string]string{"DB_HOST": "localhost"})
	out := o.Apply(mapped)

	if out["DB_PASSWORD"] != "s3cr3t" {
		t.Errorf("unexpected password: %q", out["DB_PASSWORD"])
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected localhost override, got %q", out["DB_HOST"])
	}
}
