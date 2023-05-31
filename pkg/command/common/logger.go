package common

import (
	"github.com/rs/zerolog/log"
)

const (
	LogLevelFlag = "log-level"
)

func InitFileLogger(global *GlobalCmdOptions) error {
	log.Info().Msgf("LEVEL=%s", global.logLevel)
	return nil
}
