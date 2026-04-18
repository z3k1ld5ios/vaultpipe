// Package preflight provides composable runtime checks that run before
// vaultpipe injects secrets into a process environment. Checks validate
// preconditions such as required environment variables, reachable Vault
// addresses, and non-empty secret payloads.
package preflight
