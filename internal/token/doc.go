// Package token provides token resolution for HashiCorp Vault authentication.
//
// Tokens can be sourced from three locations, checked in priority order:
//
//  1. A static token supplied directly at construction time.
//  2. An environment variable (defaults to VAULT_TOKEN).
//  3. A file path whose contents are read and trimmed.
//
// Use NewProvider to construct a resolver and call Resolve to obtain a Token.
package token
