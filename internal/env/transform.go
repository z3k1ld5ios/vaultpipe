package env

import (
	"fmt"
	"strings"
)

// TransformFunc is a function that transforms a secret value.
type TransformFunc func(value string) (string, error)

// Transformer applies named transformations to a secret map.
type Transformer struct {
	rules map[string]TransformFunc
}

// NewTransformer returns a Transformer with built-in transforms registered.
func NewTransformer() *Transformer {
	t := &Transformer{rules: make(map[string]TransformFunc)}
	t.Register("upper", func(v string) (string, error) { return strings.ToUpper(v), nil })
	t.Register("lower", func(v string) (string, error) { return strings.ToLower(v), nil })
	t.Register("trim", func(v string) (string, error) { return strings.TrimSpace(v), nil })
	return t
}

// Register adds a named transform function.
func (t *Transformer) Register(name string, fn TransformFunc) {
	t.rules[strings.ToLower(name)] = fn
}

// Apply runs the named transform on the given value.
func (t *Transformer) Apply(name, value string) (string, error) {
	fn, ok := t.rules[strings.ToLower(name)]
	if !ok {
		return "", fmt.Errorf("unknown transform: %q", name)
	}
	return fn(value)
}

// ApplyMap applies a transform to each key listed in keys within the secrets map.
func (t *Transformer) ApplyMap(secrets map[string]string, transform string, keys []string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for _, k := range keys {
		v, ok := out[k]
		if !ok {
			continue
		}
		result, err := t.Apply(transform, v)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		out[k] = result
	}
	return out, nil
}
