package middleware_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yourorg/vaultpipe/internal/middleware"
)

func TestChain_AllGuardsPass(t *testing.T) {
	pass := func(_ context.Context) error { return nil }
	c := middleware.NewChain(pass, pass)
	if err := c.Run(context.Background()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestChain_StopsOnFirstError(t *testing.T) {
	sentinel := errors.New("blocked")
	called := false
	fail := func(_ context.Context) error { return sentinel }
	after := func(_ context.Context) error { called = true; return nil }
	c := middleware.NewChain(fail, after)
	err := c.Run(context.Background())
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if called {
		t.Fatal("second guard should not have been called")
	}
}

func TestChain_EmptyChain(t *testing.T) {
	c := middleware.NewChain()
	if err := c.Run(context.Background()); err != nil {
		t.Fatalf("empty chain should pass, got %v", err)
	}
}

func TestChain_RateLimitGuard_Integration(t *testing.T) {
	rl := middleware.NewRateLimiter(2, time.Minute)
	c := middleware.NewChain(middleware.RateLimitGuard(rl))
	if err := c.Run(context.Background()); err != nil {
		t.Fatalf("first call should pass: %v", err)
	}
	if err := c.Run(context.Background()); err != nil {
		t.Fatalf("second call should pass: %v", err)
	}
	if err := c.Run(context.Background()); err == nil {
		t.Fatal("third call should be rate limited")
	}
}
