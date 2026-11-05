package main

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"testing"
)

func TestGenerateRandomElements(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		wantLen int
		wantNil bool
	}{
		{name: "negative", size: -1, wantLen: 0, wantNil: true},
		{name: "zero", size: 0, wantLen: 0, wantNil: true},
		{name: "one", size: 1, wantLen: 1, wantNil: false},
		{name: "ordinary", size: 1000, wantLen: 1000, wantNil: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := generateRandomElements(tt.size)
			if (data == nil) != tt.wantNil {
				t.Fatalf("generateRandomElements(%d) nil = %t, want %t", tt.size, data == nil, tt.wantNil)
			}
			if len(data) != tt.wantLen {
				t.Fatalf("generateRandomElements(%d) length = %d, want %d", tt.size, len(data), tt.wantLen)
			}
			for i, value := range data {
				if value < 1 || value > 1_000_000 {
					t.Errorf("data[%d] = %d, want value in [1, 1000000]", i, value)
				}
			}
			if tt.size > 0 {
				other := generateRandomElements(tt.size)
				if len(other) != tt.wantLen {
					t.Fatalf("second call length = %d, want %d", len(other), tt.wantLen)
				}
				if &data[0] == &other[0] {
					t.Error("separate calls share the same backing array")
				}
			}
		})
	}
}

func TestGenerateRandomElementsRange(t *testing.T) {
	const (
		sampleSize = 2_000_000
		edgeWindow = 1_000
	)

	// Включаем rand.Seed в Go 1.24, чтобы тест был воспроизводимым.
	t.Setenv("GODEBUG", os.Getenv("GODEBUG")+",randseednop=0")
	rand.Seed(1)
	data := generateRandomElements(sampleSize)
	if len(data) != sampleSize {
		t.Fatalf("length = %d, want %d", len(data), sampleSize)
	}

	minValue, maxValue := data[0], data[0]
	for i, value := range data {
		if value < 1 || value > maxRandomValue {
			t.Fatalf("data[%d] = %d, want value in [1, %d]", i, value, maxRandomValue)
		}
		if value < minValue {
			minValue = value
		}
		if value > maxValue {
			maxValue = value
		}
	}
	// Проверяем значения вблизи границ диапазона, так как сами границы могут не попасть в выборку.
	if minValue > edgeWindow {
		t.Errorf("minimum = %d, want value in [1, %d]", minValue, edgeWindow)
	}
	if maxValue <= maxRandomValue-edgeWindow {
		t.Errorf("maximum = %d, want value in [%d, %d]", maxValue, maxRandomValue-edgeWindow+1, maxRandomValue)
	}
}

func TestMaximum(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	minInt := -maxInt - 1
	tests := []struct {
		name string
		data []int
		want int
	}{
		{name: "nil", data: nil, want: 0},
		{name: "empty", data: []int{}, want: 0},
		{name: "one", data: []int{42}, want: 42},
		{name: "equal", data: []int{7, 7, 7}, want: 7},
		{name: "maximum at start", data: []int{9, 3, 1}, want: 9},
		{name: "maximum at end", data: []int{1, 3, 9}, want: 9},
		{name: "maximum in middle", data: []int{3, 9, 1}, want: 9},
		{name: "zeros", data: []int{0, 0, 0}, want: 0},
		{name: "zero and negatives", data: []int{-3, 0, -1}, want: 0},
		{name: "negatives", data: []int{-8, -2, -5}, want: -2},
		{name: "minimum int", data: []int{minInt}, want: minInt},
		{name: "maximum int", data: []int{maxInt}, want: maxInt},
		{name: "int boundaries", data: []int{minInt, 0, maxInt, -1}, want: maxInt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var before []int
			if tt.data != nil {
				before = make([]int, len(tt.data))
				copy(before, tt.data)
			}

			if got := maximum(tt.data); got != tt.want {
				t.Errorf("maximum(%v) = %d, want %d", before, got, tt.want)
			}
			if !reflect.DeepEqual(tt.data, before) {
				t.Errorf("maximum changed input: got %v, want %v", tt.data, before)
			}
		})
	}
}

