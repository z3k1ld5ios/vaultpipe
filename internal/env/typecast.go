package env

import (
	"fmt"
	"strconv"
	"strings"
)

// TypeCaster converts string values from secret maps into typed Go values.
// It is useful when downstream consumers expect non-string types (e.g. int, bool, float64).
type TypeCaster struct{}

// NewTypeCaster returns a new TypeCaster.
func NewTypeCaster() *TypeCaster {
	return &TypeCaster{}
}

// AsString returns the value as-is.
func (t *TypeCaster) AsString(m map[string]string, key string) (string, error) {
	v, ok := m[key]
	if !ok {
		return "", fmt.Errorf("typecast: key %q not found", key)
	}
	return v, nil
}

// AsInt parses the value as a base-10 integer.
func (t *TypeCaster) AsInt(m map[string]string, key string) (int64, error) {
	v, err := t.AsString(m, key)
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("typecast: key %q cannot be parsed as int: %w", key, err)
	}
	return n, nil
}

// AsBool parses the value as a boolean (accepts 1/0, true/false, yes/no, on/off).
func (t *TypeCaster) AsBool(m map[string]string, key string) (bool, error) {
	v, err := t.AsString(m, key)
	if err != nil {
		return false, err
	}
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "yes", "on":
		return true, nil
	case "0", "false", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("typecast: key %q value %q is not a recognised boolean", key, v)
	}
}

// AsFloat parses the value as a 64-bit floating-point number.
func (t *TypeCaster) AsFloat(m map[string]string, key string) (float64, error) {
	v, err := t.AsString(m, key)
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
	if err != nil {
		return 0, fmt.Errorf("typecast: key %q cannot be parsed as float: %w", key, err)
	}
	return f, nil
}
