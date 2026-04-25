package env

import (
	"testing"
)

func TestAsString_Found(t *testing.T) {
	tc := NewTypeCaster()
	m := map[string]string{"HOST": "localhost"}
	v, err := tc.AsString(m, "HOST")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "localhost" {
		t.Errorf("expected localhost, got %q", v)
	}
}

func TestAsString_Missing(t *testing.T) {
	tc := NewTypeCaster()
	_, err := tc.AsString(map[string]string{}, "MISSING")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestAsInt_Valid(t *testing.T) {
	tc := NewTypeCaster()
	m := map[string]string{"PORT": "8080"}
	n, err := tc.AsInt(m, "PORT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 8080 {
		t.Errorf("expected 8080, got %d", n)
	}
}

func TestAsInt_Invalid(t *testing.T) {
	tc := NewTypeCaster()
	m := map[string]string{"PORT": "abc"}
	_, err := tc.AsInt(m, "PORT")
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestAsBool_TrueVariants(t *testing.T) {
	tc := NewTypeCaster()
	for _, val := range []string{"1", "true", "yes", "on", "TRUE", "YES"} {
		m := map[string]string{"FLAG": val}
		b, err := tc.AsBool(m, "FLAG")
		if err != nil {
			t.Errorf("value %q: unexpected error: %v", val, err)
		}
		if !b {
			t.Errorf("value %q: expected true", val)
		}
	}
}

func TestAsBool_FalseVariants(t *testing.T) {
	tc := NewTypeCaster()
	for _, val := range []string{"0", "false", "no", "off"} {
		m := map[string]string{"FLAG": val}
		b, err := tc.AsBool(m, "FLAG")
		if err != nil {
			t.Errorf("value %q: unexpected error: %v", val, err)
		}
		if b {
			t.Errorf("value %q: expected false", val)
		}
	}
}

func TestAsBool_Invalid(t *testing.T) {
	tc := NewTypeCaster()
	m := map[string]string{"FLAG": "maybe"}
	_, err := tc.AsBool(m, "FLAG")
	if err == nil {
		t.Fatal("expected error for unrecognised boolean")
	}
}

func TestAsFloat_Valid(t *testing.T) {
	tc := NewTypeCaster()
	m := map[string]string{"RATIO": "3.14"}
	f, err := tc.AsFloat(m, "RATIO")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f < 3.13 || f > 3.15 {
		t.Errorf("expected ~3.14, got %f", f)
	}
}

func TestAsFloat_Invalid(t *testing.T) {
	tc := NewTypeCaster()
	m := map[string]string{"RATIO": "not-a-number"}
	_, err := tc.AsFloat(m, "RATIO")
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestAsInt_TrimsWhitespace(t *testing.T) {
	tc := NewTypeCaster()
	m := map[string]string{"PORT": "  9090  "}
	n, err := tc.AsInt(m, "PORT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 9090 {
		t.Errorf("expected 9090, got %d", n)
	}
}
