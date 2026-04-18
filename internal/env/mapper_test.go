package env

import (
	"testing"
)

func sampleSecretMap() map[string]string {
	return map[string]string{
		"db_password": "s3cr3t",
		"api_key":     "abc123",
		"debug":       "true",
	}
}

func TestApply_NoMappings_UsesAllKeysWithPrefix(t *testing.T) {
	m := NewMapper("APP_", nil)
	out, err := m.Apply(sampleSecretMap())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["APP_DB_PASSWORD"]; !ok {
		t.Error("expected APP_DB_PASSWORD in output")
	}
	if _, ok := out["APP_API_KEY"]; !ok {
		t.Error("expected APP_API_KEY in output")
	}
}

func TestApply_ExplicitMappings_OnlyMappedKeysIncluded(t *testing.T) {
	mappings := []Mapping{
		{SecretKey: "db_password", EnvKey: "DATABASE_PASS"},
	}
	m := NewMapper("", mappings)
	out, err := m.Apply(sampleSecretMap())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
	if out["DATABASE_PASS"] != "s3cr3t" {
		t.Errorf("unexpected value: %s", out["DATABASE_PASS"])
	}
}

func TestApply_MissingSecretKey_ReturnsError(t *testing.T) {
	mappings := []Mapping{
		{SecretKey: "nonexistent", EnvKey: "SOME_VAR"},
	}
	m := NewMapper("", mappings)
	_, err := m.Apply(sampleSecretMap())
	if err == nil {
		t.Fatal("expected error for missing secret key")
	}
}

func TestApply_NoPrefix_KeysUnchanged(t *testing.T) {
	m := NewMapper("", nil)
	out, err := m.Apply(map[string]string{"foo_bar": "val"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO_BAR"] != "val" {
		t.Errorf("expected FOO_BAR=val, got %v", out)
	}
}

func TestApply_EmptySecrets_ReturnsEmpty(t *testing.T) {
	m := NewMapper("PRE_", nil)
	out, err := m.Apply(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
