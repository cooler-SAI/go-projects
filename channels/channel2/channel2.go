package main

import (
	"fmt"
	"time"
)

type Task struct {
	ID        int
	Name      string
	Completed bool
}

func main() {

	taskChan := make(chan Task, 2)

	// producer
	go func() {
		taskChan <- Task{ID: 1, Name: "Download File", Completed: false}
		taskChan <- Task{ID: 2, Name: "Research data", Completed: false}
		close(taskChan)
	}()

	// consumer
	go func() {
		for task := range taskChan {
			fmt.Printf("Got task: %s (ID: %d)\n", task.Name, task.ID)
			// work similar to time-consuming task
			time.Sleep(500 * time.Millisecond)
			task.Completed = true
			fmt.Printf("Task '%s' completed!\n", task.Name)
		}
	}()

	// give some time for goroutines to finish
	time.Sleep(2 * time.Second)
}
