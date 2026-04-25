package env

import (
	"fmt"
	"sort"
)

// Diff represents the changes between two environment maps.
type Diff struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string][2]string // key -> [old, new]
}

// HasChanges returns true if the diff contains any additions, removals, or changes.
func (d *Diff) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}

// Summary returns a human-readable summary of the diff.
func (d *Diff) Summary() string {
	return fmt.Sprintf("added=%d removed=%d changed=%d",
		len(d.Added), len(d.Removed), len(d.Changed))
}

// ChangedKeys returns a sorted slice of keys that were added, removed, or changed.
func (d *Diff) ChangedKeys() []string {
	seen := make(map[string]struct{})
	for k := range d.Added {
		seen[k] = struct{}{}
	}
	for k := range d.Removed {
		seen[k] = struct{}{}
	}
	for k := range d.Changed {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// DiffMaps computes the difference between two environment maps.
// prev is the previous state; next is the new state.
func DiffMaps(prev, next map[string]string) *Diff {
	d := &Diff{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
	}

	for k, nv := range next {
		if pv, ok := prev[k]; !ok {
			d.Added[k] = nv
		} else if pv != nv {
			d.Changed[k] = [2]string{pv, nv}
		}
	}

	for k, pv := range prev {
		if _, ok := next[k]; !ok {
			d.Removed[k] = pv
		}
	}

	return d
}
