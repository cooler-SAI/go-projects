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

	go func() {
		zerolog.Log.Info().Msg("Starting sync.Cond demonstration...")
		cond.L.Lock()
		defer cond.L.Unlock()
		cond.Wait()
		zerolog.Log.Info().Msg("I got a signal! Gratz!")
	}()

	go func() {
		time.Sleep(5 * time.Second)
		zerolog.Log.Info().Msg("Sending a signal...")
		cond.L.Lock()
		defer cond.L.Unlock()
		cond.Signal()
	}()

	var rwMU sync.RWMutex
	condRW := sync.NewCond(rwMU.RLocker())

	go func() {
		rwMU.RLock()
		defer rwMU.RUnlock()
		condRW.Wait()
		zerolog.Log.Info().Msg("I got a RWMutex signal! Gratz!")
	}()

	go func() {
		time.Sleep(5 * time.Second)
		zerolog.Log.Info().Msg("Sending RWMutex signal...")
		rwMU.Lock()
		defer rwMU.Unlock()
		condRW.Signal()
	}()

	time.Sleep(6 * time.Second)
}
