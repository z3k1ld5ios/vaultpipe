package env

import (
	"fmt"
	"strings"
)

// ChangeKind classifies the type of environment variable change.
type ChangeKind string

const (
	ChangeKindAdded   ChangeKind = "added"
	ChangeKindRemoved ChangeKind = "removed"
	ChangeKindUpdated ChangeKind = "updated"
)

// EnvChange describes a single key-level change between two env snapshots.
type EnvChange struct {
	Key      string
	Kind     ChangeKind
	OldValue string
	NewValue string
}

// AuditSummary holds a human-readable and structured summary of env changes.
type AuditSummary struct {
	Changes []EnvChange
	Added   int
	Removed int
	Updated int
}

// HasChanges returns true when at least one change is recorded.
func (s AuditSummary) HasChanges() bool {
	return len(s.Changes) > 0
}

// Lines returns each change as a formatted string, suitable for logging.
func (s AuditSummary) Lines() []string {
	out := make([]string, 0, len(s.Changes))
	for _, c := range s.Changes {
		switch c.Kind {
		case ChangeKindAdded:
			out = append(out, fmt.Sprintf("+ %s", c.Key))
		case ChangeKindRemoved:
			out = append(out, fmt.Sprintf("- %s", c.Key))
		case ChangeKindUpdated:
			out = append(out, fmt.Sprintf("~ %s", c.Key))
		}
	}
	return out
}

// String returns a compact one-line summary.
func (s AuditSummary) String() string {
	parts := []string{}
	if s.Added > 0 {
		parts = append(parts, fmt.Sprintf("%d added", s.Added))
	}
	if s.Removed > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", s.Removed))
	}
	if s.Updated > 0 {
		parts = append(parts, fmt.Sprintf("%d updated", s.Updated))
	}
	if len(parts) == 0 {
		return "no changes"
	}
	return strings.Join(parts, ", ")
}

// BuildAuditSummary compares two env maps and produces an AuditSummary.
func BuildAuditSummary(before, after map[string]string) AuditSummary {
	var changes []EnvChange

	for k, newVal := range after {
		if oldVal, ok := before[k]; !ok {
			changes = append(changes, EnvChange{Key: k, Kind: ChangeKindAdded, NewValue: newVal})
		} else if oldVal != newVal {
			changes = append(changes, EnvChange{Key: k, Kind: ChangeKindUpdated, OldValue: oldVal, NewValue: newVal})
		}
	}

	for k, oldVal := range before {
		if _, ok := after[k]; !ok {
			changes = append(changes, EnvChange{Key: k, Kind: ChangeKindRemoved, OldValue: oldVal})
		}
	}

	sortEnvChanges(changes)

	s := AuditSummary{Changes: changes}
	for _, c := range changes {
		switch c.Kind {
		case ChangeKindAdded:
			s.Added++
		case ChangeKindRemoved:
			s.Removed++
		case ChangeKindUpdated:
			s.Updated++
		}
	}
	return s
}

func sortEnvChanges(changes []EnvChange) {
	for i := 1; i < len(changes); i++ {
		for j := i; j > 0 && changes[j].Key < changes[j-1].Key; j-- {
			changes[j], changes[j-1] = changes[j-1], changes[j]
		}
	}
}
