package model

import (
	"context"
	"io"
)

type BoxContext struct {
	Ctx      context.Context
	Template *BoxV1
	Streams  *BoxStreams
}

type BoxStreams struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	IsTty  bool // TODO tty false for tunnel only
}

// TODO ???
type BoxCallback interface {
	OnCreate()
	OnOpen()
	OnClose()
}
