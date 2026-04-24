package env

import (
	"fmt"
	"regexp"
)

// SchemaType represents the expected type for an environment variable.
type SchemaType string

const (
	TypeString SchemaType = "string"
	TypeInt    SchemaType = "int"
	TypeBool   SchemaType = "bool"
)

// FieldSchema describes a single expected environment variable.
type FieldSchema struct {
	Key      string
	Type     SchemaType
	Required bool
	Pattern  string // optional regex pattern
}

// Schema holds a collection of field definitions for validation.
type Schema struct {
	fields []FieldSchema
}

// NewSchema constructs a Schema from a slice of field definitions.
func NewSchema(fields []FieldSchema) *Schema {
	return &Schema{fields: fields}
}

// Validate checks the provided env map against the schema.
// It returns a list of validation errors, or nil if all checks pass.
func (s *Schema) Validate(env map[string]string) []error {
	var errs []error

	for _, f := range s.fields {
		val, ok := env[f.Key]

		if !ok || val == "" {
			if f.Required {
				errs = append(errs, fmt.Errorf("required key %q is missing or empty", f.Key))
			}
			continue
		}

		switch f.Type {
		case TypeInt:
			if !regexp.MustCompile(`^-?\d+$`).MatchString(val) {
				errs = append(errs, fmt.Errorf("key %q: expected int, got %q", f.Key, val))
			}
		case TypeBool:
			if !regexp.MustCompile(`^(true|false|1|0)$`).MatchString(val) {
				errs = append(errs, fmt.Errorf("key %q: expected bool, got %q", f.Key, val))
			}
		}

		if f.Pattern != "" {
			matched, err := regexp.MatchString(f.Pattern, val)
			if err != nil {
				errs = append(errs, fmt.Errorf("key %q: invalid pattern %q: %w", f.Key, f.Pattern, err))
			} else if !matched {
				errs = append(errs, fmt.Errorf("key %q: value %q does not match pattern %q", f.Key, val, f.Pattern))
			}
		}
	}

	if len(errs) == 0 {
		return nil
	}
	return errs
}
