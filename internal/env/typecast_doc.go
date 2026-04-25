// Package env provides utilities for transforming, resolving, and injecting
// environment variables sourced from HashiCorp Vault secrets.
//
// TypeCaster offers a lightweight mechanism to convert the string values that
// Vault always returns into typed Go primitives (int64, bool, float64) without
// requiring reflection or external dependencies. It is intentionally stateless
// so that a single instance can be shared across goroutines safely.
package env
