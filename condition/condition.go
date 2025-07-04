package main

import (
	"github.com/cooler-SAI/go-Tools/zerolog"
	"sync"
	"time"
)

func main() {
	var mu sync.Mutex
	cond := sync.NewCond(&mu)

	zerolog.Init()

	// Горутина, которая ждёт сигнала
	go func() {
		zerolog.Log.Info().Msg("Starting sync.Once demonstration...")
		cond.L.Lock()
		cond.Wait()
		cond.L.Unlock()
		zerolog.Log.Info().Msg("I got a signal! Gratz!")
	}()

	go func() {
		time.Sleep(5 * time.Second)
		zerolog.Log.Info().Msg("Sending a signal...")
		cond.L.Lock()
		cond.Signal()
		cond.L.Unlock()
	}()

	time.Sleep(6 * time.Second)
}
