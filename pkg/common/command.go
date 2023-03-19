package common

import "fmt"

const (
	CommandPrefix    string = "hck-"
	CommandDelimiter string = "::"
)

// TODO https://pkg.go.dev/golang.org/x/tools/cmd/stringer
type Command string

const (
	CommandBoxCreate Command = "hck-box-create" // creates a detached box
	CommandBoxOpen           = "hck-box-open"   // creates and attach+tunnel box, with optional id skip creation
	CommandBoxTunnel         = "hck-box-tunnel" // creates and tunnel box only, with optional id skip creation
)

// TODO change payload format e.g. json, protocol buffer
// TODO schema example "{"kind":"action/v1","name":"hck-box-open","body":{"template":"official/alpine","revision":"main"}}"
// TODO simple client validation
func NewCommandOpenBox(name, revision string) string {
	return fmt.Sprintf("%s%s%s:%s", CommandBoxOpen, CommandDelimiter, name, revision)
}
