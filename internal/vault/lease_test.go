package vault

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newLeaseServer(t *testing.T, renewStatus, revokeStatus int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/sys/leases/renew":
			w.WriteHeader(renewStatus)
			if renewStatus == 200 {
				w.Write([]byte(`{"lease_id":"test-lease","lease_duration":3600,"renewable":true}`))
			}
		case "/v1/sys/leases/revoke":
			w.WriteHeader(revokeStatus)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestRenewLease_EmptyID(t *testing.T) {
	c, _ := NewClient("http://127.0.0.1:8200", "test-token")
	_, err := c.RenewLease(context.Background(), "", time.Hour)
	if err == nil {
		t.Fatal("expected error for empty lease ID")
	}
}

func TestRenewLease_Success(t *testing.T) {
	srv := newLeaseServer(t, http.StatusOK, http.StatusNoContent)
	defer srv.Close()

	c, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	info, err := c.RenewLease(context.Background(), "test-lease", time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.LeaseID != "test-lease" {
		t.Errorf("expected lease ID 'test-lease', got %q", info.LeaseID)
	}
	if !info.Renewable {
		t.Error("expected lease to be renewable")
	}
	if info.RenewedAt.IsZero() {
		t.Error("expected RenewedAt to be set")
	}
}

func TestRenewLease_NotFound(t *testing.T) {
	srv := newLeaseServer(t, http.StatusNotFound, http.StatusNoContent)
	defer srv.Close()

	c, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = c.RenewLease(context.Background(), "expired-lease", time.Hour)
	if err == nil {
		t.Fatal("expected error for not-found lease")
	}
}

func TestRevokeLease_EmptyID(t *testing.T) {
	c, _ := NewClient("http://127.0.0.1:8200", "test-token")
	err := c.RevokeLease(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty lease ID")
	}
}

func TestRevokeLease_Success(t *testing.T) {
	srv := newLeaseServer(t, http.StatusOK, http.StatusNoContent)
	defer srv.Close()

	c, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	if err := c.RevokeLease(context.Background(), "test-lease"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
