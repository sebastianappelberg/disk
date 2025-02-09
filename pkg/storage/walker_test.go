package storage

import (
	"testing"
)

func BenchmarkAnalyze(b *testing.B) {
	analyzer := NewFileWalker[File](WithDecisionFilter[File](IdentityFilter), WithMapper(IdentityMapper))
	counter := int64(0)
	for i := 0; i < b.N; i++ {
		files := analyzer.GetFiles("C:\\Users\\sebastian\\src")
		for file := range files {
			counter += file.Size
		}
	}
}

func TestFilterDecision_Includes(t *testing.T) {
	decision := Include | Skip

	if !decision.Includes(Include) {
		t.Error("Includes fail")
	}
	if !decision.Includes(Skip) {
		t.Error("Includes fail")
	}
	if !decision.Includes(Include | Skip) {
		t.Error("Includes fail")
	}
	if decision.Includes(ShortCircuit) {
		t.Error("Includes fail")
	}
}
