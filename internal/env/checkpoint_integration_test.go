package env

import (
	"testing"
)

// TestCheckpoint_WithOverride_TracksPipelineStages verifies that a Checkpointer
// can capture env state before and after an Override step and correctly diff them.
func TestCheckpoint_WithOverride_TracksPipelineStages(t *testing.T) {
	base := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	overrides := map[string]string{
		"DB_HOST": "prod-db.internal",
		"API_KEY": "secret-token",
	}

	cp := NewCheckpointer()
	cp.Save("initial", base)

	ov := NewOverride(overrides)
	result, err := ov.Apply(base)
	if err != nil {
		t.Fatalf("override failed: %v", err)
	}
	cp.Save("post-override", result)

	changes, err := cp.Between("initial", "post-override")
	if err != nil {
		t.Fatalf("Between failed: %v", err)
	}

	updated := FilterByType(changes, ChangeUpdated)
	if len(updated) != 1 || updated[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST updated, got %v", updated)
	}
	added := FilterByType(changes, ChangeAdded)
	if len(added) != 1 || added[0].Key != "API_KEY" {
		t.Errorf("expected API_KEY added, got %v", added)
	}
	removed := FilterByType(changes, ChangeRemoved)
	if len(removed) != 0 {
		t.Errorf("expected no removals, got %v", removed)
	}
}

// TestCheckpoint_WithDefaults_NoChangesWhenAllPresent verifies that applying
// defaults to a fully-populated env produces no diff.
func TestCheckpoint_WithDefaults_NoChangesWhenAllPresent(t *testing.T) {
	env := map[string]string{
		"LOG_LEVEL": "info",
		"TIMEOUT":   "30s",
	}
	defaults := map[string]string{
		"LOG_LEVEL": "debug",
		"TIMEOUT":   "10s",
	}

	cp := NewCheckpointer()
	cp.Save("before", env)

	da := NewDefaultsApplier(defaults)
	result := da.Apply(env)
	cp.Save("after", result)

	changes, err := cp.Between("before", "after")
	if err != nil {
		t.Fatalf("Between failed: %v", err)
	}
	if HasChanges(changes) {
		t.Errorf("expected no changes when all keys already present, got %v", changes)
	}
}
