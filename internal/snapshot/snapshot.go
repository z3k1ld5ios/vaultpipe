// Package snapshot captures and compares secret state for change detection and rollback.
package snapshot

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"time"
)

// Snapshot holds a point-in-time copy of resolved secrets.
type Snapshot struct {
	Secrets   map[string]string
	Checksum  string
	CapturedAt time.Time
}

// Take creates a new Snapshot from the given secrets map.
func Take(secrets map[string]string) *Snapshot {
	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}
	return &Snapshot{
		Secrets:    copy,
		Checksum:   checksum(secrets),
		CapturedAt: time.Now(),
	}
}

// Equal returns true if two snapshots have the same checksum.
func (s *Snapshot) Equal(other *Snapshot) bool {
	if s == nil || other == nil {
		return s == other
	}
	return s.Checksum == other.Checksum
}

// Diff returns keys that were added, removed, or changed between s and other.
func (s *Snapshot) Diff(other *Snapshot) map[string]string {
	delta := make(map[string]string)
	for k, v := range other.Secrets {
		if old, ok := s.Secrets[k]; !ok {
			delta[k] = fmt.Sprintf("added: %q", v)
		} else if old != v {
			delta[k] = "changed"
		}
	}
	for k := range s.Secrets {
		if _, ok := other.Secrets[k]; !ok {
			delta[k] = "removed"
		}
	}
	return delta
}

func checksum(secrets map[string]string) string {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s;", k, secrets[k])
	}
	return hex.EncodeToString(h.Sum(nil))
}
