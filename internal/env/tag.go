package env

import (
	"fmt"
	"strings"
)

// Tag represents a metadata label attached to an environment key.
type Tag struct {
	Key   string
	Value string
}

// Tagger maintains a registry of tags associated with environment variable keys.
type Tagger struct {
	tags map[string][]Tag
}

// NewTagger returns a new Tagger with an empty tag registry.
func NewTagger() *Tagger {
	return &Tagger{tags: make(map[string][]Tag)}
}

// Set associates a tag key-value pair with the given environment variable name.
// Duplicate tag keys for the same env key are silently overwritten.
func (t *Tagger) Set(envKey, tagKey, tagValue string) {
	envKey = strings.ToUpper(envKey)
	existing := t.tags[envKey]
	for i, tag := range existing {
		if tag.Key == tagKey {
			existing[i].Value = tagValue
			t.tags[envKey] = existing
			return
		}
	}
	t.tags[envKey] = append(existing, Tag{Key: tagKey, Value: tagValue})
}

// Get returns all tags associated with the given environment variable name.
// Returns nil if no tags are registered for that key.
func (t *Tagger) Get(envKey string) []Tag {
	return t.tags[strings.ToUpper(envKey)]
}

// HasTag reports whether the given environment variable has a tag with the
// specified key and value.
func (t *Tagger) HasTag(envKey, tagKey, tagValue string) bool {
	for _, tag := range t.Get(envKey) {
		if tag.Key == tagKey && tag.Value == tagValue {
			return true
		}
	}
	return false
}

// Filter returns a subset of the provided env map whose keys carry the given
// tag key-value pair.
func (t *Tagger) Filter(env map[string]string, tagKey, tagValue string) map[string]string {
	out := make(map[string]string)
	for k, v := range env {
		if t.HasTag(k, tagKey, tagValue) {
			out[k] = v
		}
	}
	return out
}

// Summary returns a human-readable string listing all registered tags for
// the given environment variable.
func (t *Tagger) Summary(envKey string) string {
	tags := t.Get(envKey)
	if len(tags) == 0 {
		return fmt.Sprintf("%s: (no tags)", strings.ToUpper(envKey))
	}
	parts := make([]string, 0, len(tags))
	for _, tag := range tags {
		parts = append(parts, fmt.Sprintf("%s=%s", tag.Key, tag.Value))
	}
	return fmt.Sprintf("%s: [%s]", strings.ToUpper(envKey), strings.Join(parts, ", "))
}
