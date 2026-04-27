package env

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

// envChecksum produces a deterministic SHA-256 digest of the key=value
// pairs in m, sorted lexicographically by key so that map iteration
// order does not affect the result.
func envChecksum(m map[string]string) string {
	if len(m) == 0 {
		return "e3b0c44298fc1c14" // well-known empty hash prefix
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(m[k])
		sb.WriteByte('\n')
	}

	h := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(h[:8]) // first 8 bytes → 16 hex chars
}
