// Package env provides utilities for constructing, transforming, and
// validating environment variable maps used when injecting secrets into
// child processes.
//
// Chain allows callers to compose multiple env.ChainStep functions into a
// single ordered pipeline. Unlike Pipeline (which operates on typed stages),
// Chain works directly with plain map[string]string transform functions,
// making it easy to wire together ad-hoc transformations without defining
// dedicated stage types.
//
// Example:
//
//	chain := env.NewChain(
//		env.PrefixStep("APP_"),
//		env.UpperStep(),
//	)
//	result, err := chain.Apply(secrets)
package env
