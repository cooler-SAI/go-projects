package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Create Redis client - this is our "radio station"
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	ctx := context.Background()

	// Check Redis connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("❌ Cannot connect to Redis: %v\n", err)
		return
	}

	fmt.Println("✅ Connected to Redis!")

	var mu sync.Mutex // Mutex for output synchronization

	// Subscriber 1 - "Radio-1" (first listener)
	go func() {
		pubsub := client.Subscribe(ctx, "notifications:alerts")
		defer func(pubsub *redis.PubSub) {
			err := pubsub.Close()
			if err != nil {
				fmt.Printf("❌ Radio-1 pubsub close error: %v\n", err)
			}
		}(pubsub)

		fmt.Println("📻 Radio-1 is listening...")
		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				fmt.Printf("❌ Radio-1 error: %v\n", err)
				return
			}
			mu.Lock()
			fmt.Printf("📻 Radio-1 received: %s\n", msg.Payload)
			mu.Unlock()
		}
	}()

	// Subscriber 2 - "Radio-2" (second listener)
	go func() {
		pubsub := client.Subscribe(ctx, "notifications:alerts")
		defer func(pubsub *redis.PubSub) {
			err := pubsub.Close()
			if err != nil {
				fmt.Printf("❌ Radio-2 pubsub close error: %v\n", err)
			}
		}(pubsub)

		fmt.Println("📻 Radio-2 is listening...")
		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				fmt.Printf("❌ Radio-2 error: %v\n", err)
				return
			}
			mu.Lock()
			fmt.Printf("📻 Radio-2 received: %s\n", msg.Payload)
			mu.Unlock()
		}
	}()

	// Give subscribers time to connect
	time.Sleep(1 * time.Second)

	// Publisher - "DJ" (sends messages)
	go func() {
		time.Sleep(500 * time.Millisecond) // Wait for subscribers to be ready

		messages := []string{
			"Server CPU usage is high",
			"New user registration: john_doe",
			"Database backup completed successfully",
		}

		fmt.Println("\n🚀 Starting broadcast...")
		for _, msg := range messages {
			// First print what we're sending
			mu.Lock()
			fmt.Printf("🎤 DJ broadcasting: %s\n", msg)
			mu.Unlock()

			// Then send the message
			err := client.Publish(ctx, "notifications:alerts", msg).Err()
			if err != nil {
				fmt.Printf("❌ DJ broadcast error: %v\n", err)
				return
			}

			time.Sleep(1 * time.Second) // Give subscribers time to receive and print
		}
		fmt.Println("✅ Broadcast completed!")
	}()

	// Graceful shutdown with Ctrl+C
	fmt.Println("\n⏳ Program is running... Press Ctrl+C to exit")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n👋 Shutting down...")
	err2 := client.Close()
	if err2 != nil {
		fmt.Printf("❌ Error closing Redis client: %v\n", err)
		return
	}
	fmt.Println("✅ Redis client closed")
	fmt.Println("👋 Goodbye!")
}
