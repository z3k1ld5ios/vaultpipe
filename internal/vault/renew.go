package vault

import (
	"context"
	"fmt"
	"time"
)

// RenewConfig holds configuration for automatic lease renewal.
type RenewConfig struct {
	// RenewBeforeExpiry is how long before expiry to attempt renewal.
	RenewBeforeExpiry time.Duration
	// MaxRenewals is the maximum number of renewal cycles (0 = unlimited).
	MaxRenewals int
}

// DefaultRenewConfig returns a sensible default renewal configuration.
func DefaultRenewConfig() RenewConfig {
	return RenewConfig{
		RenewBeforeExpiry: 10 * time.Second,
		MaxRenewals:       0,
	}
}

// RenewLoop runs a background loop that renews the given lease before it
// expires. It blocks until ctx is cancelled or MaxRenewals is reached.
// renewFn is called each cycle and should return the new TTL on success.
func RenewLoop(
	ctx context.Context,
	leaseID string,
	initialTTL time.Duration,
	cfg RenewConfig,
	renewFn func(ctx context.Context, leaseID string) (time.Duration, error),
) error {
	if leaseID == "" {
		return fmt.Errorf("renew loop: leaseID must not be empty")
	}
	if initialTTL <= 0 {
		return fmt.Errorf("renew loop: initialTTL must be positive")
	}

	ttl := initialTTL
	renewals := 0

	for {
		waitFor := ttl - cfg.RenewBeforeExpiry
		if waitFor <= 0 {
			waitFor = ttl / 2
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitFor):
		}

		newTTL, err := renewFn(ctx, leaseID)
		if err != nil {
			return fmt.Errorf("renew loop: renewal failed for lease %s: %w", leaseID, err)
		}
		ttl = newTTL
		renewals++

		if cfg.MaxRenewals > 0 && renewals >= cfg.MaxRenewals {
			return nil
		}
	}
}
