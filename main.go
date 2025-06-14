package main

import (
	"fmt"
	"sync"
)

var counter int
var mu sync.Mutex
var wg sync.WaitGroup

func worker(id int) {
	defer wg.Done()
	id = 20
	mu.Lock()
	counter++
	mu.Unlock()
}

func main() {

	fmt.Println("hi all here!")
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go worker(i)
	}
	wg.Wait()
	fmt.Println("counter:", counter)

}
