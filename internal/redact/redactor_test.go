package redact_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/redact"
)

func TestRedact_ReplacesKnownSecret(t *testing.T) {
	r := redact.New(map[string]string{"KEY": "s3cr3t"})
	got := r.Redact("value is s3cr3t here")
	if got != "value is ***REDACTED*** here" {
		t.Fatalf("unexpected: %s", got)
	}
}

func TestRedact_NoMatch_Unchanged(t *testing.T) {
	r := redact.New(map[string]string{"KEY": "s3cr3t"})
	got := r.Redact("nothing sensitive")
	if got != "nothing sensitive" {
		t.Fatalf("unexpected: %s", got)
	}
}

func TestRedact_MultipleSecrets(t *testing.T) {
	r := redact.New(map[string]string{"A": "alpha", "B": "beta"})
	got := r.Redact("alpha and beta")
	if got != "***REDACTED*** and ***REDACTED***" {
		t.Fatalf("unexpected: %s", got)
	}
}

func TestRedact_EmptySecretValue_Skipped(t *testing.T) {
	r := redact.New(map[string]string{"EMPTY": ""})
	got := r.Redact("some string")
	if got != "some string" {
		t.Fatalf("unexpected: %s", got)
	}
}

func TestAdd_RegistersNewSecret(t *testing.T) {
	r := redact.New(nil)
	r.Add("newtoken")
	got := r.Redact("token=newtoken")
	if got != "token=***REDACTED***" {
		t.Fatalf("unexpected: %s", got)
	}
}

func TestRedactMap_MasksAllValues(t *testing.T) {
	m := map[string]string{"DB_PASS": "hunter2", "API_KEY": "abc123"}
	out := redact.RedactMap(m)
	for k, v := range out {
		if v != "***REDACTED***" {
			t.Fatalf("key %s not redacted: %s", k, v)
		}
	}
	if len(out) != len(m) {
		t.Fatalf("length mismatch")
	}
}
