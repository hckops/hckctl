package v1

type BoxCreateRequestBody struct {
	TemplateName string `json:"templateName"`
	Size         string `json:"size"` // expected
}

func (b BoxCreateRequestBody) method() MethodName {
	return MethodBoxCreate
}

type BoxCreateResponseBody struct {
	Name string `json:"name"`
	Size string `json:"size"` // actual
}

func (b BoxCreateResponseBody) method() MethodName {
	return MethodBoxCreate
}

func NewBoxCreateRequest(origin string, templateName string, size string) *Message[BoxCreateRequestBody] {
	return newMessage[BoxCreateRequestBody](origin, BoxCreateRequestBody{templateName, size})
}

func NewBoxCreateResponse(origin string, name string, size string) *Message[BoxCreateResponseBody] {
	return newMessage[BoxCreateResponseBody](origin, BoxCreateResponseBody{name, size})
}
