package main

import (
	"fmt"
	"sync"
	"time"
)

var config string

var once sync.Once

func initializeConfig() {
	fmt.Println("Initializing configuration...")
	time.Sleep(500 * time.Millisecond)
	config = "Application Configuration Loaded"
	fmt.Println("Configuration initialized!")
}

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d: Attempting to load config...\n", id)
	once.Do(initializeConfig)
	fmt.Printf("Worker %d: Config: %s\n", id, config)
}

func main() {
	fmt.Println("Starting sync.Once demonstration...")
	var wg sync.WaitGroup

	const numWorkers = 5

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	wg.Wait()
	fmt.Println("All workers completed.")
	fmt.Println("Demonstration finished.")

}
