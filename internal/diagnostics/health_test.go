package diagnostics_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/diagnostics"
)

func newHealthServer(statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
	}))
}

func TestPing_HealthyVault(t *testing.T) {
	srv := newHealthServer(http.StatusOK)
	defer srv.Close()

	checker := diagnostics.NewChecker(srv.URL, 3*time.Second)
	status := checker.Ping(context.Background())

	if !status.Healthy {
		t.Fatalf("expected healthy, got message: %s", status.Message)
	}
	if status.Latency <= 0 {
		t.Error("expected non-zero latency")
	}
}

func TestPing_StandbyVault(t *testing.T) {
	srv := newHealthServer(http.StatusTooManyRequests) // 429 = standby
	defer srv.Close()

	checker := diagnostics.NewChecker(srv.URL, 3*time.Second)
	status := checker.Ping(context.Background())

	if !status.Healthy {
		t.Fatalf("standby vault should be considered healthy, got: %s", status.Message)
	}
}

func TestPing_UnhealthyStatus(t *testing.T) {
	srv := newHealthServer(http.StatusServiceUnavailable)
	defer srv.Close()

	checker := diagnostics.NewChecker(srv.URL, 3*time.Second)
	status := checker.Ping(context.Background())

	if status.Healthy {
		t.Fatal("expected unhealthy for 503 response")
	}
}

func TestPing_UnreachableServer(t *testing.T) {
	checker := diagnostics.NewChecker("http://127.0.0.1:19999", 500*time.Millisecond)
	status := checker.Ping(context.Background())

	if status.Healthy {
		t.Fatal("expected unhealthy for unreachable server")
	}
}

func TestPing_ContextCancelled(t *testing.T) {
	srv := newHealthServer(http.StatusOK)
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	checker := diagnostics.NewChecker(srv.URL, 3*time.Second)
	status := checker.Ping(ctx)

	if status.Healthy {
		t.Fatal("expected unhealthy when context is cancelled")
	}
}

func TestNewChecker_DefaultTimeout(t *testing.T) {
	// zero timeout should default to 5s without panicking
	checker := diagnostics.NewChecker("http://127.0.0.1:8200", 0)
	if checker == nil {
		t.Fatal("expected non-nil checker")
	}
}
