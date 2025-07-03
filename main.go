package main

import (
	"github.com/cooler-SAI/go-Tools/zerolog"
)

func main() {

	zerolog.Init()

	zerolog.Log.Error().Msg("Hello World")
	zerolog.Log.Info().Msg("Testing")
}
