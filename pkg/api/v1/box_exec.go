package v1

// TODO add tunnel options

type BoxExecSessionBody struct {
	Name string `json:"name"`
}

func (b BoxExecSessionBody) method() MethodName {
	return MethodBoxExec
}

func NewBoxExecSession(origin string, name string) *Message[BoxExecSessionBody] {
	return newMessage[BoxExecSessionBody](origin, BoxExecSessionBody{Name: name})
}
