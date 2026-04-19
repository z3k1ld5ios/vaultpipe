package env

import (
	"fmt"
	"testing"
)

func BenchmarkApplyMap_Large(b *testing.B) {
	tr := NewTransformer()
	secrets := make(map[string]string, 100)
	keys := make([]string, 100)
	for i := 0; i < 100; i++ {
		k := fmt.Sprintf("SECRET_KEY_%d", i)
		secrets[k] = fmt.Sprintf("VALUE_%d", i)
		keys[i] = k
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tr.ApplyMap(secrets, "lower", keys)
	}
}
