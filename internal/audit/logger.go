// Package audit provides structured audit logging for secret access events.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// EventType represents the kind of audit event.
type EventType string

const (
	EventSecretRead   EventType = "secret_read"
	EventSecretDenied EventType = "secret_denied"
	EventProcessStart EventType = "process_start"
	EventProcessExit  EventType = "process_exit"
)

// Event represents a single audit log entry.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Type      EventType `json:"type"`
	Mount     string    `json:"mount,omitempty"`
	Path      string    `json:"path,omitempty"`
	Command   string    `json:"command,omitempty"`
	ExitCode  *int      `json:"exit_code,omitempty"`
	Message   string    `json:"message,omitempty"`
}

// Logger writes structured audit events to an output writer.
type Logger struct {
	out io.Writer
}

// NewLogger creates a Logger writing to the given writer.
// Pass nil to use stderr.
func NewLogger(out io.Writer) *Logger {
	if out == nil {
		out = os.Stderr
	}
	return &Logger{out: out}
}

// Log writes a single audit event as a JSON line.
func (l *Logger) Log(e Event) error {
	e.Timestamp = time.Now().UTC()
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal event: %w", err)
	}
	_, err = fmt.Fprintf(l.out, "%s\n", b)
	return err
}

// SecretRead logs a successful secret read event.
func (l *Logger) SecretRead(mount, path string) error {
	return l.Log(Event{Type: EventSecretRead, Mount: mount, Path: path})
}

// SecretDenied logs a denied or failed secret access event.
func (l *Logger) SecretDenied(mount, path, reason string) error {
	return l.Log(Event{Type: EventSecretDenied, Mount: mount, Path: path, Message: reason})
}

// ProcessStart logs a process launch event.
func (l *Logger) ProcessStart(command string) error {
	return l.Log(Event{Type: EventProcessStart, Command: command})
}

// ProcessExit logs a process exit event with its exit code.
func (l *Logger) ProcessExit(command string, code int) error {
	return l.Log(Event{Type: EventProcessExit, Command: command, ExitCode: &code})
}
