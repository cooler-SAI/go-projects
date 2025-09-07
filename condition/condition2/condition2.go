package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	items = make([]int, 0)
	mtx   sync.Mutex
	cond  = sync.NewCond(&mtx)
)

func producer(wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(2 * time.Second)

	cond.L.Lock()
	for i := 0; i < 3; i++ {
		items = append(items, i+1)
	}
	fmt.Println("Producer: Added 3 items to the queue.")
	cond.L.Unlock()

	cond.Broadcast()
	fmt.Println("Producer: Signaled all waiting consumers.")
}

func consumer(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	cond.L.Lock()
	for len(items) == 0 {
		fmt.Printf("Consumer %d: No items to process. Waiting...\n", id)
		cond.Wait()
	}

	item := items[0]
	items = items[1:]
	fmt.Printf("Consumer %d: Processed item %d\n", id, item)
	cond.L.Unlock()
}

func main() {
	var wg sync.WaitGroup
	numConsumers := 3

	wg.Add(numConsumers)
	for i := 0; i < numConsumers; i++ {
		go consumer(i, &wg)
	}

	wg.Add(1)
	go producer(&wg)

	wg.Wait()
	fmt.Println("\nAll goroutines have finished.")
}
