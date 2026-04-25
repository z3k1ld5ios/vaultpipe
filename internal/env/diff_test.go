package env

import (
	"testing"
)

func TestDiff_NoDifference(t *testing.T) {
	m := map[string]string{"A": "1", "B": "2"}
	changes := Diff(m, m)
	if len(changes) != 0 {
		t.Fatalf("expected no changes, got %d", len(changes))
	}
}

func TestDiff_DetectsAdded(t *testing.T) {
	prev := map[string]string{"A": "1"}
	next := map[string]string{"A": "1", "B": "2"}
	changes := Diff(prev, next)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != ChangeAdded || changes[0].Key != "B" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestDiff_DetectsRemoved(t *testing.T) {
	prev := map[string]string{"A": "1", "B": "2"}
	next := map[string]string{"A": "1"}
	changes := Diff(prev, next)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != ChangeRemoved || changes[0].Key != "B" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestDiff_DetectsUpdated(t *testing.T) {
	prev := map[string]string{"A": "old"}
	next := map[string]string{"A": "new"}
	changes := Diff(prev, next)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	c := changes[0]
	if c.Type != ChangeUpdated || c.OldVal != "old" || c.NewVal != "new" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestDiff_MixedChanges_SortedByKey(t *testing.T) {
	prev := map[string]string{"B": "1", "C": "old"}
	next := map[string]string{"A": "new", "C": "updated"}
	changes := Diff(prev, next)
	// A added, B removed, C updated
	if len(changes) != 3 {
		t.Fatalf("expected 3 changes, got %d", len(changes))
	}
	if changes[0].Key != "A" || changes[1].Key != "B" || changes[2].Key != "C" {
		t.Errorf("changes not sorted: %+v", changes)
	}
}

func TestHasChanges_True(t *testing.T) {
	if !HasChanges(map[string]string{"X": "1"}, map[string]string{"X": "2"}) {
		t.Error("expected HasChanges to return true")
	}
}

func TestHasChanges_False(t *testing.T) {
	m := map[string]string{"X": "1"}
	if HasChanges(m, m) {
		t.Error("expected HasChanges to return false")
	}
}

func TestFilterByType_ReturnsMatchingOnly(t *testing.T) {
	changes := []Change{
		{Key: "A", Type: ChangeAdded},
		{Key: "B", Type: ChangeRemoved},
		{Key: "C", Type: ChangeAdded},
	}
	added := FilterByType(changes, ChangeAdded)
	if len(added) != 2 {
		t.Fatalf("expected 2 added changes, got %d", len(added))
	}
}

func TestDiff_EmptyBoth(t *testing.T) {
	changes := Diff(map[string]string{}, map[string]string{})
	if len(changes) != 0 {
		t.Errorf("expected empty diff, got %d changes", len(changes))
	}
}
