package v1

type BoxCreateRequestBody struct {
	TemplateName string `json:"templateName"`
}

func (b BoxCreateRequestBody) method() methodName {
	return methodBoxCreate
}

type BoxCreateResponseBody struct {
	Name string `json:"name"`
}

func (b BoxCreateResponseBody) method() methodName {
	return methodBoxCreate
}

func NewBoxCreateRequest(origin string, templateName string) *Message[BoxCreateRequestBody] {
	return newMessage[BoxCreateRequestBody](origin, BoxCreateRequestBody{templateName})
}

func NewBoxCreateResponse(origin string, name string) *Message[BoxCreateResponseBody] {
	return newMessage[BoxCreateResponseBody](origin, BoxCreateResponseBody{name})
}