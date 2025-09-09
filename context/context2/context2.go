package main

import (
	"errors" // Package errors for creating custom errors
	"fmt"
	"math/rand" // For generating random numbers
	"sync"
	"time"
)

// workerResult - structure for returning result or error from a goroutine.
type workerResult struct {
	id    int
	value string
	err   error
}

// worker simulates a task that might succeed or fail.
// It sends its result (or error) to a results channel.
func worker(id int, results chan<- workerResult, wg *sync.WaitGroup) {
	defer wg.Done() // Decrement WaitGroup counter when goroutine completes

	fmt.Printf("Worker %d: Starting task...\n", id)
	time.Sleep(time.Duration(rand.Intn(500)+100) * time.Millisecond) // Simulate work

	// Simulate random error
	if rand.Intn(100) < 30 { // 30% chance of error
		err := errors.New(fmt.Sprintf("Worker %d: Failed due to random error", id))
		results <- workerResult{id: id, err: err} // Send result with error
		fmt.Printf("Worker %d: Task failed.\n", id)
		return // Important: exit function after sending error
	}

	value := fmt.Sprintf("Data from Worker %d", id)
	results <- workerResult{id: id, value: value} // Send successful result
	fmt.Printf("Worker %d: Task completed successfully.\n", id)
}

func main() {
	fmt.Println("Starting demonstration of error handling in concurrent programs...")

	var wg sync.WaitGroup
	const numWorkers = 5 // Number of worker goroutines

	// Channel for collecting results and errors from all worker goroutines
	resultsChan := make(chan workerResult, numWorkers) // Buffered channel

	// Start worker goroutines
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, resultsChan, &wg)
	}

	// Goroutine to close results channel after all worker goroutines complete
	go func() {
		wg.Wait()          // Wait for all worker goroutines to complete
		close(resultsChan) // Close results channel, signaling no more data
		fmt.Println("All worker goroutines completed, results channel closed.")
	}()

	// Collect and process results/errors from the channel
	var successfulResults []string
	var failedResults []error

	for res := range resultsChan { // Read from channel until it's closed
		if res.err != nil {
			fmt.Printf("Main: Received error from Worker %d: %v\n", res.id, res.err)
			failedResults = append(failedResults, res.err)
		} else {
			fmt.Printf("Main: Received successful result from Worker %d: %s\n", res.id, res.value)
			successfulResults = append(successfulResults, res.value)
		}
	}

	fmt.Println("\n--- Execution Summary ---")
	fmt.Printf("Successfully processed results: %d\n", len(successfulResults))
	for _, val := range successfulResults {
		fmt.Printf("  - %s\n", val)
	}
	fmt.Printf("Number of errors: %d\n", len(failedResults))
	for _, err := range failedResults {
		fmt.Printf("  - %v\n", err)
	}

	fmt.Println("\nError handling demonstration completed.")
}
