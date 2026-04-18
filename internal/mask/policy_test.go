package mask_test

import (
	"testing"

	"github.com/your-org/vaultpipe/internal/mask"
)

func TestApply_LevelNone(t *testing.T) {
	p := mask.Policy{Level: mask.LevelNone, Placeholder: "***"}
	if got := p.Apply("mysecret"); got != "mysecret" {
		t.Fatalf("expected mysecret, got %s", got)
	}
}

func TestApply_LevelFull(t *testing.T) {
	p := mask.Policy{Level: mask.LevelFull, Placeholder: "***"}
	if got := p.Apply("mysecret"); got != "***" {
		t.Fatalf("expected ***, got %s", got)
	}
}

func TestApply_LevelPartial_Long(t *testing.T) {
	p := mask.DefaultPolicy()
	got := p.Apply("mysecret")
	if got[0] != 'm' || got[len(got)-1] != 't' {
		t.Fatalf("expected partial mask preserving first/last, got %s", got)
	}
}

func TestApply_LevelPartial_Short(t *testing.T) {
	p := mask.DefaultPolicy()
	if got := p.Apply("ab"); got != "***" {
		t.Fatalf("expected *** for short value, got %s", got)
	}
}

func TestApply_EmptyValue(t *testing.T) {
	p := mask.DefaultPolicy()
	if got := p.Apply(""); got != "" {
		t.Fatalf("expected empty string, got %s", got)
	}
}

func TestApplyMap_MasksAllValues(t *testing.T) {
	p := mask.Policy{Level: mask.LevelFull, Placeholder: "[redacted]"}
	input := map[string]string{"KEY1": "value1", "KEY2": "value2"}
	out := p.ApplyMap(input)
	for k, v := range out {
		if v != "[redacted]" {
			t.Errorf("key %s: expected [redacted], got %s", k, v)
		}
	}
}

func TestDefaultPolicy_Values(t *testing.T) {
	p := mask.DefaultPolicy()
	if p.Level != mask.LevelPartial {
		t.Errorf("expected LevelPartial, got %d", p.Level)
	}
	if p.Placeholder != "***" {
		t.Errorf("expected ***, got %s", p.Placeholder)
	}
}
