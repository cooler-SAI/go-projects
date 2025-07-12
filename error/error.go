package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type workerResult struct {
	id    int
	value string
	err   error
}

func worker(id int, results chan<- workerResult, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d: Starting task...\n", id)
	time.Sleep(time.Duration(rand.Intn(500)+100) * time.Millisecond)

	if rand.Intn(100) < 30 {
		err := errors.New(fmt.Sprintf("Worker %d: Failed due to random error", id))
		results <- workerResult{id: id, err: err}
		fmt.Printf("Worker %d: Task failed.\n", id)
	}

	value := fmt.Sprintf("Data from Worker %d", id)
	results <- workerResult{id: id, value: value}
	fmt.Printf("Worker %d: Task succeeded.\n", id)

}

func main() {
	fmt.Println("Start demonstration of sync.Error.....")

	var wg sync.WaitGroup
	const numWorkers = 5
	resultsChan := make(chan workerResult, numWorkers)

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, resultsChan, &wg)

	}

	go func() {
		wg.Wait()
		close(resultsChan)
		fmt.Println("All working goroutines closed, channel closed successful")

	}()

	var successfulResults []string
	var failedResults []error

	for res := range resultsChan {
		if res.err != nil {
			fmt.Printf("Main: Received error from Worker %d: %v\n",
				res.id, res.err)
			failedResults = append(failedResults, res.err)
		} else {
			fmt.Printf("Main: Received successful result from Worker %d:"+
				" %s\n", res.id, res.value)
			successfulResults = append(successfulResults, res.value)
		}
	}

	fmt.Println("\n--- Results: ---")
	fmt.Printf("Successful completed Results: %d\n", len(successfulResults))
	for _, val := range successfulResults {
		fmt.Printf("  - %s\n", val)
	}
	fmt.Printf("Counts of Errors: %d\n", len(failedResults))
	for _, err := range failedResults {
		fmt.Printf("  - %v\n", err)
	}

	fmt.Println("\nDemonstration of sync.Error completed.")

}
