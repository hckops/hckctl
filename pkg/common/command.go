package common

import "fmt"

const (
	CommandRequestType   string = "hck-v1"
	CommandResponseError string = "error"
	CommandDelimiter     string = "::"
)

type Command int

const (
	CommandBoxCreate Command = iota
	CommandBoxExec
	CommandBoxTunnel
	CommandBoxOpen
	CommandBoxList
	CommandBoxDelete
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
	case CommandBoxList: // list existing boxes
		return "hck-box-list"
	case CommandBoxDelete: // delete a box
		return "hck-box-delete"
	default:
		return ""
	}
}

// TODO change payload format e.g. json, protocol buffer
// TODO schema example "{"kind":"action/v1","name":"hck-box-open","body":{"template":"official/alpine","revision":"main","resources":{"memory:"512MB","cpu":"0.5"}}}"
// TODO basic client side validation

func NewCommandCreateBox(name, revision string) string {
	return fmt.Sprintf("%s%s%s%s%s", CommandBoxCreate, CommandDelimiter, name, CommandDelimiter, revision)
}

func NewCommandExecBox(name, revision, boxId string) string {
	return fmt.Sprintf("%s%s%s%s%s%s%s", CommandBoxExec, CommandDelimiter, name, CommandDelimiter, revision, CommandDelimiter, boxId)
}

func NewCommandListBox() string {
	return fmt.Sprintf("%s", CommandBoxList)
}

func NewCommandDeleteBox(name, revision, boxId string) string {
	return fmt.Sprintf("%s%s%s%s%s%s%s", CommandBoxDelete, CommandDelimiter, name, CommandDelimiter, revision, CommandDelimiter, boxId)
}
