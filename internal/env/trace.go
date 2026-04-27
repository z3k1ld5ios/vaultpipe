package env

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// TraceLevel controls the verbosity of a trace entry.
type TraceLevel int

const (
	TraceLevelInfo  TraceLevel = iota
	TraceLevelWarn
	TraceLevelError
)

// TraceEntry records a single transformation step applied to an env map.
type TraceEntry struct {
	Step      string
	Key       string
	OldValue  string
	NewValue  string
	Level     TraceLevel
	Message   string
	Timestamp time.Time
}

// Tracer collects trace entries produced during env pipeline execution.
type Tracer struct {
	mu      sync.Mutex
	entries []TraceEntry
}

// NewTracer returns an initialised Tracer.
func NewTracer() *Tracer {
	return &Tracer{}
}

// Record appends a new trace entry.
func (t *Tracer) Record(step, key, oldVal, newVal string, level TraceLevel, msg string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = append(t.entries, TraceEntry{
		Step:      step,
		Key:       key,
		OldValue:  oldVal,
		NewValue:  newVal,
		Level:     level,
		Message:   msg,
		Timestamp: time.Now().UTC(),
	})
}

// Entries returns a copy of all recorded trace entries.
func (t *Tracer) Entries() []TraceEntry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]TraceEntry, len(t.entries))
	copy(out, t.entries)
	return out
}

// Filter returns entries whose Step matches the provided step name.
func (t *Tracer) Filter(step string) []TraceEntry {
	t.mu.Lock()
	defer t.mu.Unlock()
	var out []TraceEntry
	for _, e := range t.entries {
		if e.Step == step {
			out = append(out, e)
		}
	}
	return out
}

// Summary returns a human-readable multi-line summary of all trace entries.
func (t *Tracer) Summary() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.entries) == 0 {
		return "(no trace entries)"
	}
	var sb strings.Builder
	for _, e := range t.entries {
		sb.WriteString(fmt.Sprintf("[%s] step=%s key=%s old=%q new=%q msg=%s\n",
			e.Timestamp.Format(time.RFC3339Nano), e.Step, e.Key, e.OldValue, e.NewValue, e.Message))
	}
	return sb.String()
}

// Reset clears all recorded entries.
func (t *Tracer) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = t.entries[:0]
}
