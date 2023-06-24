package v1

type BoxCreateBody struct {
	TemplateName string `json:"templateName"`
}

func (b BoxCreateBody) name() requestName {
	return requestBoxCreate
}

func NewBoxCreateRequest(templateName string) *Request[BoxCreateBody] {
	body := BoxCreateBody{
		TemplateName: templateName,
	}
	return newRequest[BoxCreateBody](body)
}
