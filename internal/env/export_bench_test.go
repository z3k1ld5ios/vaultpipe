package env

import (
	"fmt"
	"testing"
)

func BenchmarkExport_Shell_Large(b *testing.B) {
	env := make(map[string]string, 500)
	for i := 0; i < 500; i++ {
		env[fmt.Sprintf("KEY_%04d", i)] = fmt.Sprintf("value_%d", i)
	}
	exporter := NewExporter(FormatShell, true)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = exporter.Export(env)
	}
}

func BenchmarkExport_JSON_Large(b *testing.B) {
	env := make(map[string]string, 500)
	for i := 0; i < 500; i++ {
		env[fmt.Sprintf("KEY_%04d", i)] = fmt.Sprintf("value_%d", i)
	}
	exporter := NewExporter(FormatJSON, false)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = exporter.Export(env)
	}
}
