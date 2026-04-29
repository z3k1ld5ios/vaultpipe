package env

import (
	"fmt"
	"sync"
	"time"
)

// Checkpoint captures a named snapshot of an env map at a point in time.
type Checkpoint struct {
	Name      string
	CapturedAt time.Time
	Env       map[string]string
}

// Checkpointer stores named env checkpoints and supports diffing between them.
type Checkpointer struct {
	mu          sync.RWMutex
	checkpoints []Checkpoint
}

// NewCheckpointer returns an empty Checkpointer.
func NewCheckpointer() *Checkpointer {
	return &Checkpointer{}
}

// Save captures a copy of env under the given name.
func (c *Checkpointer) Save(name string, env map[string]string) {
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.checkpoints = append(c.checkpoints, Checkpoint{
		Name:       name,
		CapturedAt: time.Now(),
		Env:        copy,
	})
}

// Get returns the checkpoint with the given name, or false if not found.
func (c *Checkpointer) Get(name string) (Checkpoint, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for i := len(c.checkpoints) - 1; i >= 0; i-- {
		if c.checkpoints[i].Name == name {
			return c.checkpoints[i], true
		}
	}
	return Checkpoint{}, false
}

// Names returns all saved checkpoint names in order.
func (c *Checkpointer) Names() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	names := make([]string, len(c.checkpoints))
	for i, cp := range c.checkpoints {
		names[i] = cp.Name
	}
	return names
}

// Between diffs two named checkpoints. Returns error if either is missing.
func (c *Checkpointer) Between(from, to string) ([]Change, error) {
	a, ok := c.Get(from)
	if !ok {
		return nil, fmt.Errorf("checkpoint %q not found", from)
	}
	b, ok := c.Get(to)
	if !ok {
		return nil, fmt.Errorf("checkpoint %q not found", to)
	}
	return Diff(a.Env, b.Env), nil
}

// Clear removes all saved checkpoints.
func (c *Checkpointer) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.checkpoints = nil
}
