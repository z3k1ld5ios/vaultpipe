// Package env provides helpers for constructing and manipulating process
// environment variables from Vault secrets.
//
// The Expander type resolves shell-style variable references ($VAR and ${VAR})
// within secret values or configuration strings, drawing substitutions from a
// secrets map and optionally from the host OS environment.
//
// Typical usage:
//
//	expander := env.NewExpander(secrets, true)
//	resolved, err := expander.ExpandMap(rawEnv)
package env
