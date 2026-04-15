package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "vaultpipe-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTemp(t, `
vault_address: "http://127.0.0.1:8200"
secrets:
  - path: myapp/db
    mount: secret
    env_map:
      password: DB_PASSWORD
      user: DB_USER
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VaultAddress != "http://127.0.0.1:8200" {
		t.Errorf("expected vault_address to be set, got %q", cfg.VaultAddress)
	}
	if len(cfg.Secrets) != 1 {
		t.Fatalf("expected 1 secret mapping, got %d", len(cfg.Secrets))
	}
	if cfg.Secrets[0].EnvMap["password"] != "DB_PASSWORD" {
		t.Errorf("env_map not parsed correctly")
	}
}

func TestLoad_DefaultMount(t *testing.T) {
	path := writeTemp(t, `
secrets:
  - path: myapp/api
    env_map:
      key: API_KEY
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Secrets[0].Mount != "secret" {
		t.Errorf("expected default mount 'secret', got %q", cfg.Secrets[0].Mount)
	}
}

func TestLoad_EmptyPath(t *testing.T) {
	_, err := Load("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_NoSecrets(t *testing.T) {
	path := writeTemp(t, `vault_address: "http://127.0.0.1:8200"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing secrets")
	}
}

func TestLoad_EmptySecretPath(t *testing.T) {
	path := writeTemp(t, `
secrets:
  - path: ""
    mount: secret
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for empty secret path")
	}
}
