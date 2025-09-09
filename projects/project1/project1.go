package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	ID       int
	Duration time.Duration
	willFail bool
}

type TaskResult struct {
	TaskID int
	Result string
	Err    error
}

func processTask(ctx context.Context, task Task, results chan<- TaskResult, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Task %d: started (%s)\n", task.ID, task.Duration)

	select {
	case <-time.After(task.Duration):
		if task.willFail {
			err := fmt.Errorf("task %d failed", task.ID)
			results <- TaskResult{TaskID: task.ID, Err: err}
			fmt.Printf("Task %d: failed\n", task.ID)
		} else {
			results <- TaskResult{TaskID: task.ID, Result: fmt.Sprintf("task %d success", task.ID)}
			fmt.Printf("Task %d: completed\n", task.ID)
		}
	case <-ctx.Done():
		results <- TaskResult{TaskID: task.ID, Err: ctx.Err()}
		fmt.Printf("Task %d: cancelled (%v)\n", task.ID, ctx.Err())
	}
}

func main() {
	fmt.Println("Starting task processor...")
	rand.Seed(time.Now().UnixNano())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	const numTasks = 5
	results := make(chan TaskResult, numTasks)
	var wg sync.WaitGroup

	// Create and start tasks
	for i := 1; i <= numTasks; i++ {
		wg.Add(1)
		task := Task{
			ID:       i,
			Duration: time.Duration(rand.Intn(4)+1) * time.Second,
			willFail: rand.Intn(5) == 0, // 20% chance to fail
		}
		go processTask(ctx, task, results, &wg)
	}

	// Close results channel when all tasks done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var success, failed, cancelled int
	for res := range results {
		if res.Err != nil {
			if errors.Is(res.Err, context.DeadlineExceeded) {
				cancelled++
			} else {
				failed++
			}
		} else {
			success++
		}
	}

	// Print summary
	fmt.Println("\n--- Results ---")
	fmt.Printf("Successful: %d\n", success)
	fmt.Printf("Failed: %d\n", failed)
	fmt.Printf("Cancelled: %d\n", cancelled)
}
