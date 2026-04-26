package env_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpipe/internal/env"
)

// TestAuditBridge_WithOverride_TracksChanges verifies that applying an override
// and then building an audit summary correctly identifies all mutations.
func TestAuditBridge_WithOverride_TracksChanges(t *testing.T) {
	base := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"OLD_KEY": "will-be-removed",
	}

	ovr := env.NewOverride(map[string]string{
		"DB_HOST": "prod-db.internal",
		"API_KEY": "secret-token",
	})

	after, err := ovr.Apply(base)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	// Simulate removal by not including OLD_KEY in after.
	delete(after, "OLD_KEY")

	summary := env.BuildAuditSummary(base, after)

	if summary.Added != 1 {
		t.Errorf("expected 1 added, got %d", summary.Added)
	}
	if summary.Updated != 1 {
		t.Errorf("expected 1 updated, got %d", summary.Updated)
	}
	if summary.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", summary.Removed)
	}

	lines := summary.Lines()
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "API_KEY") {
		t.Errorf("expected API_KEY in lines: %s", joined)
	}
	if !strings.Contains(joined, "DB_HOST") {
		t.Errorf("expected DB_HOST in lines: %s", joined)
	}
	if !strings.Contains(joined, "OLD_KEY") {
		t.Errorf("expected OLD_KEY in lines: %s", joined)
	}
}

// TestAuditBridge_WithDefaults_NoChangesWhenKeyPresent verifies that applying
// defaults to a map that already has all keys produces no audit changes.
func TestAuditBridge_WithDefaults_NoChangesWhenKeyPresent(t *testing.T) {
	base := map[string]string{
		"LOG_LEVEL": "info",
		"TIMEOUT":   "30s",
	}

	da := env.NewDefaultsApplier(map[string]string{
		"LOG_LEVEL": "debug",
		"TIMEOUT":   "10s",
	})

	after := da.Apply(base)
	summary := env.BuildAuditSummary(base, after)

	if summary.HasChanges() {
		t.Errorf("expected no changes, got: %s", summary.String())
	}
}
