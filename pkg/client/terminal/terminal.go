package terminal

import (
	"fmt"
	"io"

	"github.com/moby/term"
	"github.com/pkg/errors"
)

type RawTerminal struct {
	fileDescriptor uintptr
	previousState  *term.State
}

func NewRawTerminal(in io.Reader) (*RawTerminal, error) {

	if fd, isTerminal := term.GetFdInfo(in); isTerminal {
		previousState, err := term.SetRawTerminal(fd)
		if err != nil {
			return nil, errors.Wrap(err, "error raw terminal")
		}

		return &RawTerminal{
			fileDescriptor: fd,
			previousState:  previousState,
		}, nil
	}

	return nil, fmt.Errorf("error invalid terminal")
}

func (t *RawTerminal) Restore() error {
	return term.RestoreTerminal(t.fileDescriptor, t.previousState)
}
