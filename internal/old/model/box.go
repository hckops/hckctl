package model

import (
	"io"
	"os"
)

type BoxStreams struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	IsTty  bool // tty false for tunnel only
}

func NewDefaultStreams(tty bool) *BoxStreams {
	return &BoxStreams{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		IsTty:  tty,
	}
}
