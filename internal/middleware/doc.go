// Package middleware provides cross-cutting controls applied around
// Vault secret operations in vaultpipe.
//
// Currently includes:
//   - RateLimiter: caps the number of secret reads within a rolling time window
//     to avoid overwhelming Vault or triggering policy-based denials.
//
// Middleware components are stateful and safe for concurrent use.
package middleware
