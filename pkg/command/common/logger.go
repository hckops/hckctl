package common

import (
	"github.com/rs/zerolog/log"
)

func SetupLogger(config *LogConfig) error {
	log.Info().Msgf("config LEVEL=%s", config.Level)
	return nil
}
