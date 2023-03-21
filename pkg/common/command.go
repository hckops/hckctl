package common

import "fmt"

const (
	CommandPrefix    string = "hck-"
	CommandDelimiter string = "::"
)

// TODO https://pkg.go.dev/golang.org/x/tools/cmd/stringer
type Command string

const (
	CommandBoxCreate Command = "hck-box-create" // create a long-running detached box
	CommandBoxExec           = "hck-box-exec"   // attach to a box
	CommandBoxTunnel         = "hck-box-tunnel" // tunnel a box
	CommandBoxOpen           = "hck-box-open"   // create, attach and tunnel to an ephemeral box
)

func NewCommandCreateBox(name, revision string) string {
	return fmt.Sprintf("%s%s%s", name, CommandDelimiter, revision)
}

// TODO change payload format e.g. json, protocol buffer
// TODO schema example "{"kind":"action/v1","name":"hck-box-open","body":{"template":"official/alpine","revision":"main"}}"
// TODO simple client validation
func NewCommandOpenBox(name, revision string) string {
	return fmt.Sprintf("%s%s%s:%s", CommandBoxOpen, CommandDelimiter, name, revision)
}
