package env

import (
	"fmt"
	"strings"
)

// Label represents a key-value metadata annotation attached to an env entry.
type Label struct {
	Key   string
	Value string
}

// Labeler manages a set of labels associated with environment variable keys.
type Labeler struct {
	entries map[string][]Label
}

// NewLabeler returns an initialised Labeler.
func NewLabeler() *Labeler {
	return &Labeler{entries: make(map[string][]Label)}
}

// Set attaches a label to the given env key. Duplicate label keys are
// overwritten for the same env key.
func (l *Labeler) Set(envKey, labelKey, labelValue string) {
	envKey = strings.ToUpper(envKey)
	for i, lbl := range l.entries[envKey] {
		if strings.EqualFold(lbl.Key, labelKey) {
			l.entries[envKey][i].Value = labelValue
			return
		}
	}
	l.entries[envKey] = append(l.entries[envKey], Label{Key: labelKey, Value: labelValue})
}

// Get returns the value for a label key attached to an env key, and whether
// it was found.
func (l *Labeler) Get(envKey, labelKey string) (string, bool) {
	envKey = strings.ToUpper(envKey)
	for _, lbl := range l.entries[envKey] {
		if strings.EqualFold(lbl.Key, labelKey) {
			return lbl.Value, true
		}
	}
	return "", false
}

// All returns every label attached to the given env key.
func (l *Labeler) All(envKey string) []Label {
	return l.entries[strings.ToUpper(envKey)]
}

// Delete removes all labels for the given env key.
func (l *Labeler) Delete(envKey string) {
	delete(l.entries, strings.ToUpper(envKey))
}

// Annotate returns a copy of the supplied map with each value annotated by
// appending any labels as a comment-style suffix (useful for debug output).
func (l *Labeler) Annotate(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		labels := l.All(k)
		if len(labels) == 0 {
			out[k] = v
			continue
		}
		parts := make([]string, len(labels))
		for i, lbl := range labels {
			parts[i] = fmt.Sprintf("%s=%s", lbl.Key, lbl.Value)
		}
		out[k] = fmt.Sprintf("%s # %s", v, strings.Join(parts, ", "))
	}
	return out
}
