package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	coffeeCups []string
	door       = sync.Mutex{}
	bell       = sync.NewCond(&door)
)

func barista(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("🐌 Barista: Prepare to make coffee...")
	time.Sleep(2 * time.Second)

	door.Lock()
	coffeeCups = append(coffeeCups, "☕ Espresso")
	coffeeCups = append(coffeeCups, "☕ Cappuccino")
	coffeeCups = append(coffeeCups, "☕ Latte")
	fmt.Println("🍵 Barista: Coffee is ready!")
	door.Unlock()

	bell.Broadcast()

}

func programmer(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	door.Lock()
	defer door.Unlock()

	for len(coffeeCups) == 0 {
		fmt.Printf("💻 Programmer %d: No coffee available. Waiting...\n", id)
		bell.Wait()
	}
	cup := coffeeCups[0]
	coffeeCups = coffeeCups[1:]
	fmt.Printf("💻 Programmer %d: Enjoying my %s\n", id, cup)
}

func main() {
	var wg sync.WaitGroup
	numProgrammers := 3

	wg.Add(numProgrammers)
	for i := 0; i < numProgrammers; i++ {
		go programmer(i, &wg)
	}

	wg.Add(1)
	go barista(&wg)

	wg.Wait()
	fmt.Println("\nAll programmers have their coffee and are coding happily.")

}
