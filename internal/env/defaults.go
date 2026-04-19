package env

// DefaultsApplier merges a set of default key/value pairs into a secret map,
// only setting values for keys that are not already present.
type DefaultsApplier struct {
	defaults map[string]string
}

// NewDefaultsApplier creates a DefaultsApplier with the provided defaults.
func NewDefaultsApplier(defaults map[string]string) *DefaultsApplier {
	d := make(map[string]string, len(defaults))
	for k, v := range defaults {
		d[k] = v
	}
	return &DefaultsApplier{defaults: d}
}

// Apply returns a new map containing all entries from base, with any missing
// keys filled in from the configured defaults.
func (a *DefaultsApplier) Apply(base map[string]string) map[string]string {
	out := make(map[string]string, len(base)+len(a.defaults))
	for k, v := range a.defaults {
		out[k] = v
	}
	for k, v := range base {
		out[k] = v
	}
	return out
}

// Keys returns the list of default keys managed by this applier.
func (a *DefaultsApplier) Keys() []string {
	keys := make([]string, 0, len(a.defaults))
	for k := range a.defaults {
		keys = append(keys, k)
	}
	return keys
}
