// Package mask provides configurable secret masking policies for log output and display.
package mask

import "strings"

// Level controls how aggressively values are masked.
type Level int

const (
	// LevelNone disables masking.
	LevelNone Level = iota
	// LevelPartial shows first and last characters.
	LevelPartial
	// LevelFull replaces the entire value with a placeholder.
	LevelFull
)

// Policy defines masking behaviour for secret values.
type Policy struct {
	Level       Level
	Placeholder string
}

// DefaultPolicy returns a sensible default masking policy.
func DefaultPolicy() Policy {
	return Policy{
		Level:       LevelPartial,
		Placeholder: "***",
	}
}

// Apply masks a secret value according to the policy.
func (p Policy) Apply(value string) string {
	if len(value) == 0 {
		return value
	}
	switch p.Level {
	case LevelNone:
		return value
	case LevelFull:
		return p.Placeholder
	case LevelPartial:
		if len(value) <= 4 {
			return p.Placeholder
		}
		return string(value[0]) + strings.Repeat("*", len(value)-2) + string(value[len(value)-1])
	default:
		return p.Placeholder
	}
}

// ApplyMap masks all values in a map according to the policy.
func (p Policy) ApplyMap(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = p.Apply(v)
	}
	return out
}
