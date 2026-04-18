package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/vaultpipe/internal/audit"
)

func newTestLogger() (*audit.Logger, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	return audit.NewLogger(buf), buf
}

func decodeEvent(t *testing.T, buf *bytes.Buffer) audit.Event {
	t.Helper()
	var e audit.Event
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		t.Fatalf("failed to decode audit event: %v", err)
	}
	return e
}

func TestSecretRead_LogsCorrectFields(t *testing.T) {
	logger, buf := newTestLogger()
	if err := logger.SecretRead("secret", "myapp/db"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := decodeEvent(t, buf)
	if e.Type != audit.EventSecretRead {
		t.Errorf("expected type %q, got %q", audit.EventSecretRead, e.Type)
	}
	if e.Mount != "secret" || e.Path != "myapp/db" {
		t.Errorf("unexpected mount/path: %q / %q", e.Mount, e.Path)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestSecretDenied_LogsMessage(t *testing.T) {
	logger, buf := newTestLogger()
	_ = logger.SecretDenied("secret", "myapp/db", "permission denied")
	e := decodeEvent(t, buf)
	if e.Type != audit.EventSecretDenied {
		t.Errorf("expected type %q, got %q", audit.EventSecretDenied, e.Type)
	}
	if !strings.Contains(e.Message, "permission denied") {
		t.Errorf("expected message to contain 'permission denied', got %q", e.Message)
	}
}

func TestProcessStart_LogsCommand(t *testing.T) {
	logger, buf := newTestLogger()
	_ = logger.ProcessStart("myapp --serve")
	e := decodeEvent(t, buf)
	if e.Type != audit.EventProcessStart {
		t.Errorf("expected type %q, got %q", audit.EventProcessStart, e.Type)
	}
	if e.Command != "myapp --serve" {
		t.Errorf("expected command %q, got %q", "myapp --serve", e.Command)
	}
}

func TestProcessExit_LogsExitCode(t *testing.T) {
	logger, buf := newTestLogger()
	_ = logger.ProcessExit("myapp", 1)
	e := decodeEvent(t, buf)
	if e.Type != audit.EventProcessExit {
		t.Errorf("expected type %q, got %q", audit.EventProcessExit, e.Type)
	}
	if e.ExitCode == nil || *e.ExitCode != 1 {
		t.Errorf("expected exit code 1, got %v", e.ExitCode)
	}
}

func TestLog_OutputIsNewlineTerminated(t *testing.T) {
	logger, buf := newTestLogger()
	_ = logger.SecretRead("kv", "foo/bar")
	if !strings.HasSuffix(buf.String(), "\n") {
		t.Error("expected log output to end with newline")
	}
}

func TestProcessExit_ZeroExitCode(t *testing.T) {
	logger, buf := newTestLogger()
	_ = logger.ProcessExit("myapp", 0)
	e := decodeEvent(t, buf)
	if e.ExitCode == nil || *e.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %v", e.ExitCode)
	}
}
