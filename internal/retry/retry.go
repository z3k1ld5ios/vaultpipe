// Package retry provides configurable retry logic for transient failures
// such as Vault connectivity issues or temporary secret unavailability.
package retry

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrMaxAttemptsReached is returned when all retry attempts are exhausted.
var ErrMaxAttemptsReached = errors.New("max retry attempts reached")

// Config holds retry behaviour settings.
type Config struct {
	MaxAttempts int
	Delay       time.Duration
	Multiplier  float64 // backoff multiplier; 1.0 = constant delay
}

// DefaultConfig returns a sensible default retry configuration.
func DefaultConfig() Config {
	return Config{
		MaxAttempts: 3,
		Delay:       500 * time.Millisecond,
		Multiplier:  2.0,
	}
}

// Do executes fn up to cfg.MaxAttempts times, backing off between attempts.
// If ctx is cancelled the function returns immediately with ctx.Err().
// fn receives the current attempt number (1-based).
func Do(ctx context.Context, cfg Config, fn func(attempt int) error) error {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 1
	}

	delay := cfg.Delay

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := fn(attempt); err == nil {
			return nil
		} else if attempt == cfg.MaxAttempts {
			return fmt.Errorf("%w: last error: %w", ErrMaxAttemptsReached, err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}

		if cfg.Multiplier > 0 {
			delay = time.Duration(float64(delay) * cfg.Multiplier)
		}
	}

	return ErrMaxAttemptsReached
}
