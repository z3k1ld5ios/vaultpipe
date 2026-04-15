// Package cache implements a lightweight in-memory TTL cache for Vault secret
// maps. It reduces redundant round-trips to the Vault API when the same secret
// path is referenced multiple times within a single vaultpipe invocation.
//
// Usage:
//
//	c := cache.New(30 * time.Second)
//	if secrets, ok := c.Get(path); ok {
//		return secrets, nil
//	}
//	secrets, err := vaultClient.ReadSecret(path)
//	if err == nil {
//		c.Set(path, secrets)
//	}
package cache
