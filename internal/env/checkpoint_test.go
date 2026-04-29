package env

import (
	"testing"
)

func TestSave_And_Get_ReturnsCheckpoint(t *testing.T) {
	cp := NewCheckpointer()
	env := map[string]string{"KEY": "val"}
	cp.Save("snap1", env)

	got, ok := cp.Get("snap1")
	if !ok {
		t.Fatal("expected checkpoint to be found")
	}
	if got.Name != "snap1" {
		t.Errorf("expected name snap1, got %s", got.Name)
	}
	if got.Env["KEY"] != "val" {
		t.Errorf("expected KEY=val, got %s", got.Env["KEY"])
	}
}

func TestGet_MissingCheckpoint_ReturnsFalse(t *testing.T) {
	cp := NewCheckpointer()
	_, ok := cp.Get("missing")
	if ok {
		t.Error("expected false for missing checkpoint")
	}
}

func TestSave_IsolatesOriginalMap(t *testing.T) {
	cp := NewCheckpointer()
	env := map[string]string{"A": "1"}
	cp.Save("snap", env)
	env["A"] = "mutated"

	got, _ := cp.Get("snap")
	if got.Env["A"] != "1" {
		t.Errorf("checkpoint should not reflect mutation, got %s", got.Env["A"])
	}
}

func TestNames_ReturnsAllInOrder(t *testing.T) {
	cp := NewCheckpointer()
	cp.Save("first", map[string]string{})
	cp.Save("second", map[string]string{})
	cp.Save("third", map[string]string{})

	names := cp.Names()
	if len(names) != 3 {
		t.Fatalf("expected 3 names, got %d", len(names))
	}
	if names[0] != "first" || names[1] != "second" || names[2] != "third" {
		t.Errorf("unexpected order: %v", names)
	}
}

func TestBetween_DetectsChanges(t *testing.T) {
	cp := NewCheckpointer()
	cp.Save("before", map[string]string{"X": "old", "Y": "same"})
	cp.Save("after", map[string]string{"X": "new", "Y": "same", "Z": "added"})

	changes, err := cp.Between("before", "after")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !HasChanges(changes) {
		t.Error("expected changes between checkpoints")
	}
	updated := FilterByType(changes, ChangeUpdated)
	if len(updated) != 1 || updated[0].Key != "X" {
		t.Errorf("expected X to be updated, got %v", updated)
	}
	added := FilterByType(changes, ChangeAdded)
	if len(added) != 1 || added[0].Key != "Z" {
		t.Errorf("expected Z to be added, got %v", added)
	}
}

func TestBetween_MissingFrom_ReturnsError(t *testing.T) {
	cp := NewCheckpointer()
	cp.Save("after", map[string]string{})
	_, err := cp.Between("ghost", "after")
	if err == nil {
		t.Error("expected error for missing 'from' checkpoint")
	}
}

func TestBetween_MissingTo_ReturnsError(t *testing.T) {
	cp := NewCheckpointer()
	cp.Save("before", map[string]string{})
	_, err := cp.Between("before", "ghost")
	if err == nil {
		t.Error("expected error for missing 'to' checkpoint")
	}
}

func TestClear_RemovesAllCheckpoints(t *testing.T) {
	cp := NewCheckpointer()
	cp.Save("a", map[string]string{"K": "v"})
	cp.Clear()

	if names := cp.Names(); len(names) != 0 {
		t.Errorf("expected no checkpoints after clear, got %v", names)
	}
}

func TestGet_ReturnsLatestWhenDuplicateName(t *testing.T) {
	cp := NewCheckpointer()
	cp.Save("snap", map[string]string{"V": "first"})
	cp.Save("snap", map[string]string{"V": "second"})

	got, ok := cp.Get("snap")
	if !ok {
		t.Fatal("expected checkpoint")
	}
	if got.Env["V"] != "second" {
		t.Errorf("expected latest value 'second', got %s", got.Env["V"])
	}
}
