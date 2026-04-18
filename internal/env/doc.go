// Package env provides utilities for building and injecting environment
// variable sets into child processes.
//
// Injector merges a base environment with resolved secrets, ensuring
// secret values take precedence and no duplicate keys are emitted.
//
// Resolver maps user-defined ENV_VAR names to secret keys, supporting
// optional default values using the "key:default" syntax.
package env
