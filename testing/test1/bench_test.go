package test1

import "testing"

// Простой бенчмарк для проверки
func BenchmarkSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Простая операция
		_ = i * i
	}
}
