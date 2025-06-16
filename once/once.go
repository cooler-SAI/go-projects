package main

import (
	"github.com/cooler-SAI/go-Tools/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

var config string
var once sync.Once

func initializeConfig() {
	log.Info().Msg("Initializing configuration...")
	time.Sleep(500 * time.Millisecond)
	config = "Application Configuration Loaded"
	log.Info().Msg("Configuration initialized!")
}

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Info().
		Int("worker_id", id).
		Msg("Worker attempting to load config...")

	once.Do(initializeConfig)

	log.Info().
		Int("worker_id", id).
		Str("config", config).
		Msg("Worker accessed config")
}

func main() {

	zerolog.ConfigureZerologConsoleWriter()

	log.Info().Msg("Starting sync.Once demonstration...")

	var wg sync.WaitGroup
	const numWorkers = 5

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	wg.Wait()

	log.Info().
		Int("workers_count", numWorkers).
		Msg("All workers completed")

	log.Info().Msg("Demonstration finished")
}
