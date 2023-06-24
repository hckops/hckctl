package command

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/schema"
	"github.com/hckops/hckctl/pkg/util"
)

type commandName int

const (
	commandPing commandName = iota
	commandBoxCreate
	commandBoxDelete
	commandBoxList
)

var commands = map[commandName]string{
	commandPing:      "hck-ping",
	commandBoxCreate: "hck-box-create",
	commandBoxDelete: "hck-box-delete",
	commandBoxList:   "hck-box-list",
}

func (c commandName) String() string {
	return commands[c]
}

type Request[T body] struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
	Body T      `json:"body,omitempty"`
}

type body interface {
	cmdName() commandName
}

func (c Request[T]) Encode() (string, error) {
	return util.EncodeJson(c)
}

func Decode[T body](value string) (*Request[T], error) {
	var remoteCommand Request[T]
	if err := json.Unmarshal([]byte(value), &remoteCommand); err != nil {
		return nil, errors.Wrap(err, "error decoding json")
	}
	return &remoteCommand, nil
}

func newRequest[T body](body T) *Request[T] {
	return &Request[T]{
		Kind: schema.KindCommandV1.String(),
		Name: body.cmdName().String(),
		Body: body,
	}
}
