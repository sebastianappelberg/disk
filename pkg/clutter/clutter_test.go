package clutter

import (
	"testing"
	"time"
)

func BenchmarkAnalyze(b *testing.B) {
	minAge := time.Now().AddDate(0, 0, -90)
	clutterAnalyzer := NewAnalyzer(
		WithSizeFilter(50),
		WithMinAgeFilter(minAge),
	)
	counter := int64(0)
	for i := 0; i < b.N; i++ {
		files := clutterAnalyzer.Analyze("C:\\Users\\sebastian\\src")
		for _, file := range files {
			counter += file.Size
		}
	}
}
