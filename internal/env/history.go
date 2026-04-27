package env

import (
	"fmt"
	"sync"
	"time"
)

// HistoryEntry records a snapshot of an environment map at a point in time,
// along with an optional label describing the operation that produced it.
type HistoryEntry struct {
	Timestamp time.Time
	Label     string
	Env       map[string]string
}

// History maintains an ordered, bounded log of environment snapshots.
// It is safe for concurrent use. Each entry is a deep copy of the map
// provided at the time of recording, so later mutations do not affect
// previously stored entries.
type History struct {
	mu      sync.RWMutex
	entries []HistoryEntry
	maxSize int
}

// NewHistory creates a History that retains at most maxSize entries.
// If maxSize is zero or negative, a default of 32 is used.
func NewHistory(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = 32
	}
	return &History{maxSize: maxSize}
}

// Record appends a labelled snapshot of env to the history.
// If the history is at capacity the oldest entry is evicted.
func (h *History) Record(label string, env map[string]string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}

	entry := HistoryEntry{
		Timestamp: time.Now(),
		Label:     label,
		Env:       copy,
	}

	if len(h.entries) >= h.maxSize {
		// Evict the oldest entry by shifting the slice.
		h.entries = append(h.entries[1:], entry)
	} else {
		h.entries = append(h.entries, entry)
	}
}

// Entries returns a copy of all recorded entries in chronological order.
func (h *History) Entries() []HistoryEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()

	out := make([]HistoryEntry, len(h.entries))
	copy(out, h.entries)
	return out
}

// Len returns the number of entries currently stored.
func (h *History) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.entries)
}

// Latest returns the most recently recorded entry and true, or a zero
// HistoryEntry and false if no entries have been recorded yet.
func (h *History) Latest() (HistoryEntry, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.entries) == 0 {
		return HistoryEntry{}, false
	}
	return h.entries[len(h.entries)-1], true
}

// At returns the entry at index i (0 = oldest). It returns an error if the
// index is out of range.
func (h *History) At(i int) (HistoryEntry, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if i < 0 || i >= len(h.entries) {
		return HistoryEntry{}, fmt.Errorf("history: index %d out of range [0, %d)", i, len(h.entries))
	}
	return h.entries[i], nil
}

// Clear removes all recorded entries.
func (h *History) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries = h.entries[:0]
}
