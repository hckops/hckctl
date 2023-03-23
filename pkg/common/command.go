package common

import "fmt"

const (
	CommandPrefix    string = "hck-"
	CommandDelimiter string = "::"
)

type Command int

const (
	CommandBoxCreate Command = iota
	CommandBoxExec
	CommandBoxTunnel
	CommandBoxOpen
)

// TODO https://pkg.go.dev/golang.org/x/tools/cmd/stringer
func (c Command) String() string {
	switch c {
	case CommandBoxCreate: // create a long-running detached box
		return "hck-box-create"
	case CommandBoxExec: // attach to a box
		return "hck-box-exec"
	case CommandBoxTunnel: // tunnel a box
		return "hck-box-tunnel"
	case CommandBoxOpen: // create, attach and tunnel to an ephemeral box
		return "hck-box-open"
	default:
		return ""
	}
}

func NewCommandCreateBox(name, revision string) string {
	return fmt.Sprintf("%s%s%s", name, CommandDelimiter, revision)
}

// TODO change payload format e.g. json, protocol buffer
// TODO schema example "{"kind":"action/v1","name":"hck-box-open","body":{"template":"official/alpine","revision":"main"}}"
// TODO simple client validation
func NewCommandOpenBox(name, revision string) string {
	return fmt.Sprintf("%s%s%s:%s", CommandBoxOpen.String(), CommandDelimiter, name, revision)
}
