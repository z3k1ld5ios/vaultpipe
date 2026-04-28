package env

import "fmt"

// Cloner produces deep copies of environment maps with optional key filtering
// and transformation applied during the clone operation.
type Cloner struct {
	filterFn func(key string) bool
	transformFn func(key, value string) (string, string, error)
}

// CloneOption configures a Cloner.
type CloneOption func(*Cloner)

// WithCloneFilter restricts which keys are included in the cloned map.
func WithCloneFilter(fn func(key string) bool) CloneOption {
	return func(c *Cloner) {
		c.filterFn = fn
	}
}

// WithCloneTransform applies a key/value transformation during cloning.
func WithCloneTransform(fn func(key, value string) (string, string, error)) CloneOption {
	return func(c *Cloner) {
		c.transformFn = fn
	}
}

// NewCloner constructs a Cloner with the given options.
func NewCloner(opts ...CloneOption) *Cloner {
	c := &Cloner{}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Clone returns a deep copy of src, applying any configured filter and transform.
// Returns an error if the transform function fails for any key.
func (c *Cloner) Clone(src map[string]string) (map[string]string, error) {
	if src == nil {
		return map[string]string{}, nil
	}
	out := make(map[string]string, len(src))
	for k, v := range src {
		if c.filterFn != nil && !c.filterFn(k) {
			continue
		}
		outKey, outVal := k, v
		if c.transformFn != nil {
			var err error
			outKey, outVal, err = c.transformFn(k, v)
			if err != nil {
				return nil, fmt.Errorf("clone transform failed for key %q: %w", k, err)
			}
		}
		out[outKey] = outVal
	}
	return out, nil
}

// MustClone is like Clone but panics on error.
func (c *Cloner) MustClone(src map[string]string) map[string]string {
	result, err := c.Clone(src)
	if err != nil {
		panic(err)
	}
	return result
}
