// Package diagnostics provides health check utilities for verifying
// connectivity to a HashiCorp Vault instance before vaultpipe attempts
// to read secrets.
//
// Usage:
//
//	checker := diagnostics.NewChecker("http://127.0.0.1:8200", 5*time.Second)
//	status := checker.Ping(ctx)
//	if !status.Healthy {
//		log.Fatalf("vault unreachable: %s", status.Message)
//	}
package diagnostics
