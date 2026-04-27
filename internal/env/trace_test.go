package env

import (
	"strings"
	"sync"
	"testing"
)

func TestNewTracer_StartsEmpty(t *testing.T) {
	tr := NewTracer()
	if len(tr.Entries()) != 0 {
		t.Fatal("expected no entries on a new tracer")
	}
}

func TestRecord_AppendsEntry(t *testing.T) {
	tr := NewTracer()
	tr.Record("coerce", "PORT", "8080", "8080", TraceLevelInfo, "no change")
	entries := tr.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Step != "coerce" || e.Key != "PORT" {
		t.Errorf("unexpected entry fields: %+v", e)
	}
}

func TestFilter_ReturnsByStep(t *testing.T) {
	tr := NewTracer()
	tr.Record("coerce", "A", "", "a", TraceLevelInfo, "")
	tr.Record("override", "B", "", "b", TraceLevelInfo, "")
	tr.Record("coerce", "C", "", "c", TraceLevelInfo, "")

	got := tr.Filter("coerce")
	if len(got) != 2 {
		t.Fatalf("expected 2 coerce entries, got %d", len(got))
	}
	for _, e := range got {
		if e.Step != "coerce" {
			t.Errorf("unexpected step: %s", e.Step)
		}
	}
}

func TestFilter_NoMatch_ReturnsNil(t *testing.T) {
	tr := NewTracer()
	tr.Record("coerce", "A", "", "a", TraceLevelInfo, "")
	if tr.Filter("missing") != nil {
		t.Error("expected nil for unmatched step")
	}
}

func TestSummary_ContainsKeyAndStep(t *testing.T) {
	tr := NewTracer()
	tr.Record("transform", "FOO", "bar", "BAR", TraceLevelInfo, "uppercased")
	s := tr.Summary()
	if !strings.Contains(s, "transform") || !strings.Contains(s, "FOO") {
		t.Errorf("summary missing expected content: %s", s)
	}
}

func TestSummary_EmptyTracer_ReturnsPlaceholder(t *testing.T) {
	tr := NewTracer()
	if tr.Summary() != "(no trace entries)" {
		t.Error("expected placeholder for empty tracer")
	}
}

func TestReset_ClearsEntries(t *testing.T) {
	tr := NewTracer()
	tr.Record("step", "K", "", "v", TraceLevelInfo, "")
	tr.Reset()
	if len(tr.Entries()) != 0 {
		t.Error("expected entries to be cleared after Reset")
	}
}

func TestEntries_ReturnsCopy(t *testing.T) {
	tr := NewTracer()
	tr.Record("step", "K", "old", "new", TraceLevelWarn, "msg")
	copy1 := tr.Entries()
	copy1[0].Key = "mutated"
	copy2 := tr.Entries()
	if copy2[0].Key == "mutated" {
		t.Error("Entries should return an isolated copy")
	}
}

func TestRecord_ConcurrentSafe(t *testing.T) {
	tr := NewTracer()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			tr.Record("step", "K", "", "v", TraceLevelInfo, "")
		}(i)
	}
	wg.Wait()
	if len(tr.Entries()) != 50 {
		t.Errorf("expected 50 entries, got %d", len(tr.Entries()))
	}
}
