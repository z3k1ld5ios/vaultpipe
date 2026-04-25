package env

// Deduplicator removes duplicate keys from an environment map,
// applying a configurable conflict resolution strategy.
//
// When multiple sources contribute the same key, the strategy
// determines which value survives: first-wins or last-wins.
type Deduplicator struct {
	lastWins bool
}

// DedupeOption configures a Deduplicator.
type DedupeOption func(*Deduplicator)

// WithLastWins configures the deduplicator to keep the last value
// seen for a given key. Default behaviour is first-wins.
func WithLastWins() DedupeOption {
	return func(d *Deduplicator) {
		d.lastWins = true
	}
}

// NewDeduplicator returns a Deduplicator with the given options applied.
func NewDeduplicator(opts ...DedupeOption) *Deduplicator {
	d := &Deduplicator{}
	for _, o := range opts {
		o(d)
	}
	return d
}

// Apply removes duplicate keys from each successive map in sources,
// merging them into a single result map. Order of sources matters:
// earlier sources are processed first.
//
// With first-wins (default), the value from the earliest source is kept.
// With last-wins, the value from the latest source is kept.
func (d *Deduplicator) Apply(sources ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, src := range sources {
		for k, v := range src {
			if _, exists := result[k]; !exists || d.lastWins {
				result[k] = v
			}
		}
	}
	return result
}

// Keys returns a deduplicated slice of keys present in the given map.
// The returned slice is sorted for deterministic output.
func Keys(m map[string]string) []string {
	seen := make(map[string]struct{}, len(m))
	out := make([]string, 0, len(m))
	for k := range m {
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			out = append(out, k)
		}
	}
	sortStrings(out)
	return out
}

// sortStrings sorts a string slice in-place (insertion sort for small sets).
func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
