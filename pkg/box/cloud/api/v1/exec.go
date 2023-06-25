package v1

type BoxExecRequestBody struct {
	Name string `json:"name"`
}

func (b BoxExecRequestBody) method() methodName {
	return methodBoxExec
}

func NewBoxExecRequest(origin string, name string) *Message[BoxExecRequestBody] {
	return newMessage[BoxExecRequestBody](origin, BoxExecRequestBody{name})
}