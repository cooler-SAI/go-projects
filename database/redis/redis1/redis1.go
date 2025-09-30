package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis client instance
var rdb *redis.Client

// initRedis initializes the Redis client connection
func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis address. Make sure Docker is running.
		DB:   0,
	})

	// Test connection
	ctx := context.Background()
	if rdb.Ping(ctx).Err() != nil {
		log.Fatalf("Redis connection error: %v. Please check Docker.", rdb.Ping(ctx).Err())
	}
	fmt.Println("Successfully connected to Redis!")
}

func main() {
	initRedis()
	ctx := context.Background()

	// --- 1. SET (Write to cache / Add to "Bartender's Cheat Sheet") ---
	key := "vip_order:coffee"
	value := "Large Latte, Oat Milk, extra shot"
	expiration := 5 * time.Second // Set short TTL (5 seconds)

	// SET - Write order with time-to-live
	if err := rdb.Set(ctx, key, value, expiration).Err(); err != nil {
		log.Fatalf("Error writing to Redis: %v", err)
	}
	fmt.Printf("\n[1] SET: VIP order written to cache. Time-to-live: %v\n", expiration)

	// --- 2. GET (First read / Lookup in "Cheat Sheet") ---
	// Read value immediately. Should be Cache Hit.
	val1, err := rdb.Get(ctx, key).Result()
	if err == nil {
		fmt.Printf("[2] GET (Cache Hit): VIP order quickly found: '%s'\n", val1)
	} else {
		// Should not happen in this scenario
		fmt.Printf("[2] GET (Cache Miss - Unexpected): Error reading: %v\n", err)
	}

	// --- 3. WAIT (Wait for TTL expiration / VIP customer left) ---
	fmt.Printf("\n[3] Waiting %v for order to be removed from cache...\n", expiration+1*time.Second)
	time.Sleep(expiration + 1*time.Second)

	// --- 4. GET after TTL (Second read / Lookup for expired order) ---
	// Attempt to read the value again. Should result in Cache Miss.
	val2, err := rdb.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		// redis.Nil is a special error indicating key not found (Cache Miss)
		fmt.Println("[4] GET (Cache Miss - Success): Order expired (TTL), key removed from Redis. Received error: redis.Nil")
	} else if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	} else {
		fmt.Printf("[4] GET (Unexpected Hit): Cache not expired (Value: %s). TTL logic error.\n", val2)
	}

	err2 := rdb.Close()
	if err2 != nil {
		return
	}
	fmt.Println("\nRedis demonstration completed.")
	fmt.Println("Stop Redis container: docker stop my-redis")
}
