package process

import (
	"testing"
)

func TestRun_EmptyCommand(t *testing.T) {
	r := NewRunner([]string{})
	code, err := r.Run("", nil)
	if err == nil {
		t.Fatal("expected error for empty command, got nil")
	}
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestRun_SuccessfulCommand(t *testing.T) {
	env := []string{"PATH=/usr/bin:/bin", "HOME=/tmp"}
	r := NewRunner(env)
	code, err := r.Run("true", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}

func TestRun_FailingCommand(t *testing.T) {
	env := []string{"PATH=/usr/bin:/bin"}
	r := NewRunner(env)
	code, err := r.Run("false", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestRun_CommandWithArgs(t *testing.T) {
	env := []string{"PATH=/usr/bin:/bin"}
	r := NewRunner(env)
	code, err := r.Run("sh", []string{"-c", "exit 42"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 42 {
		t.Fatalf("expected exit code 42, got %d", code)
	}
}

func TestRun_EnvIsInjected(t *testing.T) {
	env := []string{"PATH=/usr/bin:/bin", "VAULTPIPE_TEST_VAR=hello"}
	r := NewRunner(env)
	// sh -c 'test "$VAULTPIPE_TEST_VAR" = "hello"' exits 0 on match
	code, err := r.Run("sh", []string{"-c", `test "$VAULTPIPE_TEST_VAR" = "hello"`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 0 {
		t.Fatalf("expected exit code 0 (env var present), got %d", code)
	}
}

func TestRun_InvalidBinary(t *testing.T) {
	r := NewRunner([]string{"PATH=/usr/bin:/bin"})
	_, err := r.Run("/nonexistent/binary", nil)
	if err == nil {
		t.Fatal("expected error for invalid binary, got nil")
	}
}
