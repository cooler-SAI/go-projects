package main

import (
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var counter int

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
}

func mainDemonstrationWithoutMutex() {
	var wg sync.WaitGroup
	numGoroutines := 1000
	localCounter := 0

	log.Info().
		Int("num_goroutines", numGoroutines).
		Msg("Starting race condition demonstration WITHOUT mutex")

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			localCounter++
		}()
	}
	wg.Wait()

	log.Warn().
		Int("final_counter_value", localCounter).
		Int("expected_value", 1000).
		Msg("Race condition result: final value does NOT match expected")
}

func mainDemonstrationWithMutex() {
	var mu sync.Mutex
	var wg sync.WaitGroup
	counter = 0
	numGoroutines := 1000

	log.Info().
		Int("num_goroutines", numGoroutines).
		Msg("Starting demonstration WITH mutex")

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()

	log.Info().
		Int("final_counter_value", counter).
		Int("expected_value", 1000).
		Msg("Mutex demonstration completed successfully")
}

func main() {
	log.Info().Msg("Program started: Mutex Demonstration")

	mainDemonstrationWithoutMutex()
	time.Sleep(1 * time.Second)

	mainDemonstrationWithMutex()

	log.Info().Msg("Program completed")
}
