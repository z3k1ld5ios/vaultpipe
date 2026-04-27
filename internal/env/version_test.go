package env

import (
	"testing"
)

func TestNewVersioner_InitialGeneration(t *testing.T) {
	v := NewVersioner(map[string]string{"A": "1"})
	if v.Version().Generation != 0 {
		t.Fatalf("expected generation 0, got %d", v.Version().Generation)
	}
}

func TestApply_NoChange_GenerationUnchanged(t *testing.T) {
	base := map[string]string{"KEY": "value"}
	v := NewVersioner(base)
	ver1 := v.Apply(map[string]string{"KEY": "value"})
	if ver1.Generation != 0 {
		t.Fatalf("generation should not increment on identical content, got %d", ver1.Generation)
	}
}

func TestApply_Change_IncrementsGeneration(t *testing.T) {
	v := NewVersioner(map[string]string{"KEY": "old"})
	ver := v.Apply(map[string]string{"KEY": "new"})
	if ver.Generation != 1 {
		t.Fatalf("expected generation 1 after change, got %d", ver.Generation)
	}
}

func TestApply_MultipleChanges_GenerationMonotonic(t *testing.T) {
	v := NewVersioner(nil)
	v.Apply(map[string]string{"A": "1"})
	v.Apply(map[string]string{"A": "2"})
	ver := v.Apply(map[string]string{"A": "3"})
	if ver.Generation != 3 {
		t.Fatalf("expected generation 3, got %d", ver.Generation)
	}
}

func TestApply_ChecksumChangesWithContent(t *testing.T) {
	v := NewVersioner(map[string]string{"X": "a"})
	before := v.Version().Checksum
	v.Apply(map[string]string{"X": "b"})
	after := v.Version().Checksum
	if before == after {
		t.Fatal("expected checksum to change after content update")
	}
}

func TestCurrent_ReturnsCopy(t *testing.T) {
	v := NewVersioner(map[string]string{"K": "v"})
	copy := v.Current()
	copy["K"] = "mutated"
	if v.Current()["K"] != "v" {
		t.Fatal("Versioner internal state was mutated through Current()")
	}
}

func TestVersionString_Format(t *testing.T) {
	v := NewVersioner(map[string]string{"A": "1"})
	v.Apply(map[string]string{"A": "2"})
	s := v.Version().String()
	if len(s) == 0 {
		t.Fatal("Version.String() returned empty string")
	}
}

func TestEnvChecksum_Deterministic(t *testing.T) {
	m := map[string]string{"B": "2", "A": "1", "C": "3"}
	if envChecksum(m) != envChecksum(m) {
		t.Fatal("checksum is not deterministic")
	}
}

func TestEnvChecksum_EmptyMap(t *testing.T) {
	sum := envChecksum(map[string]string{})
	if sum == "" {
		t.Fatal("expected non-empty checksum for empty map")
	}
}
