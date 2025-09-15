package test1

import "testing"

func BenchmarkSum(b *testing.B) {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Sum(numbers)
	}
}

func BenchmarkSumSmallSlice(b *testing.B) {
	numbers := []int{1, 2, 3}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Sum(numbers)
	}
}

func BenchmarkSumLargeSlice(b *testing.B) {
	numbers := make([]int, 1000)
	for i := range numbers {
		numbers[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Sum(numbers)
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		expected int
	}{
		{"positive numbers", []int{1, 2, 3, 4, 5}, 15},
		{"with zero", []int{0, 1, 2, 3}, 6},
		{"negative numbers", []int{-1, -2, -3}, -6},
		{"empty slice", []int{}, 0},
		{"single element", []int{42}, 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := Sum(tt.numbers)
			if actual != tt.expected {
				t.Errorf("Sum(%v): expected %d, got %d", tt.numbers, tt.expected, actual)
			}
		})
	}
}

func TestSumEdgeCases(t *testing.T) {
	if Sum(nil) != 0 {
		t.Error("Sum(nil) should return 0")
	}

	if Sum([]int{}) != 0 {
		t.Error("Sum([]int{}) should return 0")
	}
}
