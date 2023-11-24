//go:build local

package main

import "github.com/rs/zerolog/log"

func main() {
	err := Init("123", "456")
	if err != nil {
		log.Fatal().Msg("Failed to save connection details")
	}

	log.Info().Msg("Successfully saved connection details")
}
