package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newMockVaultServer(t *testing.T, responseBody map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
			t.Fatalf("failed to encode mock response: %v", err)
		}
	}))
}

func TestNewClient_MissingAddress(t *testing.T) {
	t.Setenv("VAULT_ADDR", "")
	_, err := NewClient(Config{})
	if err == nil {
		t.Fatal("expected error when vault address is missing, got nil")
	}
}

func TestNewClient_UsesEnvFallback(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	t.Setenv("VAULT_ADDR", server.URL)
	t.Setenv("VAULT_TOKEN", "test-token")

	client, err := NewClient(Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestReadSecretKV2_Success(t *testing.T) {
	mockResponse := map[string]interface{}{
		"data": map[string]interface{}{
			"data": map[string]interface{}{
				"DB_PASSWORD": "supersecret",
				"API_KEY":     "abc123",
			},
		},
	}
	server := newMockVaultServer(t, mockResponse)
	defer server.Close()

	client, err := NewClient(Config{Address: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	secrets, err := client.ReadSecretKV2("secret", "myapp/config")
	if err != nil {
		t.Fatalf("unexpected error reading secret: %v", err)
	}

	if secrets["DB_PASSWORD"] != "supersecret" {
		t.Errorf("expected DB_PASSWORD=supersecret, got %q", secrets["DB_PASSWORD"])
	}
	if secrets["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", secrets["API_KEY"])
	}
}

func TestReadSecretKV2_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`null`))
	}))
	defer server.Close()

	client, err := NewClient(Config{Address: server.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	_, err = client.ReadSecretKV2("secret", "nonexistent/path")
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}
