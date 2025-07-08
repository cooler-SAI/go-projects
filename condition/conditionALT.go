package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("Start demonstration of sync.Cond (Simple version)...")

	var wg sync.WaitGroup

	// sync.Cond with Mutex
	fmt.Println("\n--- Demonstration Cond with Mutex ---")

	mu := sync.Mutex{}
	cond := sync.NewCond(&mu)

	// goroutine Waiter
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Waiting goroutine (Mutex): Waiting signal...")
		cond.L.Lock()
		cond.Wait()
		cond.L.Unlock()
		fmt.Println("Waiting goroutine (Mutex): Received signal! Yay!")
	}()
	// goroutine Sender
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(2 * time.Second)
		fmt.Println("Signaling goroutine (Mutex): Sending signal...")
		cond.L.Lock()
		cond.Signal()
		cond.L.Unlock()
	}()

	// --- Demonstration of sync.Cond with sync.RWMutex ---
	// WARNING: Cond must always be associated with the base RWMutex,
	// not its RLocker(). cond.Wait() requires Lock/Unlock semantics,
	// which RLocker() doesn't provide.
	fmt.Println("\n--- Cond with RWMutex Demonstration ---")
	rwMu := sync.RWMutex{}
	condRW := sync.NewCond(&rwMu)

	// goroutine Waiter
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Waiting goroutine (RWMutex): Waiting for signal...")
		rwMu.Lock()
		defer rwMu.Unlock()
		condRW.Wait()
		fmt.Println("Waiting goroutine (RWMutex): Received RWMutex signal! Success!")
	}()

	// goroutine Sender
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(4 * time.Second)
		fmt.Println("Signaling goroutine (RWMutex): Sending RWMutex " +
			"notification...")
		rwMu.Lock()
		condRW.Signal()
		rwMu.Unlock()
	}()

	wg.Wait()
	fmt.Println("\nDemonstration of sync.Cond completed.")

}
