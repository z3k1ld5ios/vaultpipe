package env

import (
	"fmt"
	"os"
	"strings"
)

// SourcePriority defines the precedence order when merging from multiple sources.
type SourcePriority int

const (
	PriorityLow    SourcePriority = iota // e.g. defaults
	PriorityNormal                        // e.g. vault secrets
	PriorityHigh                          // e.g. explicit overrides
)

// Source represents a named provider of environment key-value pairs.
type Source struct {
	Name     string
	Priority SourcePriority
	Values   map[string]string
}

// Sourcer merges multiple Sources according to priority, with higher-priority
// sources winning on key conflicts.
type Sourcer struct {
	sources []Source
}

// NewSourcer returns a Sourcer with the given sources. Sources are sorted
// internally by priority (ascending) so higher-priority values overwrite lower
// ones during merge.
func NewSourcer(sources ...Source) *Sourcer {
	return &Sourcer{sources: sources}
}

// FromOS returns a Source populated from the current process environment.
func FromOS(name string, priority SourcePriority) Source {
	vals := make(map[string]string)
	for _, entry := range os.Environ() {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 2 {
			vals[parts[0]] = parts[1]
		}
	}
	return Source{Name: name, Priority: priority, Values: vals}
}

// Merge combines all registered sources into a single map. Sources with higher
// priority overwrite keys from lower-priority sources. Returns an error if any
// source has a nil Values map.
func (s *Sourcer) Merge() (map[string]string, error) {
	// Sort by priority ascending so higher-priority sources are applied last.
	sorted := make([]Source, len(s.sources))
	copy(sorted, s.sources)
	for i := 1; i < len(sorted); i++ {
		for j := i; j > 0 && sorted[j].Priority < sorted[j-1].Priority; j-- {
			sorted[j], sorted[j-1] = sorted[j-1], sorted[j]
		}
	}

	result := make(map[string]string)
	for _, src := range sorted {
		if src.Values == nil {
			return nil, fmt.Errorf("sourcer: source %q has nil values map", src.Name)
		}
		for k, v := range src.Values {
			result[k] = v
		}
	}
	return result, nil
}

// Names returns the names of all registered sources in registration order.
func (s *Sourcer) Names() []string {
	names := make([]string, len(s.sources))
	for i, src := range s.sources {
		names[i] = src.Name
	}
	return names
}
