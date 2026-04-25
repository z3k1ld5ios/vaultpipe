package env

import (
	"testing"
)

func TestDiffMaps_NoDifference(t *testing.T) {
	prev := map[string]string{"A": "1", "B": "2"}
	next := map[string]string{"A": "1", "B": "2"}
	d := DiffMaps(prev, next)
	if d.HasChanges() {
		t.Errorf("expected no changes, got: %s", d.Summary())
	}
}

func TestDiffMaps_DetectsAdded(t *testing.T) {
	prev := map[string]string{"A": "1"}
	next := map[string]string{"A": "1", "B": "2"}
	d := DiffMaps(prev, next)
	if len(d.Added) != 1 || d.Added["B"] != "2" {
		t.Errorf("expected B=2 in Added, got: %v", d.Added)
	}
	if len(d.Removed) != 0 || len(d.Changed) != 0 {
		t.Errorf("unexpected removals or changes: %s", d.Summary())
	}
}

func TestDiffMaps_DetectsRemoved(t *testing.T) {
	prev := map[string]string{"A": "1", "B": "2"}
	next := map[string]string{"A": "1"}
	d := DiffMaps(prev, next)
	if len(d.Removed) != 1 || d.Removed["B"] != "2" {
		t.Errorf("expected B=2 in Removed, got: %v", d.Removed)
	}
}

func TestDiffMaps_DetectsChanged(t *testing.T) {
	prev := map[string]string{"A": "old"}
	next := map[string]string{"A": "new"}
	d := DiffMaps(prev, next)
	if len(d.Changed) != 1 {
		t.Fatalf("expected 1 changed key, got %d", len(d.Changed))
	}
	pair := d.Changed["A"]
	if pair[0] != "old" || pair[1] != "new" {
		t.Errorf("expected [old new], got %v", pair)
	}
}

func TestDiffMaps_EmptyBoth(t *testing.T) {
	d := DiffMaps(map[string]string{}, map[string]string{})
	if d.HasChanges() {
		t.Error("expected no changes for two empty maps")
	}
}

func TestDiff_ChangedKeys_Sorted(t *testing.T) {
	prev := map[string]string{"C": "1", "A": "1"}
	next := map[string]string{"B": "2", "A": "changed"}
	d := DiffMaps(prev, next)
	keys := d.ChangedKeys()
	if len(keys) != 3 {
		t.Fatalf("expected 3 changed keys, got %d: %v", len(keys), keys)
	}
	if keys[0] != "A" || keys[1] != "B" || keys[2] != "C" {
		t.Errorf("expected sorted [A B C], got %v", keys)
	}
}

func TestDiff_Summary_Format(t *testing.T) {
	prev := map[string]string{"A": "1"}
	next := map[string]string{"B": "2"}
	d := DiffMaps(prev, next)
	got := d.Summary()
	want := "added=1 removed=1 changed=0"
	if got != want {
		t.Errorf("Summary() = %q, want %q", got, want)
	}
}
