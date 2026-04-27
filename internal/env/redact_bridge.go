package env

import (
	"github.com/yourusername/vaultpipe/internal/redact"
)

// RedactBridge applies a redactor to an environment map, replacing any known
// secret values with a redacted placeholder before the map is used in logs,
// output, or audit events. It does not mutate the input map.
type RedactBridge struct {
	redactor *redact.Redactor
}

// NewRedactBridge creates a RedactBridge backed by the given Redactor.
func NewRedactBridge(r *redact.Redactor) *RedactBridge {
	return &RedactBridge{redactor: r}
}

// Apply returns a copy of env with all registered secret values replaced by
// the redactor's placeholder string. Keys are never redacted.
func (b *RedactBridge) Apply(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = b.redactor.Redact(v)
	}
	return out
}

// RegisterSecrets registers all values from the provided secrets map with the
// underlying redactor so that subsequent Apply calls will mask them.
func (b *RedactBridge) RegisterSecrets(secrets map[string]string) {
	for _, v := range secrets {
		if v != "" {
			b.redactor.Add(v)
		}
	}
}

// RedactMap is a convenience function that builds a one-shot RedactBridge,
// registers the provided secrets, and returns a redacted copy of env.
func RedactMap(env map[string]string, secrets map[string]string) map[string]string {
	r := redact.New()
	b := NewRedactBridge(r)
	b.RegisterSecrets(secrets)
	return b.Apply(env)
}
