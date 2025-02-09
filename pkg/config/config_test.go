package config

import "testing"

func BenchmarkGetAppDir(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetAppDir()
	}
}
