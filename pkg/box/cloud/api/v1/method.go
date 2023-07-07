package v1

type methodName int

const (
	MethodPing methodName = iota
	MethodBoxCreate
	MethodBoxExec
	MethodBoxDelete
	MethodBoxList
)

var methods = map[methodName]string{
	MethodPing:      "hck-ping",
	MethodBoxCreate: "hck-box-create",
	MethodBoxExec:   "hck-box-exec",
	MethodBoxDelete: "hck-box-delete",
	MethodBoxList:   "hck-box-list",
}

func (c methodName) String() string {
	return methods[c]
}
