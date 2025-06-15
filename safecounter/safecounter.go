package main

import (
	"fmt"
	"sync"
	"time"
)

// SafeCounter - streamSafe App
type SafeCounter struct {
	mu      sync.RWMutex
	counter int
}

func (sc *SafeCounter) IncrementCount() {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.counter++
	fmt.Printf("-> Writing: %d\n", sc.counter)
}

func (sc *SafeCounter) ReadCount() int {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	value := sc.counter
	fmt.Printf("<- Reading: %d\n", value)
	return value
}

func main() {
	counter := SafeCounter{}
	var wg sync.WaitGroup

	// Readers
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 3; j++ {
				val := counter.ReadCount()
				time.Sleep(100 * time.Millisecond)
				fmt.Printf("Reader %d see: %d\n", id, val)
			}
		}(i)
	}

	// Writers
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 2; j++ {
				counter.IncrementCount()
				time.Sleep(200 * time.Millisecond)
				fmt.Printf("Writer %d increase counter\n", id)
			}
		}(i)
	}

	wg.Wait()

	fmt.Printf("Job is done! Gratz! Final number is:%d\n", counter.ReadCount())

}
