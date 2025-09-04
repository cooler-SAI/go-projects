package main

import (
	"fmt"
	"sync"
	"time"
)

var databaseConnection string
var once sync.Once

func initializeDatabase() {
	fmt.Println("Initializing database connection...")
	time.Sleep(2 * time.Second)
	databaseConnection = "Database connected successfully!"
	fmt.Println("Database initialization completed.")
}

func GetDatabaseConnection() {
	once.Do(initializeDatabase)
	fmt.Printf("Connection accessed. Status: %s\n", databaseConnection)
}

func main() {
	fmt.Println("Starting demonstration of sync.Once...")
	var wg sync.WaitGroup

	numGoroutines := 5
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Goroutine %d trying to get connection...\n", id)
			GetDatabaseConnection()
		}(i)
	}

	wg.Wait()
	fmt.Println("\nAll goroutines have finished.")
}
