// Package env provides utilities for injecting, transforming, and validating
// environment variables sourced from HashiCorp Vault secrets.
//
// The Schema type allows callers to declare the expected shape of an environment
// map — including required keys, expected types (string, int, bool), and optional
// regex patterns — and validate a resolved env map before it is injected into a
// child process.
//
// Example usage:
//
//	schema := env.NewSchema([]env.FieldSchema{
//		{Key: "PORT", Type: env.TypeInt, Required: true},
//		{Key: "LOG_LEVEL", Type: env.TypeString, Pattern: `^(debug|info|warn|error)$`},
//	})
//	if errs := schema.Validate(envMap); errs != nil {
//		for _, e := range errs { log.Println(e) }
//	}
package env
