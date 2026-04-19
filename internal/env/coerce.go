package env

import (
	"fmt"
	"strconv"
	"strings"
)

// CoerceType describes how a secret value should be coerced before injection.
type CoerceType string

const (
	CoerceString CoerceType = "string"
	CoerceBool   CoerceType = "bool"
	CoerceInt    CoerceType = "int"
	CoerceUpper  CoerceType = "upper"
	CoerceLower  CoerceType = "lower"
)

// Coercer applies type coercions to a map of environment values.
type Coercer struct {
	rules map[string]CoerceType
}

// NewCoercer creates a Coercer with the given key→type rules.
func NewCoercer(rules map[string]CoerceType) *Coercer {
	return &Coercer{rules: rules}
}

// Apply returns a new map with coercions applied to matching keys.
func (c *Coercer) Apply(env map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}
	for key, ct := range c.rules {
		v, ok := out[key]
		if !ok {
			continue
		}
		coerced, err := coerceValue(v, ct)
		if err != nil {
			return nil, fmt.Errorf("coerce %q as %s: %w", key, ct, err)
		}
		out[key] = coerced
	}
	return out, nil
}

func coerceValue(v string, ct CoerceType) (string, error) {
	switch ct {
	case CoerceString:
		return v, nil
	case CoerceBool:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return "", fmt.Errorf("cannot parse %q as bool", v)
		}
		return strconv.FormatBool(b), nil
	case CoerceInt:
		_, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return "", fmt.Errorf("cannot parse %q as int", v)
		}
		return v, nil
	case CoerceUpper:
		return strings.ToUpper(v), nil
	case CoerceLower:
		return strings.ToLower(v), nil
	default:
		return "", fmt.Errorf("unknown coerce type %q", ct)
	}
}
