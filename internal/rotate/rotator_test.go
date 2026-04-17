package rotate_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/rotate"
)

func staticFetcher(secrets map[string]string) rotate.SecretFetcher {
	return func(_ context.Context, _ string) (map[string]string, error) {
		return secrets, nil
	}
}

func TestDefaultConfig_Values(t *testing.T) {
	cfg := rotate.DefaultConfig()
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", cfg.Interval)
	}
	if cfg.MaxRetries != 3 {
		t.Errorf("expected 3 retries, got %d", cfg.MaxRetries)
	}
}

func TestWatch_EmptyPath_ReturnsError(t *testing.T) {
	r := rotate.New(rotate.DefaultConfig(), staticFetcher(nil))
	err := r.Watch(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestWatch_FiresCallbackOnChange(t *testing.T) {
	var calls atomic.Int32

	call := 0
	fetcher := func(_ context.Context, _ string) (map[string]string, error) {
		call++
		if call == 1 {
			return map[string]string{"key": "v1"}, nil
		}
		return map[string]string{"key": "v2"}, nil
	}

	cfg := rotate.Config{Interval: 20 * time.Millisecond, MaxRetries: 1}
	r := rotate.New(cfg, fetcher)
	r.OnRotate(func(_ string, _ map[string]string) {
		calls.Add(1)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	_ = r.Watch(ctx, "secret/app")

	if calls.Load() == 0 {
		t.Error("expected at least one rotation callback")
	}
}

func TestWatch_NoCallbackWhenUnchanged(t *testing.T) {
	var calls atomic.Int32

	cfg := rotate.Config{Interval: 20 * time.Millisecond, MaxRetries: 1}
	r := rotate.New(cfg, staticFetcher(map[string]string{"key": "stable"}))
	r.OnRotate(func(_ string, _ map[string]string) {
		calls.Add(1)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()

	_ = r.Watch(ctx, "secret/app")

	if calls.Load() != 0 {
		t.Errorf("expected no callbacks, got %d", calls.Load())
	}
}

func TestWatch_ContextCancel_Exits(t *testing.T) {
	cfg := rotate.Config{Interval: 10 * time.Millisecond, MaxRetries: 1}
	r := rotate.New(cfg, staticFetcher(map[string]string{"x": "1"}))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := r.Watch(ctx, "secret/app")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
