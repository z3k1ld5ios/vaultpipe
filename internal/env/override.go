package env

// Override applies a set of explicit key-value overrides on top of a base
// environment map. Override values always win, regardless of what is already
// present in base.
type Override struct {
	overrides map[string]string
}

// NewOverride creates an Override with the supplied key-value pairs.
func NewOverride(overrides map[string]string) *Override {
	copy := make(map[string]string, len(overrides))
	for k, v := range overrides {
		copy[k] = v
	}
	return &Override{overrides: copy}
}

// Apply merges overrides into base, returning a new map. The original maps are
// not modified.
func (o *Override) Apply(base map[string]string) map[string]string {
	out := make(map[string]string, len(base)+len(o.overrides))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range o.overrides {
		out[k] = v
	}
	return out
}

// Keys returns the set of keys that will be overridden.
func (o *Override) Keys() []string {
	keys := make([]string, 0, len(o.overrides))
	for k := range o.overrides {
		keys = append(keys, k)
	}
	return keys
}

// Len returns the number of override entries.
func (o *Override) Len() int {
	return len(o.overrides)
}

// Get returns the override value for the given key and whether it was found.
func (o *Override) Get(key string) (string, bool) {
	v, ok := o.overrides[key]
	return v, ok
}
