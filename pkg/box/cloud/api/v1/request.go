package v1

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/schema"
	"github.com/hckops/hckctl/pkg/util"
)

type requestName int

const (
	requestPing requestName = iota
	requestBoxCreate
	requestBoxDelete
	requestBoxList
)

var requests = map[requestName]string{
	requestPing:      "hck-ping",
	requestBoxCreate: "hck-box-create",
	requestBoxDelete: "hck-box-delete",
	requestBoxList:   "hck-box-list",
}

func (c requestName) String() string {
	return requests[c]
}

type Request[T body] struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
	Body T      `json:"body,omitempty"`
}

type body interface {
	name() requestName
}

func (req Request[T]) Encode() (string, error) {
	return util.EncodeJson(req)
}

func Decode[T body](value string) (*Request[T], error) {
	var request Request[T]
	if err := json.Unmarshal([]byte(value), &request); err != nil {
		return nil, errors.Wrap(err, "error decoding json")
	}
	return &request, nil
}

func newRequest[T body](body T) *Request[T] {
	return &Request[T]{
		Kind: schema.KindApiV1.String(),
		Name: body.name().String(),
		Body: body,
	}
}
