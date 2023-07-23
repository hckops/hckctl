package v1

type BoxDeleteRequestBody struct {
	Names []string `json:"names"`
}

func (b BoxDeleteRequestBody) method() MethodName {
	return MethodBoxDelete
}

type BoxDeleteResponseBody struct {
	Names []string `json:"names"`
}

func (b BoxDeleteResponseBody) method() MethodName {
	return MethodBoxDelete
}

func NewBoxDeleteRequest(origin string, names []string) *Message[BoxDeleteRequestBody] {
	return newMessage[BoxDeleteRequestBody](origin, BoxDeleteRequestBody{Names: names})
}

func NewBoxDeleteResponse(origin string, names []string) *Message[BoxDeleteResponseBody] {
	return newMessage[BoxDeleteResponseBody](origin, BoxDeleteResponseBody{Names: names})
}
