package common

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

func (l *Loader) update(message string, values ...any) {
	msg := fmt.Sprintf(message, values...)
	log.Debug().Msgf("update: %s", msg)
	l.spinner.Suffix = fmt.Sprintf("  %s", msg)
}

func (l *Loader) Start(message string, values ...any) {
	l.update(message, values...)
	l.spinner.Start()
}

func (l *Loader) Reload() {
	l.spinner.Reverse()
	l.spinner.Restart()
}

func (l *Loader) Refresh(message string, values ...any) {
	l.update(message, values...)
	l.Reload()
}

func (l *Loader) Sleep(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

func (l *Loader) Stop() {
	l.spinner.Stop()
}

func (l *Loader) Halt(err error, message string) {
	l.Stop()
	fmt.Println(message)
	log.Fatal().Err(err).Msg(message)
}
