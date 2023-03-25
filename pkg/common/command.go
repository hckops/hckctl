package common

import "fmt"

const (
	CommandRequestType   string = "hck-v1"
	CommandResponseError string = "error"
	CommandDelimiter     string = "::"
)

type Command int

const (
	CommandBoxCreate Command = iota // create a long-running detached box
	CommandBoxExec                  // attach to an existing box
	CommandBoxOpen                  // create and attach to a box
	CommandBoxTunnel                // tunnel a box
	CommandBoxList                  // list existing boxes
	CommandBoxDelete                // delete a box
)

func Values() []string {
	return []string{
		"hck-box-create",
		"hck-box-exec",
		"hck-box-open",
		"hck-box-tunnel",
		"hck-box-list",
		"hck-box-delete",
	}
}

// TODO https://pkg.go.dev/golang.org/x/tools/cmd/stringer
// Command is automatically converted in fmt.Sprintf with "%s" String
func (c Command) String() string {
	return Values()[c]
}

func FromString(str string) (Command, error) {
	for index, value := range Values() {
		if str == value {
			return Command(index), nil
		}
	}
	return -1, fmt.Errorf("invalid command: %s", str)
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

func NewCommandOpenBox(name, revision string) string {
	return fmt.Sprintf("%s%s%s%s%s", CommandBoxOpen, CommandDelimiter, name, CommandDelimiter, revision)
}

func NewCommandListBox() string {
	return fmt.Sprintf("%s", CommandBoxList)
}

func NewCommandDeleteBox(name, revision, boxId string) string {
	return fmt.Sprintf("%s%s%s%s%s%s%s", CommandBoxDelete, CommandDelimiter, name, CommandDelimiter, revision, CommandDelimiter, boxId)
}
