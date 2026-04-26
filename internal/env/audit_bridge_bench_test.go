package env

import (
	"fmt"
	"testing"
)

func BenchmarkBuildAuditSummary_Large(b *testing.B) {
	before := make(map[string]string, 500)
	after := make(map[string]string, 500)

	for i := 0; i < 400; i++ {
		k := fmt.Sprintf("KEY_%04d", i)
		before[k] = fmt.Sprintf("value_%d", i)
		after[k] = fmt.Sprintf("value_%d", i)
	}
	// 50 updates
	for i := 0; i < 50; i++ {
		k := fmt.Sprintf("KEY_%04d", i)
		after[k] = fmt.Sprintf("new_value_%d", i)
	}
	// 50 additions
	for i := 400; i < 450; i++ {
		after[fmt.Sprintf("NEW_%04d", i)] = "added"
	}
	// 50 removals — just don't add them to after
	for i := 350; i < 400; i++ {
		delete(after, fmt.Sprintf("KEY_%04d", i))
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = BuildAuditSummary(before, after)
	}
}
