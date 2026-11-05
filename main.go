package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	SIZE   = 100_000_000
	CHUNKS = 8

	maxRandomValue = 1_000_000
)

// generateRandomElements создаёт слайс из size случайных чисел в диапазоне [1, maxRandomValue].
func generateRandomElements(size int) []int {
	if size <= 0 {
		return nil
	}

	data := make([]int, size)
	for i := range data {
		data[i] = rand.Intn(maxRandomValue) + 1
	}
	return data
}

// maximum возвращает наибольший элемент среза data или 0, если срез пуст.
func maximum(data []int) int {
	if len(data) == 0 {
		return 0
	}

	max := data[0]
	for _, value := range data[1:] {
		if value > max {
			max = value
		}
	}
	return max
}

// maxChunks находит наибольший элемент среза, распределяя поиск между CHUNKS горутинами.
func maxChunks(data []int) int {
	if len(data) == 0 {
		return 0
	}

	chunkSize := len(data) / CHUNKS
	results := make([]int, CHUNKS)
	for index := range results {
		results[index] = data[0]
	}

	var wg sync.WaitGroup
	wg.Add(CHUNKS)
	for index := 0; index < CHUNKS; index++ {
		start := index * chunkSize
		end := start + chunkSize
		if index == CHUNKS-1 {
			end = len(data)
		}
		go func(resultIndex int, part []int) {
			defer wg.Done()
			if len(part) > 0 {
				// Если часть среза пустая, результат сохраняет начальное значение.
				results[resultIndex] = maximum(part)
			}
		}(index, data[start:end])
	}
	wg.Wait()
	return maximum(results)
}

func main() {
	data := generateRandomElements(SIZE)

	start := time.Now()
	sequentialMax := maximum(data)
	sequentialTime := time.Since(start)

	start = time.Now()
	parallelMax := maxChunks(data)
	parallelTime := time.Since(start)

	fmt.Printf("Размер входа: %d элементов\n", len(data))
	fmt.Printf("maximum: максимум %d, время %d мкс\n", sequentialMax, sequentialTime.Microseconds())
	fmt.Printf("maxChunks (%d горутин): максимум %d, время %d мкс\n", CHUNKS, parallelMax, parallelTime.Microseconds())
	if sequentialMax != parallelMax {
		panic("Результаты поиска не совпали")
	}
}
