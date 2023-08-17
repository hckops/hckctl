package v1

type LabCreateRequestBody struct {
	TemplateName string            `json:"templateName"`
	Parameters   map[string]string `json:"params"` // expected
}

func (b LabCreateRequestBody) method() MethodName {
	return MethodLabCreate
}

type LabCreateResponseBody struct {
	Name string `json:"name"`
}

func (b LabCreateResponseBody) method() MethodName {
	return MethodLabCreate
}

func NewLabCreateRequest(origin string, templateName string, parameters map[string]string) *Message[LabCreateRequestBody] {
	return newMessage[LabCreateRequestBody](origin, LabCreateRequestBody{TemplateName: templateName, Parameters: parameters})
}

func NewLabCreateResponse(origin string, name string) *Message[LabCreateResponseBody] {
	return newMessage[LabCreateResponseBody](origin, LabCreateResponseBody{Name: name})
}
