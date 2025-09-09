package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type SafeMap struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		data: make(map[string]string),
	}
}

func (sm *SafeMap) Set(key, value string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[key] = value
	fmt.Printf("Set: Key '%s' implement into '%s'\n", key, value)
}

func (sm *SafeMap) Get(key string) string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	value, ok := sm.data[key]
	fmt.Printf("Get: Reading key '%s', got '%s' (sucseed: %t)\n", key, value, ok)
	return value
}

func readWorker(id int, sm *SafeMap, wg *sync.WaitGroup) {
	defer wg.Done()
	keys := []string{"keyA", "keyB", "keyC", "keyD"}
	for i := 0; i < 5; i++ {
		key := keys[rand.Intn(len(keys))]
		value := sm.Get(key)
		fmt.Printf("Worker %d: got %s => %v\n", id, key, value)
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}
}

func writeWorker(id int, sm *SafeMap, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 2; i++ {
		key := fmt.Sprintf("key%d", rand.Intn(5))
		value := fmt.Sprintf("value%d-%d", id, i)
		sm.Set(key, value)
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
	}
}

func main() {

	fmt.Println("Start demonstration of RWMutex")
	rand.Seed(time.Now().UnixNano())

	safeMap := NewSafeMap()
	var wg = sync.WaitGroup{}

	safeMap.Set("keyA", "valueA")
	safeMap.Set("keyB", "valueB")
	safeMap.Set("keyC", "valueC")

	numReaders := 10
	for i := 1; i <= numReaders; i++ {
		wg.Add(1)
		go readWorker(i, safeMap, &wg)
	}

	numWriters := 2
	for i := 1; i <= numWriters; i++ {
		wg.Add(1)
		go writeWorker(i, safeMap, &wg)
	}

	wg.Wait()

	fmt.Println("\nDemonstration of sync.RWMutex completed.")
	fmt.Println("Final map state:")
	for key, value := range safeMap.data {
		fmt.Printf("Key: '%s', Value: '%s'\n", key, value)
	}
	fmt.Println(
		"Job is done")

}
