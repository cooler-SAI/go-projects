package main

import (
	"fmt"
	"sync"
)

func withoutMutexNew() {
	var wg sync.WaitGroup
	numGoroutines := 1000
	counter := 0

	fmt.Printf("Running demo without a mutex. Launching %d goroutines.\n", numGoroutines)
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

func main() {
	withoutMutexNew()

}
