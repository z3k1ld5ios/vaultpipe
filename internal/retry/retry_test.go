package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/retry"
)

var errTransient = errors.New("transient error")

func TestDo_SucceedsFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), retry.DefaultConfig(), func(attempt int) error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesOnFailure(t *testing.T) {
	calls := 0
	cfg := retry.Config{MaxAttempts: 3, Delay: time.Millisecond, Multiplier: 1.0}
	err := retry.Do(context.Background(), cfg, func(attempt int) error {
		calls++
		if calls < 3 {
			return errTransient
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	cfg := retry.Config{MaxAttempts: 2, Delay: time.Millisecond, Multiplier: 1.0}
	err := retry.Do(context.Background(), cfg, func(attempt int) error {
		return errTransient
	})
	if !errors.Is(err, retry.ErrMaxAttemptsReached) {
		t.Fatalf("expected ErrMaxAttemptsReached, got %v", err)
	}
}

func TestDo_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := retry.Do(ctx, retry.DefaultConfig(), func(attempt int) error {
		return errTransient
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestDo_ZeroMaxAttempts_TreatedAsOne(t *testing.T) {
	calls := 0
	cfg := retry.Config{MaxAttempts: 0, Delay: time.Millisecond, Multiplier: 1.0}
	err := retry.Do(context.Background(), cfg, func(attempt int) error {
		calls++
		return errTransient
	})
	if !errors.Is(err, retry.ErrMaxAttemptsReached) {
		t.Fatalf("expected ErrMaxAttemptsReached, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected exactly 1 call, got %d", calls)
	}
}
