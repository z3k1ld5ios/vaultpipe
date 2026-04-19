// Package env provides utilities for building, filtering, and injecting
// environment variables into child processes.
//
// It includes:
//   - Injector: merges secret maps into a base environment slice.
//   - Resolver: resolves individual keys from a secret map with optional defaults.
//   - Mapper: maps secret keys to env var names with optional prefix.
//   - Whitelist: restricts which keys are passed through to the process.
//   - PrefixFilter: applies or strips a namespace prefix from env var keys.
package env
