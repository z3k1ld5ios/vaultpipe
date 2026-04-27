package env_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/env"
	"github.com/yourusername/vaultpipe/internal/redact"
)

func TestRedactBridge_Apply_ReplacesSecretValues(t *testing.T) {
	r := redact.New()
	r.Add("s3cr3t")
	b := env.NewRedactBridge(r)

	input := map[string]string{"TOKEN": "s3cr3t", "HOST": "localhost"}
	out := b.Apply(input)

	if out["TOKEN"] == "s3cr3t" {
		t.Error("expected TOKEN value to be redacted")
	}
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST to be unchanged, got %q", out["HOST"])
	}
}

func TestRedactBridge_Apply_DoesNotMutateInput(t *testing.T) {
	r := redact.New()
	r.Add("topsecret")
	b := env.NewRedactBridge(r)

	input := map[string]string{"PASS": "topsecret"}
	_ = b.Apply(input)

	if input["PASS"] != "topsecret" {
		t.Error("Apply must not mutate the input map")
	}
}

func TestRedactBridge_RegisterSecrets_MasksAll(t *testing.T) {
	r := redact.New()
	b := env.NewRedactBridge(r)

	secrets := map[string]string{"db_pass": "hunter2", "api_key": "abc123"}
	b.RegisterSecrets(secrets)

	envMap := map[string]string{"DB_PASSWORD": "hunter2", "API_KEY": "abc123", "USER": "admin"}
	out := b.Apply(envMap)

	if out["DB_PASSWORD"] == "hunter2" {
		t.Error("expected DB_PASSWORD to be redacted")
	}
	if out["API_KEY"] == "abc123" {
		t.Error("expected API_KEY to be redacted")
	}
	if out["USER"] != "admin" {
		t.Errorf("expected USER to be unchanged, got %q", out["USER"])
	}
}

func TestRedactBridge_RegisterSecrets_SkipsEmptyValues(t *testing.T) {
	r := redact.New()
	b := env.NewRedactBridge(r)

	b.RegisterSecrets(map[string]string{"empty": ""})

	out := b.Apply(map[string]string{"KEY": ""})
	if out["KEY"] != "" {
		t.Errorf("empty value should not be redacted, got %q", out["KEY"])
	}
}

func TestRedactMap_ConvenienceFunction(t *testing.T) {
	secrets := map[string]string{"pw": "mysecret"}
	envMap := map[string]string{"PASSWORD": "mysecret", "MODE": "prod"}

	out := env.RedactMap(envMap, secrets)

	if out["PASSWORD"] == "mysecret" {
		t.Error("expected PASSWORD to be redacted by RedactMap")
	}
	if out["MODE"] != "prod" {
		t.Errorf("expected MODE unchanged, got %q", out["MODE"])
	}
}
