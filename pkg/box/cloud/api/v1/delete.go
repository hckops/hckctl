package v1

type BoxDeleteRequestBody struct {
	Names []string `json:"names"`
}

func (b BoxDeleteRequestBody) method() methodName {
	return MethodBoxDelete
}

type BoxDeleteResponseBody struct {
	Names []string `json:"names"`
}

func (b BoxDeleteResponseBody) method() methodName {
	return MethodBoxDelete
}

func NewBoxDeleteRequest(origin string, names []string) *Message[BoxDeleteRequestBody] {
	return newMessage[BoxDeleteRequestBody](origin, BoxDeleteRequestBody{names})
}

func NewBoxDeleteResponse(origin string, names []string) *Message[BoxDeleteResponseBody] {
	return newMessage[BoxDeleteResponseBody](origin, BoxDeleteResponseBody{names})
}
