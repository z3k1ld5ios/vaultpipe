package middleware_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourorg/vaultpipe/internal/middleware"
)

func TestAllow_WithinLimit(t *testing.T) {
	rl := middleware.NewRateLimiter(3, time.Minute)
	for i := 0; i < 3; i++ {
		if err := rl.Allow(context.Background()); err != nil {
			t.Fatalf("expected no error on call %d, got %v", i+1, err)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	rl := middleware.NewRateLimiter(2, time.Minute)
	_ = rl.Allow(context.Background())
	_ = rl.Allow(context.Background())
	err := rl.Allow(context.Background())
	if err == nil {
		t.Fatal("expected rate limit error, got nil")
	}
}

func TestAllow_ResetsAfterInterval(t *testing.T) {
	rl := middleware.NewRateLimiter(1, 50*time.Millisecond)
	if err := rl.Allow(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := rl.Allow(context.Background()); err == nil {
		t.Fatal("expected limit exceeded error")
	}
	time.Sleep(60 * time.Millisecond)
	if err := rl.Allow(context.Background()); err != nil {
		t.Fatalf("expected reset, got error: %v", err)
	}
}

func TestAllow_ContextCancelled(t *testing.T) {
	rl := middleware.NewRateLimiter(10, time.Minute)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := rl.Allow(ctx); err == nil {
		t.Fatal("expected context error")
	}
}

func TestRemaining_DecreasesWithCalls(t *testing.T) {
	rl := middleware.NewRateLimiter(5, time.Minute)
	_ = rl.Allow(context.Background())
	_ = rl.Allow(context.Background())
	if got := rl.Remaining(); got != 3 {
		t.Fatalf("expected 3 remaining, got %d", got)
	}
}

func TestRemaining_NeverNegative(t *testing.T) {
	rl := middleware.NewRateLimiter(1, time.Minute)
	_ = rl.Allow(context.Background())
	_ = rl.Allow(context.Background())
	_ = rl.Allow(context.Background())
	if got := rl.Remaining(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}
