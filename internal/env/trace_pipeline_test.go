package env

import (
	"errors"
	"testing"
)

// stubStep is a minimal PipelineStep for testing.
type stubStep struct {
	out map[string]string
	err error
}

func (s *stubStep) Apply(_ map[string]string) (map[string]string, error) {
	return s.out, s.err
}

func TestTracedStep_NoChanges_NoEntries(t *testing.T) {
	tr := NewTracer()
	inner := &stubStep{out: map[string]string{"A": "1", "B": "2"}}
	ts := NewTracedStep("noop", inner, tr)

	_, err := ts.Apply(map[string]string{"A": "1", "B": "2"})
	if err != nil {
		t.Fatal(err)
	}
	if len(tr.Entries()) != 0 {
		t.Errorf("expected no entries for unchanged map, got %d", len(tr.Entries()))
	}
}

func TestTracedStep_DetectsModified(t *testing.T) {
	tr := NewTracer()
	inner := &stubStep{out: map[string]string{"A": "upper"}}
	ts := NewTracedStep("transform", inner, tr)

	_, err := ts.Apply(map[string]string{"A": "lower"})
	if err != nil {
		t.Fatal(err)
	}
	entries := tr.Filter("transform")
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].OldValue != "lower" || entries[0].NewValue != "upper" {
		t.Errorf("unexpected values: %+v", entries[0])
	}
}

func TestTracedStep_DetectsAdded(t *testing.T) {
	tr := NewTracer()
	inner := &stubStep{out: map[string]string{"A": "1", "NEW": "x"}}
	ts := NewTracedStep("override", inner, tr)

	_, _ = ts.Apply(map[string]string{"A": "1"})

	entries := tr.Filter("override")
	if len(entries) != 1 || entries[0].Key != "NEW" || entries[0].Message != "added" {
		t.Errorf("expected added entry for NEW, got %+v", entries)
	}
}

func TestTracedStep_DetectsRemoved(t *testing.T) {
	tr := NewTracer()
	inner := &stubStep{out: map[string]string{}}
	ts := NewTracedStep("filter", inner, tr)

	_, _ = ts.Apply(map[string]string{"GONE": "val"})

	entries := tr.Filter("filter")
	if len(entries) != 1 || entries[0].Key != "GONE" || entries[0].Message != "removed" {
		t.Errorf("expected removed entry for GONE, got %+v", entries)
	}
}

func TestTracedStep_InnerError_RecordsErrorEntry(t *testing.T) {
	tr := NewTracer()
	inner := &stubStep{err: errors.New("boom")}
	ts := NewTracedStep("fail", inner, tr)

	_, err := ts.Apply(map[string]string{"K": "v"})
	if err == nil {
		t.Fatal("expected error")
	}
	entries := tr.Entries()
	if len(entries) != 1 || entries[0].Level != TraceLevelError {
		t.Errorf("expected one error entry, got %+v", entries)
	}
}

func TestTracedStep_StepNamePropagated(t *testing.T) {
	tr := NewTracer()
	inner := &stubStep{out: map[string]string{"X": "new"}}
	ts := NewTracedStep("myStep", inner, tr)

	_, _ = ts.Apply(map[string]string{"X": "old"})

	for _, e := range tr.Entries() {
		if e.Step != "myStep" {
			t.Errorf("expected step=myStep, got %s", e.Step)
		}
	}
}
