package v1

import (
	"fmt"
)

type MethodName int

const (
	MethodPing MethodName = iota
	MethodBoxCreate
	MethodBoxDelete
	MethodBoxDescribe
	MethodBoxExec
	MethodBoxList
	MethodLabCreate
)

var methods = map[MethodName]string{
	MethodPing:        "hck-ping",
	MethodBoxCreate:   "hck-box-create",
	MethodBoxDelete:   "hck-box-delete",
	MethodBoxDescribe: "hck-box-describe",
	MethodBoxExec:     "hck-box-exec",
	MethodBoxList:     "hck-box-list",
	MethodLabCreate:   "hck-lab-create",
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
