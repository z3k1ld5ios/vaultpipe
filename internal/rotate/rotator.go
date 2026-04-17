// Package rotate provides secret rotation detection and callback triggering.
package rotate

import (
	"context"
	"fmt"
	"time"
)

// SecretFetcher retrieves a map of secrets by path.
type SecretFetcher func(ctx context.Context, path string) (map[string]string, error)

// OnRotateFunc is called when a secret change is detected.
type OnRotateFunc func(path string, updated map[string]string)

// Config holds rotation polling configuration.
type Config struct {
	Interval   time.Duration
	MaxRetries int
}

// DefaultConfig returns sensible rotation defaults.
func DefaultConfig() Config {
	return Config{
		Interval:   30 * time.Second,
		MaxRetries: 3,
	}
}

// Rotator polls a secret path and fires callbacks on change.
type Rotator struct {
	cfg     Config
	fetch   SecretFetcher
	callbacks []OnRotateFunc
}

// New creates a new Rotator.
func New(cfg Config, fetch SecretFetcher) *Rotator {
	return &Rotator{cfg: cfg, fetch: fetch}
}

// OnRotate registers a callback invoked when rotation is detected.
func (r *Rotator) OnRotate(fn OnRotateFunc) {
	r.callbacks = append(r.callbacks, fn)
}

// Watch polls the given path until ctx is cancelled.
func (r *Rotator) Watch(ctx context.Context, path string) error {
	if path == "" {
		return fmt.Errorf("rotate: path must not be empty")
	}

	current, err := r.fetch(ctx, path)
	if err != nil {
		return fmt.Errorf("rotate: initial fetch failed: %w", err)
	}

	ticker := time.NewTicker(r.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			updated, err := r.fetch(ctx, path)
			if err != nil {
				continue
			}
			if !equal(current, updated) {
				for _, cb := range r.callbacks {
					cb(path, updated)
				}
				current = updated
			}
		}
	}
}

func equal(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
