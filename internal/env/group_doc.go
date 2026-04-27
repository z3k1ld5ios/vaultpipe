// Package env provides utilities for constructing, transforming, and inspecting
// environment variable maps used throughout vaultpipe.
//
// Group partitions a flat env map into named logical groups using key prefixes.
// This is useful when secrets from multiple Vault paths are merged into a single
// map and downstream consumers need to distinguish their origin or domain.
//
// Example usage:
//
//	g := env.NewGroup(map[string]string{
//		"DB_":    "database",
//		"CACHE_": "cache",
//	})
//	result, err := g.Apply(envMap)
//	dbKeys := result.Keys("database")
package env
