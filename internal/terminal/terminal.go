package terminal

import (
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

type RawTerminal struct {
	fileDescriptor int
	previousState  *terminal.State
}

func NewRawTerminal() *RawTerminal {
	fd := int(os.Stdin.Fd())

	if terminal.IsTerminal(fd) {
		previousState, err := terminal.MakeRaw(fd)
		if err != nil {
			return nil
		}
		// must be invoked from the caller, use Restore()
		//defer terminal.Restore(fd, previousState)

		return &RawTerminal{
			fileDescriptor: fd,
			previousState:  previousState,
		}
	}
	return nil
}

func (t *RawTerminal) Restore() {
	terminal.Restore(t.fileDescriptor, t.previousState)
}
