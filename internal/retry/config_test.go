package retry_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/retry"
)

func TestDefaultConfig_Values(t *testing.T) {
	cfg := retry.DefaultConfig()

	if cfg.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", cfg.MaxAttempts)
	}
	if cfg.Delay != 500*time.Millisecond {
		t.Errorf("expected Delay=500ms, got %v", cfg.Delay)
	}
	if cfg.Multiplier != 2.0 {
		t.Errorf("expected Multiplier=2.0, got %f", cfg.Multiplier)
	}
}

func TestDo_ExponentialBackoff_IncreasesDelay(t *testing.T) {
	var timestamps []time.Time
	cfg := retry.Config{MaxAttempts: 3, Delay: 10 * time.Millisecond, Multiplier: 2.0}

	_ = retry.Do(context.Background(), cfg, func(attempt int) error {
		timestamps = append(timestamps, time.Now())
		return errTransient
	})

	if len(timestamps) != 3 {
		t.Fatalf("expected 3 timestamps, got %d", len(timestamps))
	}

	gap1 := timestamps[1].Sub(timestamps[0])
	gap2 := timestamps[2].Sub(timestamps[1])

	if gap2 < gap1 {
		t.Errorf("expected second gap (%v) >= first gap (%v) due to backoff", gap2, gap1)
	}
}
