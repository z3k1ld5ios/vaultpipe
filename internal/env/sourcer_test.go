package env

import (
	"testing"
)

func TestMerge_PriorityOrderWins(t *testing.T) {
	low := Source{Name: "defaults", Priority: PriorityLow, Values: map[string]string{"KEY": "low", "ONLY_LOW": "yes"}}
	high := Source{Name: "vault", Priority: PriorityHigh, Values: map[string]string{"KEY": "high"}}

	s := NewSourcer(low, high)
	got, err := s.Merge()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "high" {
		t.Errorf("expected KEY=high, got %q", got["KEY"])
	}
	if got["ONLY_LOW"] != "yes" {
		t.Errorf("expected ONLY_LOW=yes, got %q", got["ONLY_LOW"])
	}
}

func TestMerge_LowPriorityPreservedWhenNoConflict(t *testing.T) {
	low := Source{Name: "defaults", Priority: PriorityLow, Values: map[string]string{"A": "1"}}
	normal := Source{Name: "vault", Priority: PriorityNormal, Values: map[string]string{"B": "2"}}

	s := NewSourcer(low, normal)
	got, err := s.Merge()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["A"] != "1" || got["B"] != "2" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestMerge_NilValuesReturnsError(t *testing.T) {
	bad := Source{Name: "broken", Priority: PriorityNormal, Values: nil}
	s := NewSourcer(bad)
	_, err := s.Merge()
	if err == nil {
		t.Fatal("expected error for nil values map, got nil")
	}
}

func TestMerge_EmptySources_ReturnsEmptyMap(t *testing.T) {
	s := NewSourcer()
	got, err := s.Merge()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestNames_ReturnsRegistrationOrder(t *testing.T) {
	s := NewSourcer(
		Source{Name: "alpha", Priority: PriorityLow, Values: map[string]string{}},
		Source{Name: "beta", Priority: PriorityHigh, Values: map[string]string{}},
	)
	names := s.Names()
	if len(names) != 2 || names[0] != "alpha" || names[1] != "beta" {
		t.Errorf("unexpected names: %v", names)
	}
}

func TestMerge_ThreePriorities_HighestWins(t *testing.T) {
	low := Source{Name: "low", Priority: PriorityLow, Values: map[string]string{"X": "low"}}
	norm := Source{Name: "norm", Priority: PriorityNormal, Values: map[string]string{"X": "normal"}}
	high := Source{Name: "high", Priority: PriorityHigh, Values: map[string]string{"X": "high"}}

	s := NewSourcer(high, low, norm) // intentionally unordered
	got, err := s.Merge()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["X"] != "high" {
		t.Errorf("expected X=high, got %q", got["X"])
	}
}
