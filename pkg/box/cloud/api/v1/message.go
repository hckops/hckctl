package v1

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/schema"
	"github.com/hckops/hckctl/pkg/util"
)

type Message[T body] struct {
	Kind   string `json:"kind"`
	Origin string `json:"origin"`
	Method string `json:"method"`
	Body   T      `json:"body"` // TODO omitempty to remove "body":{}
}

type body interface {
	method() methodName
}

func (req *Message[T]) Protocol() string {
	return fmt.Sprintf("%s/%s", req.Kind, req.Method)
}

func (req *Message[T]) Encode() (string, error) {
	return util.EncodeJson(req)
}

func Decode[T body](value string) (*Message[T], error) {
	var request Message[T]
	if err := json.Unmarshal([]byte(value), &request); err != nil {
		return nil, errors.Wrap(err, "error decoding json")
	}
	return &request, nil
}

func newMessage[T body](origin string, body T) *Message[T] {
	return &Message[T]{
		Kind:   schema.KindApiV1.String(),
		Origin: origin,
		Method: body.method().String(),
		Body:   body,
	}
}

func IsValidProtocol(value string) (string, error) {
	schemaPrefix := fmt.Sprintf("%s/", schema.KindApiV1.String())

	if !strings.HasPrefix(value, schemaPrefix) {
		return "", errors.New("invalid protocol")
	}

	methodValue := strings.ReplaceAll(value, schemaPrefix, "")
	methodName, err := toMethodName(methodValue)
	if err != nil {
		return "", errors.Wrap(err, "invalid method")
	}

	return methodName.String(), nil
}
