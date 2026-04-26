package env

import (
	"strings"
	"testing"
)

func TestBuildAuditSummary_NoChanges(t *testing.T) {
	m := map[string]string{"A": "1", "B": "2"}
	s := BuildAuditSummary(m, m)
	if s.HasChanges() {
		t.Fatal("expected no changes")
	}
	if s.String() != "no changes" {
		t.Fatalf("unexpected string: %s", s.String())
	}
}

func TestBuildAuditSummary_DetectsAdded(t *testing.T) {
	before := map[string]string{"A": "1"}
	after := map[string]string{"A": "1", "B": "2"}
	s := BuildAuditSummary(before, after)
	if s.Added != 1 || s.Removed != 0 || s.Updated != 0 {
		t.Fatalf("unexpected counts: %+v", s)
	}
	if s.Changes[0].Kind != ChangeKindAdded || s.Changes[0].Key != "B" {
		t.Fatalf("unexpected change: %+v", s.Changes[0])
	}
}

func TestBuildAuditSummary_DetectsRemoved(t *testing.T) {
	before := map[string]string{"A": "1", "B": "2"}
	after := map[string]string{"A": "1"}
	s := BuildAuditSummary(before, after)
	if s.Removed != 1 {
		t.Fatalf("expected 1 removed, got %d", s.Removed)
	}
	if s.Changes[0].OldValue != "2" {
		t.Fatalf("expected old value preserved")
	}
}

func TestBuildAuditSummary_DetectsUpdated(t *testing.T) {
	before := map[string]string{"A": "old"}
	after := map[string]string{"A": "new"}
	s := BuildAuditSummary(before, after)
	if s.Updated != 1 {
		t.Fatalf("expected 1 updated, got %d", s.Updated)
	}
	c := s.Changes[0]
	if c.OldValue != "old" || c.NewValue != "new" {
		t.Fatalf("unexpected values: %+v", c)
	}
}

func TestBuildAuditSummary_MixedChanges_SortedByKey(t *testing.T) {
	before := map[string]string{"C": "3", "A": "1"}
	after := map[string]string{"A": "updated", "B": "new"}
	s := BuildAuditSummary(before, after)
	if s.Added != 1 || s.Removed != 1 || s.Updated != 1 {
		t.Fatalf("unexpected counts: %+v", s)
	}
	keys := make([]string, len(s.Changes))
	for i, c := range s.Changes {
		keys[i] = c.Key
	}
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Fatalf("changes not sorted: %v", keys)
		}
	}
}

func TestAuditSummary_Lines(t *testing.T) {
	before := map[string]string{"A": "1", "C": "3"}
	after := map[string]string{"A": "2", "B": "new"}
	s := BuildAuditSummary(before, after)
	lines := s.Lines()
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %v", len(lines), lines)
	}
	for _, l := range lines {
		if !strings.HasPrefix(l, "+") && !strings.HasPrefix(l, "-") && !strings.HasPrefix(l, "~") {
			t.Fatalf("unexpected line format: %q", l)
		}
	}
}

func TestAuditSummary_String_MultipleKinds(t *testing.T) {
	before := map[string]string{"A": "1", "C": "3"}
	after := map[string]string{"A": "2", "B": "new"}
	s := BuildAuditSummary(before, after)
	got := s.String()
	if !strings.Contains(got, "added") || !strings.Contains(got, "removed") || !strings.Contains(got, "updated") {
		t.Fatalf("unexpected summary string: %s", got)
	}
}

func TestBuildAuditSummary_EmptyBoth(t *testing.T) {
	s := BuildAuditSummary(map[string]string{}, map[string]string{})
	if s.HasChanges() {
		t.Fatal("expected no changes for empty maps")
	}
}
