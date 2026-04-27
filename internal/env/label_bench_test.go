package env

import (
	"fmt"
	"testing"
)

func BenchmarkAnnotate_Large(b *testing.B) {
	l := NewLabeler()
	m := make(map[string]string, 500)
	for i := 0; i < 500; i++ {
		key := fmt.Sprintf("ENV_VAR_%d", i)
		m[key] = fmt.Sprintf("value_%d", i)
		l.Set(key, "source", "vault")
		l.Set(key, "index", fmt.Sprintf("%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = l.Annotate(m)
	}
}
