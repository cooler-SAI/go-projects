package main

import (
	"fmt"
	"runtime"
	"sync"
)

var counter int = 0

func incrementor(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 20; i++ {
		counter++
		fmt.Println(counter)
	}
}

func main() {
	fmt.Println("CPU cores:", runtime.NumCPU())
	fmt.Println("Goroutines:", runtime.NumGoroutine())

	var wg sync.WaitGroup
	const numGoroutines = 1000
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go incrementor(&wg)
	}

	wg.Wait()
	fmt.Println("Final counter value (without mutex):", counter)

	fmt.Println("Expected value: 1000000")

	fmt.Println("Goroutines at end:", runtime.NumGoroutine())

}
