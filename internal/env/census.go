package env

import (
	"fmt"
	"sort"
)

// Census collects statistics about a secret or environment map,
// such as key counts, empty values, and duplicate detection across sources.
type Census struct {
	TotalKeys    int
	EmptyValues  int
	UniqueKeys   []string
	DuplicateKeys []string
}

// Analyze inspects a map and returns a Census report.
func Analyze(m map[string]string) Census {
	keys := make([]string, 0, len(m))
	empty := 0

	for k, v := range m {
		keys = append(keys, k)
		if v == "" {
			empty++
		}
	}

	sort.Strings(keys)

	return Census{
		TotalKeys:   len(m),
		EmptyValues: empty,
		UniqueKeys:  keys,
	}
}

// AnalyzeMultiple inspects multiple sources and returns a Census that includes
// keys that appear in more than one source.
func AnalyzeMultiple(sources ...map[string]string) Census {
	counts := make(map[string]int)
	values := make(map[string]string)

	for _, src := range sources {
		for k, v := range src {
			counts[k]++
			values[k] = v
		}
	}

	var unique, dupes []string
	empty := 0

	for k, c := range counts {
		if c > 1 {
			dupes = append(dupes, k)
		} else {
			unique = append(unique, k)
		}
		if values[k] == "" {
			empty++
		}
	}

	sort.Strings(unique)
	sort.Strings(dupes)

	return Census{
		TotalKeys:     len(counts),
		EmptyValues:   empty,
		UniqueKeys:    unique,
		DuplicateKeys: dupes,
	}
}

// Summary returns a human-readable one-line summary of the census.
func (c Census) Summary() string {
	return fmt.Sprintf("total=%d empty=%d duplicates=%d",
		c.TotalKeys, c.EmptyValues, len(c.DuplicateKeys))
}
