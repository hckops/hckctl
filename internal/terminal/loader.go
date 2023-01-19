package terminal

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
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
	l.spinner.Suffix = fmt.Sprintf("  %s", message)
}

func (l *Loader) Refresh(message string) {
	l.update(message)
	l.spinner.Reverse()
	l.spinner.Restart()
}

func (l *Loader) Stop() {
	l.spinner.Stop()
}
