package main

import (
	"github.com/cooler-SAI/go-Tools/zerolog"
)

func main() {

	zerolog.Init()

	zerolog.Log.Info().Msg("Application started")

	zerolog.Log.Warn().
		Str("component", "main").
		Int("count", 42).
		Msg("Warning message")
	zerolog.Log.Printf("Warning!!!")
}
