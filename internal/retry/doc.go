// Package retry implements exponential-backoff retry logic used throughout
// vaultpipe when communicating with HashiCorp Vault.
//
// Usage:
//
//	err := retry.Do(ctx, retry.DefaultConfig(), func(attempt int) error {
//		return vault.ReadSecret(path)
//	})
//
// The multiplier field on Config controls backoff growth. A multiplier of 1.0
// produces a constant delay, while values above 1.0 grow the delay
// exponentially on each failed attempt.
package retry
