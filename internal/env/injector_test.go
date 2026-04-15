package env

import (
	"sort"
	"strings"
	"testing"
)

func TestMerge_AddsSecrets(t *testing.T) {
	base := []string{"HOME=/root", "PATH=/usr/bin"}
	inj := NewInjector(base)

	secrets := SecretMap{
		"DB_PASSWORD": "supersecret",
		"API_KEY":     "abc123",
	}

	result := inj.Merge(secrets)

	env := toMap(result)

	if env["HOME"] != "/root" {
		t.Errorf("expected HOME=/root, got %q", env["HOME"])
	}
	if env["DB_PASSWORD"] != "supersecret" {
		t.Errorf("expected DB_PASSWORD=supersecret, got %q", env["DB_PASSWORD"])
	}
	if env["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", env["API_KEY"])
	}
}

func TestMerge_SecretsOverrideBase(t *testing.T) {
	base := []string{"DB_PASSWORD=old", "PATH=/usr/bin"}
	inj := NewInjector(base)

	secrets := SecretMap{"DB_PASSWORD": "new"}
	result := inj.Merge(secrets)

	env := toMap(result)
	if env["DB_PASSWORD"] != "new" {
		t.Errorf("expected DB_PASSWORD=new, got %q", env["DB_PASSWORD"])
	}
}

func TestMerge_EmptySecrets(t *testing.T) {
	base := []string{"FOO=bar"}
	inj := NewInjector(base)

	result := inj.Merge(SecretMap{})

	if len(result) != 1 {
		t.Errorf("expected 1 entry, got %d", len(result))
	}
}

func TestMerge_EmptyBase(t *testing.T) {
	inj := NewInjector([]string{})
	secrets := SecretMap{"TOKEN": "xyz"}

	result := inj.Merge(secrets)
	env := toMap(result)

	if env["TOKEN"] != "xyz" {
		t.Errorf("expected TOKEN=xyz, got %q", env["TOKEN"])
	}
}

func TestMerge_NoDuplicateKeys(t *testing.T) {
	base := []string{"KEY=a", "KEY=b"}
	inj := NewInjector(base)

	result := inj.Merge(SecretMap{})

	sort.Strings(result)
	count := 0
	for _, e := range result {
		if strings.HasPrefix(e, "KEY=") {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1 KEY entry after dedup, got %d", count)
	}
}

// toMap converts a slice of KEY=VALUE strings into a map for easy assertion.
func toMap(env []string) map[string]string {
	m := make(map[string]string, len(env))
	for _, e := range env {
		k, v := splitEntry(e)
		m[k] = v
	}
	return m
}
