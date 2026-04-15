// Package diagnostics provides health check and connectivity
// verification utilities for vaultpipe.
package diagnostics

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Status represents the result of a health check.
type Status struct {
	Healthy bool
	Latency time.Duration
	Message string
}

// Checker performs health checks against a Vault instance.
type Checker struct {
	address    string
	httpClient *http.Client
}

// NewChecker creates a Checker for the given Vault address.
func NewChecker(address string, timeout time.Duration) *Checker {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &Checker{
		address: address,
		httpClient: &http.Client{Timeout: timeout},
	}
}

// Ping checks whether the Vault server is reachable and initialized
// by calling the /v1/sys/health endpoint.
func (c *Checker) Ping(ctx context.Context) Status {
	url := fmt.Sprintf("%s/v1/sys/health", c.address)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Status{Healthy: false, Message: fmt.Sprintf("build request: %v", err)}
	}

	start := time.Now()
	resp, err := c.httpClient.Do(req)
	latency := time.Since(start)
	if err != nil {
		return Status{Healthy: false, Latency: latency, Message: fmt.Sprintf("request failed: %v", err)}
	}
	defer resp.Body.Close()

	// Vault returns 200 (active), 429 (standby), or 473 (performance standby).
	// All indicate the server is reachable; treat anything else as unhealthy.
	switch resp.StatusCode {
	case http.StatusOK, http.StatusTooManyRequests, 473:
		return Status{Healthy: true, Latency: latency, Message: "vault reachable"}
	default:
		return Status{
			Healthy: false,
			Latency: latency,
			Message: fmt.Sprintf("unexpected status %d", resp.StatusCode),
		}
	}
}
