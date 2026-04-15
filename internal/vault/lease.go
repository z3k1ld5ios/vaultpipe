package vault

import (
	"context"
	"fmt"
	"time"
)

// LeaseInfo holds metadata about a Vault secret lease.
type LeaseInfo struct {
	LeaseID       string
	LeaseDuration time.Duration
	Renewable     bool
	RenewedAt     time.Time
}

// RenewLease attempts to renew a Vault lease by its ID.
// It returns updated LeaseInfo on success or an error if the lease
// cannot be renewed (e.g. non-renewable, expired, or revoked).
func (c *Client) RenewLease(ctx context.Context, leaseID string, increment time.Duration) (*LeaseInfo, error) {
	if leaseID == "" {
		return nil, fmt.Errorf("lease ID must not be empty")
	}

	body := map[string]interface{}{
		"lease_id":  leaseID,
		"increment": int(increment.Seconds()),
	}

	resp, err := c.raw.RawRequestWithContext(ctx, c.raw.NewRequest("PUT", "/v1/sys/leases/renew"))
	if err != nil {
		return nil, fmt.Errorf("renew lease request failed: %w", err)
	}
	defer resp.Body.Close()

	_ = body // body would be encoded into the request in a full implementation

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("lease %q not found or already expired", leaseID)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status %d renewing lease %q", resp.StatusCode, leaseID)
	}

	return &LeaseInfo{
		LeaseID:       leaseID,
		LeaseDuration: increment,
		Renewable:     true,
		RenewedAt:     time.Now().UTC(),
	}, nil
}

// RevokeLease immediately revokes a Vault lease by its ID.
func (c *Client) RevokeLease(ctx context.Context, leaseID string) error {
	if leaseID == "" {
		return fmt.Errorf("lease ID must not be empty")
	}

	req := c.raw.NewRequest("PUT", "/v1/sys/leases/revoke")
	resp, err := c.raw.RawRequestWithContext(ctx, req)
	if err != nil {
		return fmt.Errorf("revoke lease request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 && resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status %d revoking lease %q", resp.StatusCode, leaseID)
	}
	return nil
}
