package env

import (
	"fmt"
	"testing"
)

func BenchmarkCheckpointer_SaveAndGet_Large(b *testing.B) {
	env := make(map[string]string, 500)
	for i := 0; i < 500; i++ {
		env[fmt.Sprintf("KEY_%d", i)] = fmt.Sprintf("value_%d", i)
	}

	cp := NewCheckpointer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		name := fmt.Sprintf("snap_%d", i)
		cp.Save(name, env)
		cp.Get(name)
	}
}

func BenchmarkCheckpointer_Between_Large(b *testing.B) {
	envA := make(map[string]string, 500)
	envB := make(map[string]string, 500)
	for i := 0; i < 500; i++ {
		envA[fmt.Sprintf("KEY_%d", i)] = fmt.Sprintf("old_%d", i)
		envB[fmt.Sprintf("KEY_%d", i)] = fmt.Sprintf("new_%d", i)
	}

	cp := NewCheckpointer()
	cp.Save("before", envA)
	cp.Save("after", envB)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cp.Between("before", "after")
	}
}
