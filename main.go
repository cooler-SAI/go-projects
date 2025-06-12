package main

import (
	"fmt"
	"runtime"
	"sync"
)

var counterWithoutMutex int

var (
	counterWithMutex int
	mu               sync.Mutex
)

const (
	numGoroutines   = 1000
	incrementsPerGo = 1000
)

func incrementorNoMutex(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < incrementsPerGo; i++ {
		counterWithoutMutex++
	}
}

func incrementorMutex(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < incrementsPerGo; i++ {
		mu.Lock()
		counterWithMutex++
		mu.Unlock()
	}
}

func main() {
	fmt.Println("CPU cores available:", runtime.NumCPU())
	fmt.Println("Initial active goroutines:", runtime.NumGoroutine())
	fmt.Println("--------------------------------------------------")

	fmt.Println("Running counter WITHOUT mutex (expect race condition)...")
	counterWithoutMutex = 0
	var wgRace sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wgRace.Add(1)
		go incrementorNoMutex(&wgRace)
	}

	wgRace.Wait()

	fmt.Printf("Final counter value (without mutex): %d\n", counterWithoutMutex)
	expectedValue := numGoroutines * incrementsPerGo
	fmt.Printf("Expected value: %d (but likely less due to race condition)\n", expectedValue)
	fmt.Println("--------------------------------------------------")

	fmt.Println("Running counter WITH mutex (expect correct result)...")
	counterWithMutex = 0
	var wgMutex sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wgMutex.Add(1)
		go incrementorMutex(&wgMutex)
	}

	wgMutex.Wait()

	fmt.Printf("Final counter value (with mutex): %d\n", counterWithMutex)
	fmt.Printf("Expected value: %d (should be correct)\n", expectedValue)
	fmt.Println("--------------------------------------------------")

	fmt.Println("Total active goroutines at end:", runtime.NumGoroutine())
	fmt.Println("Demonstration finished.")
}
