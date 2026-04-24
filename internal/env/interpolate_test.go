package env

import (
	"os"
	"testing"
)

func TestInterpolate_NoPlaceholder(t *testing.T) {
	i := NewInterpolator(map[string]string{"FOO": "bar"}, false)
	out, err := i.Interpolate("hello world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "hello world" {
		t.Errorf("expected %q, got %q", "hello world", out)
	}
}

func TestInterpolate_ResolvesFromSecrets(t *testing.T) {
	i := NewInterpolator(map[string]string{"DB_PASS": "s3cr3t"}, false)
	out, err := i.Interpolate("pass=${DB_PASS}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "pass=s3cr3t" {
		t.Errorf("expected %q, got %q", "pass=s3cr3t", out)
	}
}

func TestInterpolate_UsesDefault(t *testing.T) {
	i := NewInterpolator(map[string]string{}, false)
	out, err := i.Interpolate("${MISSING:-fallback}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "fallback" {
		t.Errorf("expected %q, got %q", "fallback", out)
	}
}

func TestInterpolate_EmptyDefault_AllowedWhenSyntaxPresent(t *testing.T) {
	i := NewInterpolator(map[string]string{}, false)
	out, err := i.Interpolate("${MISSING:-}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestInterpolate_UnresolvedNoDefault_ReturnsError(t *testing.T) {
	i := NewInterpolator(map[string]string{}, false)
	_, err := i.Interpolate("${UNDEFINED}")
	if err == nil {
		t.Fatal("expected error for unresolved placeholder")
	}
}

func TestInterpolate_FallbackToEnv(t *testing.T) {
	t.Setenv("VAULTPIPE_TEST_KEY", "envval")
	i := NewInterpolator(map[string]string{}, true)
	out, err := i.Interpolate("${VAULTPIPE_TEST_KEY}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "envval" {
		t.Errorf("expected %q, got %q", "envval", out)
	}
}

func TestInterpolate_SecretsPreferredOverEnv(t *testing.T) {
	os.Setenv("VAULTPIPE_PREF", "from-env")
	t.Cleanup(func() { os.Unsetenv("VAULTPIPE_PREF") })
	i := NewInterpolator(map[string]string{"VAULTPIPE_PREF": "from-secret"}, true)
	out, err := i.Interpolate("${VAULTPIPE_PREF}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "from-secret" {
		t.Errorf("expected %q, got %q", "from-secret", out)
	}
}

func TestInterpolateMap_AllResolved(t *testing.T) {
	i := NewInterpolator(map[string]string{"HOST": "localhost", "PORT": "5432"}, false)
	in := map[string]string{
		"DSN": "postgres://${HOST}:${PORT}/db",
		"PLAIN": "nochange",
	}
	out, err := i.InterpolateMap(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DSN"] != "postgres://localhost:5432/db" {
		t.Errorf("unexpected DSN: %q", out["DSN"])
	}
	if out["PLAIN"] != "nochange" {
		t.Errorf("unexpected PLAIN: %q", out["PLAIN"])
	}
}

func TestInterpolateMap_ErrorPropagates(t *testing.T) {
	i := NewInterpolator(map[string]string{}, false)
	in := map[string]string{"KEY": "${UNRESOLVED}"}
	_, err := i.InterpolateMap(in)
	if err == nil {
		t.Fatal("expected error to propagate from InterpolateMap")
	}
}
