package main

import (
	"sync"

	"github.com/cooler-SAI/go-Tools/zerolog"
)

func producerGoroutine(outChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	zerolog.Log.Info().Msg("Producer goroutine started")
	for i := 0; i < 10; i++ {
		zerolog.Log.Info().Int("number", i).Msg("Producer goroutine sending number")
		outChan <- i
	}
	close(outChan)
}

func consumerGoroutine(inChan <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	zerolog.Log.Info().Msg("Consumer goroutine started")
	for num := range inChan {
		zerolog.Log.Info().Int("number", num).Msg("Consumer goroutine received number")
	}
	zerolog.Log.Info().Msg("Consumer goroutine stopped")
}

func main() {
	var wg sync.WaitGroup
	zerolog.Init()

	zerolog.Log.Info().Msg("Starting channel demonstration...")

	dataChan := make(chan int, 5)

	wg.Add(1)
	go producerGoroutine(dataChan, &wg)

	wg.Add(1)
	go consumerGoroutine(dataChan, &wg)

	wg.Wait()
	zerolog.Log.Info().Msg("All goroutines finished. Exiting main.")
}
