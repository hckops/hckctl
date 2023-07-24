package v1

type BoxExecSessionBody struct {
	Name string `json:"name"`
}

func (b BoxExecSessionBody) method() MethodName {
	return MethodBoxExec
}

func NewBoxExecSession(origin string, name string) *Message[BoxExecSessionBody] {
	return newMessage[BoxExecSessionBody](origin, BoxExecSessionBody{Name: name})
}
