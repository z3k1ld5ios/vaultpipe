package snapshot_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/snapshot"
)

func base() map[string]string {
	return map[string]string{"DB_PASS": "secret1", "API_KEY": "abc"}
}

func TestTake_ChecksumDeterministic(t *testing.T) {
	a := snapshot.Take(base())
	b := snapshot.Take(base())
	if a.Checksum != b.Checksum {
		t.Fatalf("expected same checksum, got %q vs %q", a.Checksum, b.Checksum)
	}
}

func TestEqual_SameSecrets(t *testing.T) {
	a := snapshot.Take(base())
	b := snapshot.Take(base())
	if !a.Equal(b) {
		t.Fatal("expected snapshots to be equal")
	}
}

func TestEqual_DifferentSecrets(t *testing.T) {
	a := snapshot.Take(base())
	mod := base()
	mod["DB_PASS"] = "rotated"
	b := snapshot.Take(mod)
	if a.Equal(b) {
		t.Fatal("expected snapshots to differ")
	}
}

func TestEqual_NilHandling(t *testing.T) {
	a := snapshot.Take(base())
	if a.Equal(nil) {
		t.Fatal("expected non-nil to differ from nil")
	}
}

func TestDiff_DetectsAdded(t *testing.T) {
	old := snapshot.Take(base())
	newSecrets := base()
	newSecrets["NEW_KEY"] = "val"
	next := snapshot.Take(newSecrets)
	delta := old.Diff(next)
	if _, ok := delta["NEW_KEY"]; !ok {
		t.Fatal("expected NEW_KEY in diff")
	}
}

func TestDiff_DetectsRemoved(t *testing.T) {
	old := snapshot.Take(base())
	reduced := map[string]string{"DB_PASS": "secret1"}
	next := snapshot.Take(reduced)
	delta := old.Diff(next)
	if delta["API_KEY"] != "removed" {
		t.Fatalf("expected API_KEY removed, got %q", delta["API_KEY"])
	}
}

func TestDiff_DetectsChanged(t *testing.T) {
	old := snapshot.Take(base())
	mod := base()
	mod["DB_PASS"] = "rotated"
	next := snapshot.Take(mod)
	delta := old.Diff(next)
	if delta["DB_PASS"] != "changed" {
		t.Fatalf("expected DB_PASS changed, got %q", delta["DB_PASS"])
	}
}

func TestDiff_NoChanges_EmptyDelta(t *testing.T) {
	a := snapshot.Take(base())
	b := snapshot.Take(base())
	if len(a.Diff(b)) != 0 {
		t.Fatal("expected empty diff for identical snapshots")
	}
}
