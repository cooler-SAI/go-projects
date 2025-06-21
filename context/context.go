package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func performLongTask(ctx context.Context, taskName string, duration time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("%s: Starting task for %s...\n", taskName, duration)

	select {
	case <-time.After(duration): // Wait for simulated task completion
		fmt.Printf("%s: Task completed successfully!\n", taskName)

	case <-ctx.Done(): // Listen for context cancellation
		// If ctx.Done() channel is closed, it means context was canceled.
		// ctx.Err() returns cancellation reason (e.g., context.DeadlineExceeded).

		fmt.Printf("%s: Task canceled! Reason: %v\n", taskName, ctx.Err())
	}

}

func main() {
	fmt.Println("Demonstrating context.Context with timeout...")

	var wg sync.WaitGroup

	// --- Scenario 1: Task completes before timeout ---
	fmt.Println("\n--- Scenario 1: Task finishes before timeout (2 seconds) ---")
	// Create context with 3-second timeout
	ctx1, cancel1 := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel1() // Important: always call cancel function to release resources,
	// even if context expires by itself.

	wg.Add(1)
	go performLongTask(ctx1, "Task A", 2*time.Second, &wg) // Task takes 2 seconds

	wg.Wait() // Wait for Task A completion

	// --- Scenario 2: Task gets canceled by timeout ---
	fmt.Println("\n--- Scenario 2: Task canceled by timeout (1 second) ---")
	// Create context with 1-second timeout
	ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel2() // Important: always call cancel function.

	wg.Add(1)
	go performLongTask(ctx2, "Task B", 3*time.Second, &wg) // Task takes 3 seconds

	wg.Wait() // Wait for Task B completion (or cancellation)

	fmt.Println("\ncontext.Context demonstration completed.")

	ctx3, cancel3 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel3()

	wg.Add(1)
	go performLongTask(ctx3, "Task C", 3*time.Second, &wg)

	wg.Wait()

	fmt.Println("\ncontext.Context demonstration completed.")

}
