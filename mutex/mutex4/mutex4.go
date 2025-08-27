package main

import (
	"fmt"
	"sync"
)

func withoutRWMutexNew4() {
	numGoroutines := 1000
	counter := 0
	var wg sync.WaitGroup

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			counter++
		}()

	}
	wg.Wait()
	fmt.Printf("Final counter value: %d\n", counter)
	fmt.Println("Expected value: 1000")
}

func withRWMutexNew4() {
	numGoroutines := 1000
	counter := 0
	var rwmu sync.RWMutex
	var wg sync.WaitGroup

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			rwmu.Lock()
			counter++
			rwmu.Unlock()
		}()

	}
	wg.Wait()
	fmt.Printf("Final counter value: %d\n", counter)
	fmt.Println("Expected value: 1000")
}

func main() {
	withoutRWMutexNew4()
	withRWMutexNew4()
}
