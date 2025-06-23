package main

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/cooler-SAI/go-Tools/random"
	"github.com/rs/zerolog"
)

func performLongTask(ctx context.Context, taskName string, duration time.Duration, wg *sync.WaitGroup, logger zerolog.Logger) {
	defer wg.Done()

	taskLogger := logger.With().
		Str("task", taskName).
		Str("duration", duration.String()).
		Logger()

	taskLogger.Info().Msg("Starting task")

	select {
	case <-time.After(duration):
		taskLogger.Info().Msg("Task completed successfully")
	case <-ctx.Done():
		taskLogger.Warn().
			Err(ctx.Err()).
			Msg("Task canceled")
	}
}

func main() {
	// Configure console logger
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().
		Timestamp().
		Logger()

	var wg sync.WaitGroup

	logger.Info().Msg("Starting scenarios")

	// Scenario 1
	scenarioLog := logger.With().Str("scenario", "1").Logger()
	scenarioLog.Info().Msg("Normal completion (2s task, 3s timeout)")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel1()
	wg.Add(1)
	go performLongTask(ctx1, "TaskA", 2*time.Second, &wg, scenarioLog)

	// Scenario 2
	scenarioLog = logger.With().Str("scenario", "2").Logger()
	scenarioLog.Info().Msg("Timeout case (3s task, 1s timeout)")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel2()
	wg.Add(1)
	go performLongTask(ctx2, "TaskB", 3*time.Second, &wg, scenarioLog)

	// Scenario 3 - Random duration
	scenarioLog = logger.With().Str("scenario", "3").Logger()
	scenarioLog.Info().Msg("Random duration (1-15s)")
	randomDuration := time.Duration(random.RandRange(1, 15)) * time.Second
	ctx3, cancel3 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel3()
	wg.Add(1)
	go performLongTask(ctx3, "TaskC", randomDuration, &wg, scenarioLog)

	// Scenario 4 - Random timeout
	scenarioLog = logger.With().Str("scenario", "4").Logger()
	scenarioLog.Info().Msg("Random timeout (1-15s)")
	randomTimeout := time.Duration(random.RandRange(1, 15)) * time.Second
	ctx4, cancel4 := context.WithTimeout(context.Background(), randomTimeout)
	defer cancel4()
	wg.Add(1)
	go performLongTask(ctx4, "TaskD", 5*time.Second, &wg, scenarioLog)

	wg.Wait()
	logger.Info().Msg("All scenarios completed")
}
