// Package env provides utilities for constructing, transforming, and injecting
// environment variable sets derived from Vault secrets.
//
// The Interpolator type resolves ${KEY} and ${KEY:-default} placeholder syntax
// embedded in secret values or configuration strings. It is intentionally
// decoupled from OS environment access so that callers control the fallback
// behaviour — pass useFallback=true to allow resolution via os.Getenv when a
// key is absent from the secrets map.
package env
