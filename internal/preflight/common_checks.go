package preflight

import (
	"fmt"
	"net/http"
	"time"
)

// ReachableVault returns a Check that performs an HTTP GET against the Vault
// health endpoint and fails if the server is unreachable or returns a 5xx.
func ReachableVault(addr string) Check {
	return Check{
		Name: "vault:reachable",
		Run: func() error {
			url := fmt.Sprintf("%s/v1/sys/health?standbyok=true&sealedok=false", addr)
			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Get(url) //nolint:noctx
			if err != nil {
				return fmt.Errorf("vault unreachable at %s: %w", addr, err)
			}
			defer resp.Body.Close()
			if resp.StatusCode >= 500 {
				return fmt.Errorf("vault returned unhealthy status %d", resp.StatusCode)
			}
			return nil
		},
	}
}

// TokenPresent returns a Check that fails when the provided token is empty.
func TokenPresent(token string) Check {
	return Check{
		Name: "token:present",
		Run: func() error {
			if token == "" {
				return fmt.Errorf("vault token is empty; set VAULT_TOKEN or provide a token file")
			}
			return nil
		},
	}
}

// CommandNonEmpty returns a Check that fails when no command is specified.
func CommandNonEmpty(args []string) Check {
	return Check{
		Name: "command:non-empty",
		Run: func() error {
			if len(args) == 0 {
				return fmt.Errorf("no command specified to execute")
			}
			return nil
		},
	}
}
