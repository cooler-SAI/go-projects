package main

import (
	"fmt"
	"sync"
)

func withoutRWMutex() {
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

func withRWMutex() {
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

func withRWMutexFull() {
	var rwmu sync.RWMutex
	var wg sync.WaitGroup
	counter := 0
	numWriters := 300
	numReaders := 700
	totalGoroutines := numWriters + numReaders

	fmt.Printf("\nRunning demo with a RWMutex. Launching %d goroutines (%d readers, %d writers).\n", totalGoroutines, numReaders, numWriters)
	wg.Add(totalGoroutines)
	for i := 0; i < numWriters; i++ {
		go func(i int) {
			defer wg.Done()
			rwmu.Lock()
			counter++
			rwmu.Unlock()

		}(i)
	}

	for i := 0; i < numReaders; i++ {
		go func(id int) {
			defer wg.Done()
			rwmu.RLock()
			value := counter
			fmt.Printf("Reader %d: read value %d\n", id, value)
			rwmu.RUnlock()
		}(i)
	}
	wg.Wait()
	fmt.Printf("Final counter value: %d\n", counter)
	fmt.Println("Expected value: 300")
	fmt.Println("RWMutex allows many concurrent readers but only one writer, making it more efficient for read-heavy workloads.")

}

func Test() {
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(workerID int) {
			fmt.Printf("Worker #%d started work\n", workerID)
			defer wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("All workers finished work")
}

func main() {
	withoutRWMutex()
	withRWMutex()

	withRWMutexFull()
	Test()
}
