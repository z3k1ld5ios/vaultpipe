package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newSecretServer(t *testing.T, mount, path string, payload map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/v1/" + mount + "/data/" + path
		if r.URL.Path != expected {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
}

func TestReadSecretKV2_FlatMap(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"data": map[string]interface{}{
				"API_KEY": "abc123",
				"DB_PASS": "secret",
			},
		},
	}
	srv := newSecretServer(t, "secret", "myapp", payload)
	defer srv.Close()

	c, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	sm, err := c.ReadSecretKV2(context.Background(), "secret", "myapp")
	if err != nil {
		t.Fatalf("ReadSecretKV2: %v", err)
	}
	if sm["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", sm["API_KEY"])
	}
	if sm["DB_PASS"] != "secret" {
		t.Errorf("expected DB_PASS=secret, got %q", sm["DB_PASS"])
	}
}

func TestReadSecretKV2_NotFound_ReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer srv.Close()

	c, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = c.ReadSecretKV2(context.Background(), "secret", "missing")
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}

func TestReadMultiple_MergesAndOverrides(t *testing.T) {
	payloads := map[string]map[string]interface{}{
		"/v1/secret/data/base": {
			"data": map[string]interface{}{
				"data": map[string]interface{}{"KEY_A": "base_a", "KEY_B": "base_b"},
			},
		},
		"/v1/secret/data/override": {
			"data": map[string]interface{}{
				"data": map[string]interface{}{"KEY_B": "override_b", "KEY_C": "c"},
			},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p, ok := payloads[r.URL.Path]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(p)
	}))
	defer srv.Close()

	c, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	sm, err := c.ReadMultiple(context.Background(), "secret", []string{"base", "override"})
	if err != nil {
		t.Fatalf("ReadMultiple: %v", err)
	}
	if sm["KEY_A"] != "base_a" {
		t.Errorf("KEY_A: want base_a, got %q", sm["KEY_A"])
	}
	if sm["KEY_B"] != "override_b" {
		t.Errorf("KEY_B: want override_b, got %q", sm["KEY_B"])
	}
	if sm["KEY_C"] != "c" {
		t.Errorf("KEY_C: want c, got %q", sm["KEY_C"])
	}
}
