// Package env provides utilities for constructing, transforming, and injecting
// environment variable maps derived from Vault secrets.
//
// The DefaultsApplier type allows callers to define fallback values for secret
// keys that may not be present in a fetched secret map. Defaults are applied
// non-destructively: existing keys in the base map always take precedence.
package env
