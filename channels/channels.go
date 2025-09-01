package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/cooler-SAI/go-Tools/zerolog"
)

func producerGoroutine(outChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Producer goroutine started")
	for i := 0; i < 10; i++ {
		fmt.Printf("Producer goroutine sending number %d\n", i)
		outChan <- i
	}
	close(outChan)
}

func consumerGoroutine(inChan <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Consumer goroutine started")
	for num := range inChan {
		fmt.Printf("Consumer goroutine received number %d\n", num)
	}
	fmt.Println("Consumer goroutine stopped")
}

func main() {

	var wg sync.WaitGroup
	zerolog.Init()
	zerolog.Log.Info().Msg("Starting channel demonstration...")

	time.Sleep(1 * time.Second)
	fmt.Println("Starting channel demonstration...")

	dataChan := make(chan int, 5)

	wg.Add(1)
	go producerGoroutine(dataChan, &wg)

	wg.Add(1)
	go consumerGoroutine(dataChan, &wg)

	wg.Wait()
	fmt.Println("All goroutines finished. Exiting main.")

}
