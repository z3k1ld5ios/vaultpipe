package token_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpipe/internal/token"
)

func TestResolve_StaticToken(t *testing.T) {
	p := token.NewProvider("s.staticABC", "", "")
	tok, err := p.Resolve(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Value != "s.staticABC" {
		t.Errorf("got %q, want %q", tok.Value, "s.staticABC")
	}
	if tok.Source != token.SourceStatic {
		t.Errorf("expected source %q, got %q", token.SourceStatic, tok.Source)
	}
}

func TestResolve_EnvFallback(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "s.envToken")
	p := token.NewProvider("", "", "")
	tok, err := p.Resolve(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Value != "s.envToken" {
		t.Errorf("got %q, want %q", tok.Value, "s.envToken")
	}
	if tok.Source != token.SourceEnv {
		t.Errorf("expected source %q, got %q", token.SourceEnv, tok.Source)
	}
}

func TestResolve_CustomEnvVar(t *testing.T) {
	t.Setenv("MY_TOKEN", "s.customEnv")
	p := token.NewProvider("", "", "MY_TOKEN")
	tok, err := p.Resolve(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Value != "s.customEnv" {
		t.Errorf("got %q, want %q", tok.Value, "s.customEnv")
	}
}

func TestResolve_FileToken(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, ".vault-token")
	if err := os.WriteFile(f, []byte("  s.fileToken\n"), 0600); err != nil {
		t.Fatal(err)
	}
	p := token.NewProvider("", f, "VAULT_TOKEN_NONE")
	tok, err := p.Resolve(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Value != "s.fileToken" {
		t.Errorf("got %q, want %q", tok.Value, "s.fileToken")
	}
	if tok.Source != token.SourceFile {
		t.Errorf("expected source %q, got %q", token.SourceFile, tok.Source)
	}
}

func TestResolve_EmptyFile_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, ".vault-token")
	if err := os.WriteFile(f, []byte("   \n"), 0600); err != nil {
		t.Fatal(err)
	}
	p := token.NewProvider("", f, "VAULT_TOKEN_NONE")
	_, err := p.Resolve(context.Background())
	if err == nil {
		t.Fatal("expected error for empty token file")
	}
}

func TestResolve_NoSource_ReturnsError(t *testing.T) {
	p := token.NewProvider("", "", "VAULT_TOKEN_NONE")
	_, err := p.Resolve(context.Background())
	if err == nil {
		t.Fatal("expected error when no token source configured")
	}
}
