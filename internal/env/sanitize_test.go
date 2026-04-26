package env_test

import (
	"testing"

	"github.com/your-org/vaultpipe/internal/env"
)

func TestSanitizeKey_Valid(t *testing.T) {
	key := env.SanitizeKey("MY_VAR")
	if key != "MY_VAR" {
		t.Fatalf("expected MY_VAR, got %s", key)
	}
}

func TestSanitizeKey_LowercaseConverted(t *testing.T) {
	key := env.SanitizeKey("my_var")
	if key != "MY_VAR" {
		t.Fatalf("expected MY_VAR, got %s", key)
	}
}

func TestSanitizeKey_HyphensReplaced(t *testing.T) {
	key := env.SanitizeKey("my-var-name")
	if key != "MY_VAR_NAME" {
		t.Fatalf("expected MY_VAR_NAME, got %s", key)
	}
}

func TestSanitizeKey_LeadingDigitPrefixed(t *testing.T) {
	key := env.SanitizeKey("1bad")
	if key != "_1BAD" {
		t.Fatalf("expected _1BAD, got %s", key)
	}
}

func TestSanitizeKey_SpecialCharsStripped(t *testing.T) {
	key := env.SanitizeKey("my.var!name")
	if key != "MY_VAR_NAME" {
		t.Fatalf("expected MY_VAR_NAME, got %s", key)
	}
}

func TestSanitizeKey_EmptyString(t *testing.T) {
	key := env.SanitizeKey("")
	if key != "" {
		t.Fatalf("expected empty string, got %s", key)
	}
}

func TestSanitizeKey_SpacesReplaced(t *testing.T) {
	key := env.SanitizeKey("my var name")
	if key != "MY_VAR_NAME" {
		t.Fatalf("expected MY_VAR_NAME, got %s", key)
	}
}

func TestSanitizeKey_DotsAndDashes(t *testing.T) {
	key := env.SanitizeKey("app.db-host")
	if key != "APP_DB_HOST" {
		t.Fatalf("expected APP_DB_HOST, got %s", key)
	}
}
