package loader

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

func (l *Loader) Update(message string) {
	l.spinner.Suffix = fmt.Sprintf("  %s", message)
}

func (l *Loader) Start(message string) {
	l.Update(message)
	l.spinner.Start()
}

func (l *Loader) Stop() {
	l.spinner.Stop()
}