func TestMaxChunks(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	minInt := -maxInt - 1
	tests := []struct {
		name string
		data []int
		want int
	}{
		{name: "nil", data: nil, want: 0},
		{name: "empty", data: []int{}, want: 0},
		{name: "one", data: []int{42}, want: 42},
		{name: "two", data: []int{-8, -2}, want: -2},
		{name: "seven", data: []int{-8, -2, -5, -3, -6, -4, -9}, want: -2},
		{name: "eight", data: []int{1, 2, 3, 4, 5, 6, 7, 8}, want: 8},
		{name: "nine", data: []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, want: 9},
		{name: "fifteen", data: []int{-15, -14, -13, -12, -11, -10, -9, -8, -7, -6, -5, -4, -3, -2, -1}, want: -1},
		{name: "sixteen", data: []int{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}, want: 16},
		{name: "seventeen", data: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17}, want: 17},
		{name: "zeros", data: []int{0, 0, 0, 0, 0, 0, 0, 0, 0}, want: 0},
		{name: "zero and negatives", data: []int{-9, -8, -7, -6, -5, -4, -3, -2, 0}, want: 0},
		{name: "minimum int", data: []int{minInt}, want: minInt},
		{name: "int boundaries", data: []int{minInt, 0, -1, -2, -3, -4, -5, -6, maxInt}, want: maxInt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var before []int
			if tt.data != nil {
				before = make([]int, len(tt.data))
				copy(before, tt.data)
			}
			if got := maxChunks(tt.data); got != tt.want {
				t.Errorf("maxChunks(%v) = %d, want %d", before, got, tt.want)
			}
			if !reflect.DeepEqual(tt.data, before) {
				t.Errorf("maxChunks changed input: got %v, want %v", tt.data, before)
			}
		})
	}

	// Проверяем случаи, когда максимум находится на границе части или в оставшихся элементах.
	for _, size := range []int{1, 2, 7, 8, 9, 15, 16, 17} {
		for position := 0; position < size; position++ {
			t.Run(fmt.Sprintf("length_%d/maximum_at_%d", size, position), func(t *testing.T) {
				data := make([]int, size)
				for i := range data {
					data[i] = -9
				}
				data[position] = -1
				before := append([]int(nil), data...)
				if got := maxChunks(data); got != -1 {
					t.Errorf("maxChunks(%v) = %d, want -1", before, got)
				}
				if !reflect.DeepEqual(data, before) {
					t.Errorf("maxChunks changed input: got %v, want %v", data, before)
				}
			})
		}
	}
}

func TestMaxChunksMatchesMaximum(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for _, size := range []int{0, 1, 2, 7, 8, 9, 15, 16, 17, 31, 64, 127} {
		t.Run(fmt.Sprintf("length_%d", size), func(t *testing.T) {
			data := make([]int, size)
			for i := range data {
				data[i] = rng.Intn(2001) - 1000
			}
			before := make([]int, size)
			copy(before, data)
			want := maximum(data)
			if got := maxChunks(data); got != want {
				t.Errorf("maxChunks(%v) = %d, maximum = %d", before, got, want)
			}
			if !reflect.DeepEqual(data, before) {
				t.Errorf("maxChunks changed input: got %v, want %v", data, before)
			}
		})
	}
}

var benchmarkMaximumResult int

func BenchmarkMaximum(b *testing.B) {
	benchmarkSearch(b, maximum)
}

func BenchmarkMaxChunks(b *testing.B) {
	benchmarkSearch(b, maxChunks)
}

func benchmarkSearch(b *testing.B, search func([]int) int) {
	for _, size := range []int{1_000, 1_000_000, 10_000_000} {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			rng := rand.New(rand.NewSource(42))
			data := make([]int, size)
			for i := range data {
				data[i] = rng.Intn(maxRandomValue) + 1
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				benchmarkMaximumResult = search(data)
			}
		})
	}
}
