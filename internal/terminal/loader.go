package terminal

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/rs/zerolog/log"
)

type Loader struct {
	spinner *spinner.Spinner
}

func NewLoader() *Loader {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Color("green")

	return &Loader{
		spinner: s,
	}
}

func (l *Loader) Start(message string) {
	l.update(message)
	l.spinner.Start()
}

func (l *Loader) update(message string) {
	log.Debug().Msgf("loader: %s", message)
	l.spinner.Suffix = fmt.Sprintf("  %s", message)
}

func (l *Loader) Refresh(message string) {
	l.update(message)
	l.spinner.Reverse()
	l.spinner.Restart()
}

func (l *Loader) Sleep(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

func (l *Loader) Stop() {
	l.spinner.Stop()
}
