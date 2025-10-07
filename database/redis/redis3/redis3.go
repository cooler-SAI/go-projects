package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func readKey(client *redis.Client, ctx context.Context, key string) {
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		fmt.Printf("Failed to get key '%s': %v\n", key, err)
	} else {
		fmt.Printf("Retrieved value from '%s': %s\n", key, val)
	}
}

func main() {
	fmt.Println("Testing... Redis")

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
		fmt.Printf("Redis connection failed: %v\n", err)
		return
	}
	fmt.Println("Successfully connected to Redis!")

	// Key without TTL
	err := client.Set(ctx, "example_key", "Hello, Redis!", 0).Err()
	if err != nil {
		fmt.Printf("Failed to set key: %v\n", err)
		return
	}

	// Key with 5 seconds TTL
	err = client.Set(ctx, "example_key2", "Hello, Redis! 5 Seconds for Save", 5*time.Second).Err()
	if err != nil {
		fmt.Printf("Failed to set key: %v\n", err)
		return
	}

	fmt.Println("Both keys set successfully.")

	// Read function to check key existence
	fmt.Println("\n--- Immediate read ---")
	readKey(client, ctx, "example_key")
	readKey(client, ctx, "example_key2")

	// Wait for 6 seconds to let the TTL key expire
	fmt.Println("\n--- Waiting 6 seconds ---")
	time.Sleep(6 * time.Second)

	// Reading after 6 seconds
	fmt.Println("--- Read after 6 seconds ---")
	readKey(client, ctx, "example_key")
	readKey(client, ctx, "example_key2")
}
