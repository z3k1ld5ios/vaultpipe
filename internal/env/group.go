package env

import "fmt"

// Group organizes a flat env map into named logical groups based on key prefixes.
// Keys not matching any group prefix are collected under the "default" group.
type Group struct {
	prefixes map[string]string // prefix -> group name
}

// GroupResult holds the grouped keys and their values.
type GroupResult struct {
	Groups map[string]map[string]string
}

// NewGroup creates a Group with the given prefix-to-name mappings.
// Example: {"DB_": "database", "CACHE_": "cache"}
func NewGroup(prefixes map[string]string) *Group {
	p := make(map[string]string, len(prefixes))
	for k, v := range prefixes {
		p[k] = v
	}
	return &Group{prefixes: p}
}

// Apply partitions env into named groups based on registered prefixes.
// Keys matching multiple prefixes are assigned to the longest matching prefix.
// Unmatched keys go into the "default" group.
func (g *Group) Apply(env map[string]string) (*GroupResult, error) {
	if env == nil {
		return nil, fmt.Errorf("env: group: input map must not be nil")
	}

	result := &GroupResult{
		Groups: make(map[string]map[string]string),
	}

	for key, val := range env {
		groupName := g.matchGroup(key)
		if _, ok := result.Groups[groupName]; !ok {
			result.Groups[groupName] = make(map[string]string)
		}
		result.Groups[groupName][key] = val
	}

	return result, nil
}

// matchGroup returns the group name for a key by finding the longest matching prefix.
func (g *Group) matchGroup(key string) string {
	best := ""
	bestName := "default"
	for prefix, name := range g.prefixes {
		if len(prefix) > len(best) && len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			best = prefix
			bestName = name
		}
	}
	return bestName
}

// Keys returns all keys belonging to the named group, or nil if the group is absent.
func (r *GroupResult) Keys(name string) []string {
	g, ok := r.Groups[name]
	if !ok {
		return nil
	}
	keys := make([]string, 0, len(g))
	for k := range g {
		keys = append(keys, k)
	}
	sortStrings(keys)
	return keys
}
