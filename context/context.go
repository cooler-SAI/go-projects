package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func performLongTask(ctx context.Context, taskName string, duration time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("%s: Starting task for %s...\n", taskName, duration)

	select {
	case <-time.After(duration):
		fmt.Printf("%s: Task completed successfully!\n", taskName)
	case <-ctx.Done():
		fmt.Printf("%s: Task canceled! Reason: %v\n", taskName, ctx.Err())
	}
}

func RandRange(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup

	// --- Scenario 1: Task completes before timeout ---
	fmt.Println("\n--- Scenario 1: Task finishes before timeout (2 seconds) ---")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel1()

	wg.Add(1)
	go performLongTask(ctx1, "Task A", 2*time.Second, &wg)
	wg.Wait()

	// --- Scenario 2: Task gets canceled by timeout ---
	fmt.Println("\n--- Scenario 2: Task canceled by timeout (1 second) ---")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel2()

	wg.Add(1)
	go performLongTask(ctx2, "Task B", 3*time.Second, &wg)
	wg.Wait()

	// --- Scenario 3: Random duration task (1-15s) with fixed 10s timeout ---
	fmt.Println("\n--- Scenario 3: Random task duration (1-15s) with 10s timeout ---")
	randomDuration := time.Duration(RandRange(1, 15)) * time.Second
	ctx3, cancel3 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel3()

	wg.Add(1)
	go performLongTask(ctx3, "Task C", randomDuration, &wg)
	wg.Wait()

	// --- Scenario 4: Random timeout (1-15s) for 5s task ---
	fmt.Println("\n--- Scenario 4: Fixed 5s task with random timeout (1-15s) ---")
	randomTimeout := time.Duration(RandRange(1, 15)) * time.Second
	ctx4, cancel4 := context.WithTimeout(context.Background(), randomTimeout)
	defer cancel4()

	wg.Add(1)
	go performLongTask(ctx4, "Task D", 5*time.Second, &wg)
	wg.Wait()

	fmt.Println("\nAll scenarios completed.")
}
