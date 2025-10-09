package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type VIPClient struct {
	ID         string
	Name       string
	Level      string
	TotalSpent int
}

type VIPManager struct {
	rdb *redis.Client
	ctx context.Context
}

func NewVIPManager() *VIPManager {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	ctx := context.Background()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("🚨 Redis connection error:", err)
	}

	return &VIPManager{rdb: rdb, ctx: ctx}
}

// AddClient Add client (Hashes + Sets)
func (vm *VIPManager) AddClient(client VIPClient) {
	key := "vip:" + client.ID
	vm.rdb.HSet(vm.ctx, key,
		"name", client.Name,
		"level", client.Level,
		"total_spent", client.TotalSpent,
	)

	vm.rdb.SAdd(vm.ctx, "vip_level:"+client.Level, client.ID)
	fmt.Printf("✅ ADDED: %s (%s)\n", client.Name, client.Level)
}

// GetByLevel Get clients by level (Sets)
func (vm *VIPManager) GetByLevel(level string) {
	clients, _ := vm.rdb.SMembers(vm.ctx, "vip_level:"+level).Result()
	fmt.Printf("👑 LEVEL %s: %v\n", level, clients)
}

// UpdateSpending Update spending and level (Hashes + Sets)
func (vm *VIPManager) UpdateSpending(clientID string, amount int) {
	key := "vip:" + clientID // Fixed: was client.ID, now clientID
	vm.rdb.HIncrBy(vm.ctx, key, "total_spent", int64(amount))

	// Get updated data
	data, _ := vm.rdb.HGetAll(vm.ctx, key).Result()
	currentSpent, _ := strconv.Atoi(data["total_spent"])

	// Determine new level
	var newLevel string
	switch {
	case currentSpent >= 100000:
		newLevel = "VIP3"
	case currentSpent >= 50000:
		newLevel = "VIP2"
	default:
		newLevel = "VIP1"
	}

	// Update if level changed
	if newLevel != data["level"] {
		vm.rdb.SRem(vm.ctx, "vip_level:"+data["level"], clientID)
		vm.rdb.SAdd(vm.ctx, "vip_level:"+newLevel, clientID)
		vm.rdb.HSet(vm.ctx, key, "level", newLevel)
		fmt.Printf("🎉 LEVEL UP! %s → %s\n", data["name"], newLevel)
		fmt.Printf("💰 Total spent: $%d\n", currentSpent)
	}
}

// Stats Statistics (Sets)
func (vm *VIPManager) Stats() {
	levels := []string{"VIP1", "VIP2", "VIP3"}
	fmt.Println("\n📊 VIP STATISTICS:")
	for _, level := range levels {
		count, _ := vm.rdb.SCard(vm.ctx, "vip_level:"+level).Result()
		fmt.Printf("   %s 🎯: %d clients\n", level, count)
	}
}

// GetClientDetails Get client details (Hashes)
func (vm *VIPManager) GetClientDetails(clientID string) {
	key := "vip:" + clientID
	data, err := vm.rdb.HGetAll(vm.ctx, key).Result()
	if err != nil {
		fmt.Printf("❌ Error getting client: %s\n", clientID)
		return
	}

	if len(data) == 0 {
		fmt.Printf("🔍 Client not found: %s\n", clientID)
		return
	}

	fmt.Printf("\n📋 CLIENT CARD:\n")
	fmt.Printf("   🆔 ID: %s\n", clientID)
	fmt.Printf("   👤 Name: %s\n", data["name"])
	fmt.Printf("   🎯 Level: %s\n", data["level"])
	fmt.Printf("   💰 Total spent: $%s\n", data["total_spent"])
}

func main() {
	fmt.Println("🏆 VIP CLIENT MANAGEMENT SYSTEM")
	fmt.Println("🗃️ Hashes - client data | 👥 Sets - level groups\n")

	manager := NewVIPManager()

	// Add clients
	clients := []VIPClient{
		{"vip001", "Anna Petrova", "VIP1", 25000},
		{"vip002", "Boris Ivanov", "VIP2", 75000},
		{"vip003", "Victor Sidorov", "VIP1", 15000},
		{"vip004", "Maria Kozlova", "VIP3", 150000},
	}

	fmt.Println("👥 ADDING CLIENTS:")
	for _, client := range clients {
		manager.AddClient(client)
	}

	fmt.Println("\n--- LEVEL GROUPS ---")
	manager.GetByLevel("VIP1")
	manager.GetByLevel("VIP2")
	manager.GetByLevel("VIP3")

	fmt.Println("\n--- CLIENT DETAILS ---")
	manager.GetClientDetails("vip001")

	fmt.Println("\n--- UPDATING SPENDING ---")
	manager.UpdateSpending("vip001", 80000)

	fmt.Println("\n--- UPDATED LEVEL GROUPS ---")
	manager.GetByLevel("VIP1")
	manager.GetByLevel("VIP2")
	manager.GetByLevel("VIP3")

	fmt.Println("\n--- CLIENT DETAILS AFTER UPDATE ---")
	manager.GetClientDetails("vip001")

	// Final statistics
	manager.Stats()

	fmt.Println("\n🎯 SYSTEM SUMMARY:")
	fmt.Println("   ✅ Hashes store structured client data")
	fmt.Println("   ✅ Sets manage unique level groups")
	fmt.Println("   ✅ Automatic level promotion based on spending")
	fmt.Println("   ✅ Real-time statistics and tracking")
}
