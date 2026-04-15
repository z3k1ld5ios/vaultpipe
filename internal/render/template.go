package render

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// Renderer resolves template expressions in values using a secret map.
type Renderer struct {
	secrets map[string]string
}

// NewRenderer creates a Renderer backed by the provided secrets map.
func NewRenderer(secrets map[string]string) *Renderer {
	return &Renderer{secrets: secrets}
}

// RenderValue processes a single value string, replacing any
// {{ vault "KEY" }} directives with the corresponding secret value.
func (r *Renderer) RenderValue(value string) (string, error) {
	if !strings.Contains(value, "{{vault") && !strings.Contains(value, "{{ vault") {
		return value, nil
	}

	funcMap := template.FuncMap{
		"vault": func(key string) (string, error) {
			v, ok := r.secrets[key]
			if !ok {
				return "", fmt.Errorf("secret key %q not found", key)
			}
			return v, nil
		},
	}

	tmpl, err := template.New("").Funcs(funcMap).Parse(value)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return buf.String(), nil
}

// RenderMap applies RenderValue to every value in the provided map,
// returning a new map with all templates resolved.
func (r *Renderer) RenderMap(env map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		resolved, err := r.RenderValue(v)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		out[k] = resolved
	}
	return out, nil
}
