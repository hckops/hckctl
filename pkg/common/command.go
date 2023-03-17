package common

import "fmt"

const (
	CommandPrefix    string = "hck-"
	CommandDelimiter string = "::"
)

// TODO https://pkg.go.dev/golang.org/x/tools/cmd/stringer
type Command string

const (
	CommandBoxCreate Command = "hck-box-create"
	CommandBoxOpen           = "hck-box-open"
	CommandBoxTunnel         = "hck-box-tunnel"
)

// TODO example schema "{"kind":"action/v1","name":"hck-box-open","body":{"template":"alpine","revision":"main"}}"
// TODO simple client validation
func NewCommandOpenBox(name, revision string) string {
	return fmt.Sprintf("%s%s%s%s%s", CommandBoxOpen, CommandDelimiter, name, CommandDelimiter, revision)
}
