package main

import (
	"fmt"
	"sync"
)

func main() {

	var wg sync.WaitGroup

	var safeMap sync.Map

	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key_%d", id)
			value := fmt.Sprintf("value_%d", id)
			safeMap.Store(key, value)
			fmt.Printf("Goroutine writer %d: stored {%s: %s}\n", id, key, value)
		}(i)
	}

	wg.Wait()
	fmt.Println("\nAll writers finished. Starting readers...")

	wg.Add(3)
	for i := 0; i < 3; i++ {

		go func(id int) {
			defer wg.Done()

			// Try to read all keys
			for j := 0; j < 5; j++ {
				key := fmt.Sprintf("key_%d", j)
				if val, ok := safeMap.Load(key); ok {
					fmt.Printf("Reader %d: found %s for key %s\n", id, val, key)
				} else {
					fmt.Printf("Reader %d: key %s not found\n", id, key)
				}
			}

			// Alternative: iterate through all elements
			fmt.Printf("Reader %d: iterating all elements:\n", id)
			safeMap.Range(func(key, value any) bool {
				fmt.Printf("Reader %d: {%v: %v}\n", id, key, value)
				return true
			})

		}(i)
	}
	wg.Wait()
	fmt.Println("\nAll job finished. Gratz")

}
