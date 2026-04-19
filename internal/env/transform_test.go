package env

import (
	"testing"
)

func TestApply_Upper(t *testing.T) {
	tr := NewTransformer()
	got, err := tr.Apply("upper", "hello")
	if err != nil || got != "HELLO" {
		t.Fatalf("expected HELLO, got %q err %v", got, err)
	}
}

func TestApply_Lower(t *testing.T) {
	tr := NewTransformer()
	got, err := tr.Apply("lower", "WORLD")
	if err != nil || got != "world" {
		t.Fatalf("expected world, got %q err %v", got, err)
	}
}

func TestApply_Trim(t *testing.T) {
	tr := NewTransformer()
	got, err := tr.Apply("trim", "  spaces  ")
	if err != nil || got != "spaces" {
		t.Fatalf("expected 'spaces', got %q err %v", got, err)
	}
}

func TestApply_UnknownTransform_ReturnsError(t *testing.T) {
	tr := NewTransformer()
	_, err := tr.Apply("base64", "value")
	if err == nil {
		t.Fatal("expected error for unknown transform")
	}
}

func TestApply_CustomTransform(t *testing.T) {
	tr := NewTransformer()
	tr.Register("reverse", func(v string) (string, error) {
		r := []rune(v)
		for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
			r[i], r[j] = r[j], r[i]
		}
		return string(r), nil
	})
	got, err := tr.Apply("reverse", "abc")
	if err != nil || got != "cba" {
		t.Fatalf("expected cba, got %q err %v", got, err)
	}
}

func TestApplyMap_TransformsSelectedKeys(t *testing.T) {
	tr := NewTransformer()
	secrets := map[string]string{"DB_USER": "admin", "DB_PASS": "Secret123", "PORT": "5432"}
	out, err := tr.ApplyMap(secrets, "lower", []string{"DB_USER", "DB_PASS"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_USER"] != "admin" {
		t.Errorf("expected admin, got %q", out["DB_USER"])
	}
	if out["DB_PASS"] != "secret123" {
		t.Errorf("expected secret123, got %q", out["DB_PASS"])
	}
	if out["PORT"] != "5432" {
		t.Errorf("PORT should be unchanged, got %q", out["PORT"])
	}
}

func TestApplyMap_MissingKey_Skipped(t *testing.T) {
	tr := NewTransformer()
	secrets := map[string]string{"A": "hello"}
	out, err := tr.ApplyMap(secrets, "upper", []string{"A", "MISSING"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "HELLO" {
		t.Errorf("expected HELLO, got %q", out["A"])
	}
}

func TestApplyMap_DoesNotMutateInput(t *testing.T) {
	tr := NewTransformer()
	secrets := map[string]string{"KEY": "value"}
	_, _ = tr.ApplyMap(secrets, "upper", []string{"KEY"})
	if secrets["KEY"] != "value" {
		t.Error("input map was mutated")
	}
}
