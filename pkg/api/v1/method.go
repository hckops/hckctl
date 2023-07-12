package v1

import (
	"fmt"
)

type MethodName int

const (
	MethodPing MethodName = iota
	MethodBoxCreate
	MethodBoxExec
	MethodBoxDelete
	MethodBoxList
)

var methods = map[MethodName]string{
	MethodPing:      "hck-ping",
	MethodBoxCreate: "hck-box-create",
	MethodBoxExec:   "hck-box-exec",
	MethodBoxDelete: "hck-box-delete",
	MethodBoxList:   "hck-box-list",
}

func (c MethodName) String() string {
	return methods[c]
}

func toMethodName(value string) (MethodName, error) {
	for methodName, methodValue := range methods {
		if methodValue == value {
			return methodName, nil
		}
	}
	return -1, fmt.Errorf("method not found %s", value)
}
