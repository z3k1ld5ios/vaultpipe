package vault

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDefaultRenewConfig_Values(t *testing.T) {
	cfg := DefaultRenewConfig()
	if cfg.RenewBeforeExpiry != 10*time.Second {
		t.Errorf("expected RenewBeforeExpiry=10s, got %v", cfg.RenewBeforeExpiry)
	}
	if cfg.MaxRenewals != 0 {
		t.Errorf("expected MaxRenewals=0, got %d", cfg.MaxRenewals)
	}
}

func TestRenewLoop_EmptyLeaseID(t *testing.T) {
	err := RenewLoop(context.Background(), "", 30*time.Second, DefaultRenewConfig(), nil)
	if err == nil {
		t.Fatal("expected error for empty leaseID")
	}
}

func TestRenewLoop_ZeroTTL(t *testing.T) {
	err := RenewLoop(context.Background(), "lease-1", 0, DefaultRenewConfig(), nil)
	if err == nil {
		t.Fatal("expected error for zero TTL")
	}
}

func TestRenewLoop_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	cfg := DefaultRenewConfig()
	cfg.RenewBeforeExpiry = 0

	err := RenewLoop(ctx, "lease-1", 1*time.Millisecond, cfg, func(_ context.Context, _ string) (time.Duration, error) {
		return 30 * time.Second, nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestRenewLoop_MaxRenewals(t *testing.T) {
	cfg := RenewConfig{
		RenewBeforeExpiry: 0,
		MaxRenewals:       2,
	}

	calls := 0
	err := RenewLoop(context.Background(), "lease-1", 1*time.Millisecond, cfg,
		func(_ context.Context, _ string) (time.Duration, error) {
			calls++
			return 1 * time.Millisecond, nil
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 2 {
		t.Errorf("expected 2 renewal calls, got %d", calls)
	}
}

func TestRenewLoop_RenewFnError(t *testing.T) {
	cfg := RenewConfig{
		RenewBeforeExpiry: 0,
		MaxRenewals:       1,
	}

	err := RenewLoop(context.Background(), "lease-1", 1*time.Millisecond, cfg,
		func(_ context.Context, _ string) (time.Duration, error) {
			return 0, errors.New("vault unavailable")
		},
	)
	if err == nil {
		t.Fatal("expected error from renewFn")
	}
}
