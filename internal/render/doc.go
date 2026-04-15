// Package render provides lightweight Go template rendering for vaultpipe.
//
// It allows environment variable values defined in configuration to reference
// Vault secret fields using the {{ vault "KEY" }} template directive.
// The Renderer resolves these directives against a pre-fetched secrets map,
// enabling dynamic composition of env values without additional Vault calls.
package render
