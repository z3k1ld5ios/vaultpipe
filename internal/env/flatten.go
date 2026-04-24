package env

import (
	"fmt"
	"sort"
	"strings"
)

// Flattener converts a nested map[string]any into a flat map[string]string
// using a configurable separator and optional key prefix.
type Flattener struct {
	separator string
	prefix    string
}

// NewFlattener returns a Flattener with the given separator.
// A common separator is "__" for environment variable compatibility.
func NewFlattener(separator string) *Flattener {
	if separator == "" {
		separator = "__"
	}
	return &Flattener{separator: separator}
}

// WithPrefix returns a copy of the Flattener that prepends the given prefix
// to every key in the output map.
func (f *Flattener) WithPrefix(prefix string) *Flattener {
	return &Flattener{separator: f.separator, prefix: prefix}
}

// Flatten recursively walks src and returns a flat map where nested keys are
// joined with the configured separator. Non-string leaf values are formatted
// with fmt.Sprintf("%v", v).
func (f *Flattener) Flatten(src map[string]any) (map[string]string, error) {
	out := make(map[string]string)
	if err := f.walk(src, f.prefix, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Keys returns the sorted keys of the flattened map for deterministic output.
func (f *Flattener) Keys(src map[string]any) ([]string, error) {
	flat, err := f.Flatten(src)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(flat))
	for k := range flat {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys, nil
}

func (f *Flattener) walk(node map[string]any, prefix string, out map[string]string) error {
	for k, v := range node {
		if strings.ContainsAny(k, "=\x00") {
			return fmt.Errorf("flatten: invalid character in key %q", k)
		}
		full := k
		if prefix != "" {
			full = prefix + f.separator + k
		}
		switch val := v.(type) {
		case map[string]any:
			if err := f.walk(val, full, out); err != nil {
				return err
			}
		case string:
			out[full] = val
		case nil:
			out[full] = ""
		default:
			out[full] = fmt.Sprintf("%v", val)
		}
	}
	return nil
}
