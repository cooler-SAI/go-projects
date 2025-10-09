package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	fmt.Println("🚀 Redis Queue Demo")

	// Initialize Redis client
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("Warning: Error closing Redis: %v\n", err)
		}
	}(client)

	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		fmt.Printf("❌ Redis connection failed: %v\n", err)
		return
	}
	fmt.Println("✅ Connected to Redis!")

	queueName := "task_queue"

	// Producer - add tasks to queue
	fmt.Println("\n📤 Adding tasks to queue...")
	tasks := []string{
		"Send welcome email",
		"Process payment",
		"Generate report",
		"Backup data",
	}

	for i, task := range tasks {
		client.RPush(ctx, queueName, i, task)
		fmt.Printf("✅ Added task %d: %s\n", i+1, task)
		time.Sleep(300 * time.Millisecond)
	}

	// Consumer - process tasks from queue
	fmt.Println("\n📥 Processing tasks from queue...")

	for i := 0; i < len(tasks); i++ {
		// Blocking pop - waits for task
		result, err := client.BLPop(ctx, 10*time.Second, queueName).Result()
		if err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
			break
		}

		task := result[1]
		fmt.Printf("🎯 Processing: %s\n", task)

		// Simulate work
		time.Sleep(1 * time.Second)
		fmt.Printf("✅ Completed: %s\n", task)
	}

	fmt.Println("\n🎊 All tasks processed!")
}
