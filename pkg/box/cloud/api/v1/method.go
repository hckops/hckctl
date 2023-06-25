package v1

type methodName int

const (
	methodPing methodName = iota
	methodBoxCreate
	methodBoxExec
	methodBoxDelete
	methodBoxList
)

var methods = map[methodName]string{
	methodPing:      "hck-ping",
	methodBoxCreate: "hck-box-create",
	methodBoxExec:   "hck-box-exec",
	methodBoxDelete: "hck-box-delete",
	methodBoxList:   "hck-box-list",
}

func (c methodName) String() string {
	return methods[c]
}
