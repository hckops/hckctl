package common

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/hckops/hckctl/pkg/event"
)

func EventCallback(loader *Loader) func(e event.Event) {
	return func(e event.Event) {
		switch e.Kind() {
		case event.PrintConsole:
			loader.Refresh("loading")
			fmt.Println(e.String())
		case event.LoaderUpdate:
			loader.Refresh(e.String())
		case event.LoaderStop:
			loader.Stop()
		case event.LogInfo:
			log.Info().Msgf("[%v] %s", e.Source(), e.String())
		case event.LogWarning:
			log.Warn().Msgf("[%v] %s", e.Source(), e.String())
		case event.LogError:
			log.Error().Msgf("[%v] %s", e.Source(), e.String())
		default:
			log.Debug().Msgf("[%v] %s", e.Source(), e.String())
		}
	}
}
