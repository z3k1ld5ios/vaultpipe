package env

import "fmt"

// Version represents a snapshot version of an environment map,
// tracking the generation counter and a short content hash.
type Version struct {
	Generation uint64
	Checksum   string
}

// String returns a human-readable representation of the version.
func (v Version) String() string {
	return fmt.Sprintf("gen=%d chk=%s", v.Generation, v.Checksum)
}

// Versioner wraps an environment map and tracks mutations via
// an incrementing generation counter paired with a checksum.
type Versioner struct {
	current map[string]string
	gen     uint64
}

// NewVersioner creates a Versioner seeded with the provided base map.
// The initial generation is 0.
func NewVersioner(base map[string]string) *Versioner {
	copy := make(map[string]string, len(base))
	for k, v := range base {
		copy[k] = v
	}
	return &Versioner{current: copy}
}

// Apply replaces the internal map with next, bumping the generation
// counter only when the content has actually changed.
// Returns the resulting Version regardless.
func (vr *Versioner) Apply(next map[string]string) Version {
	sum := envChecksum(next)
	if sum != envChecksum(vr.current) {
		vr.gen++
		vr.current = cloneMap(next)
	}
	return Version{Generation: vr.gen, Checksum: sum}
}

// Current returns a copy of the current environment map.
func (vr *Versioner) Current() map[string]string {
	return cloneMap(vr.current)
}

// Version returns the current Version descriptor without mutating state.
func (vr *Versioner) Version() Version {
	return Version{Generation: vr.gen, Checksum: envChecksum(vr.current)}
}

func cloneMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
