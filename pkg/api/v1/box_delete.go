package v1

type BoxDeleteRequestBody struct {
	Names []string `json:"names"`
}

func (b BoxDeleteRequestBody) method() MethodName {
	return MethodBoxDelete
}

type BoxDeleteResponseBody struct {
	Items []BoxDeleteItem `json:"items"`
}

type BoxDeleteItem struct {
	Id   string
	Name string
}

func (b BoxDeleteResponseBody) method() MethodName {
	return MethodBoxDelete
}

func NewBoxDeleteRequest(origin string, names []string) *Message[BoxDeleteRequestBody] {
	return newMessage[BoxDeleteRequestBody](origin, BoxDeleteRequestBody{Names: names})
}

func NewBoxDeleteResponse(origin string, items []BoxDeleteItem) *Message[BoxDeleteResponseBody] {
	return newMessage[BoxDeleteResponseBody](origin, BoxDeleteResponseBody{Items: items})
}
