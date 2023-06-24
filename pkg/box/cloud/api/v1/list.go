package v1

type BoxListRequestBody struct{}

func (b BoxListRequestBody) method() methodName {
	return methodBoxList
}

type BoxListResponseBody struct {
	Names []string `json:"names"`
}

func (b BoxListResponseBody) method() methodName {
	return methodBoxList
}

func NewBoxListRequest(origin string) *Message[BoxListRequestBody] {
	return newMessage[BoxListRequestBody](origin, BoxListRequestBody{})
}

func NewBoxListResponse(origin string, names []string) *Message[BoxListResponseBody] {
	return newMessage[BoxListResponseBody](origin, BoxListResponseBody{names})
}
