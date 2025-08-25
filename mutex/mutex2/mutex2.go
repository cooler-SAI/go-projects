package main

import (
	"fmt"
	"sync"
	"time"
)

var counter int

func withoutMutex() {

	var wg sync.WaitGroup
	numGoroutines := 1000

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
	fmt.Println("This is a classic race condition. The final value is not 1000 because " +
		"multiple goroutines are trying to increment the counter at the same time.")
}

func withMutex() {
	var mu sync.Mutex
	var wg sync.WaitGroup

	counter := 0
	numGoroutines := 1000

	fmt.Printf("\nRunning demo with a mutex. Launching %d goroutines.\n", numGoroutines)
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Printf("Final counter value: %d\n", counter)
	fmt.Println("Expected value: 1000")
	fmt.Println("The final value is correct because the mutex ensures that only one goroutine at a time can " +
		"access and modify the counter.")

}

func main() {
	withoutMutex()
	time.Sleep(1 * time.Second)
	withMutex()
}
