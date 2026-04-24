// Package env provides utilities for constructing, transforming,
// and injecting environment variable maps derived from Vault secrets.
//
// The Scope type offers key-level access control over which environment
// variables are visible to a spawned process. It supports both exact
// key matching and prefix-based wildcard patterns (e.g. "APP_*"),
// making it suitable for multi-tenant or least-privilege deployments
// where only a subset of secrets should be exposed per command.
package env
