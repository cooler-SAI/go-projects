package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func initRedis() error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // если есть пароль
		DB:       0,
		PoolSize: 10, // размер пула соединений
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return redisClient.Ping(ctx).Err()
}

func main() {
	fmt.Println("🚀 Starting Redis Demo...")

	if err := initRedis(); err != nil {
		log.Fatalf("❌ Redis connection failed: %v", err)
	}
	defer func(redisClient *redis.Client) {
		err := redisClient.Close()
		if err != nil {
			fmt.Printf("Warning: Error closing Redis: %v\n", err)
		}
	}(redisClient)

	fmt.Println("✅ Successfully connected to Redis!")

	ctx := context.Background()

	// Test connection
	err := redisClient.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("❌ Redis connection failed: %v", err)
	}
	fmt.Println("✅ Successfully connected to Redis!")

	// 🔥 STEP 1: WRITE data to Redis
	key := "vip_order:coffee"
	value := "Large Latte, Oat Milk, extra shot"
	expiration := 5 * time.Second // Data will auto-delete after 5 seconds

	err = redisClient.Set(ctx, key, value, expiration).Err()
	if err != nil {
		log.Fatalf("❌ Failed to write to Redis: %v", err)
	}
	fmt.Printf("📝 SET: Key '%s' = '%s' (TTL: %v)\n", key, value, expiration)

	// 🔥 STEP 2: READ data immediately (Cache HIT)
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		log.Fatalf("❌ Failed to read from Redis: %v", err)
	}
	fmt.Printf("✅ GET (Cache Hit): Key '%s' = '%s'\n", key, val)

	// 🔥 STEP 3: Wait for TTL expiration
	fmt.Printf("\n⏰ Waiting %v for TTL expiration...\n", expiration+1*time.Second)
	time.Sleep(expiration + 1*time.Second)

	// 🔥 STEP 4: Try to read again (Cache MISS)
	val, err = redisClient.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		fmt.Printf("❌ GET (Cache Miss): Key '%s' expired - returned redis.Nil error\n", key)
	} else if err != nil {
		log.Fatalf("❌ Unexpected error: %v", err)
	} else {
		fmt.Printf("⚠️ Key still exists: '%s' = '%s'\n", key, val)
	}

	// 🔥 STEP 5: Clean up - close connection
	err = redisClient.Close()
	if err != nil {
		fmt.Printf("Warning: Error closing Redis: %v\n", err)
	}

	fmt.Println("\n🎉 Redis demonstration completed successfully!")
}
