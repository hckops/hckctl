package v1

type BoxListRequestBody struct{}

func (b BoxListRequestBody) method() MethodName {
	return MethodBoxList
}

type BoxListResponseBody struct {
	Items []BoxListItem `json:"items"`
}

type BoxListItem struct {
	Id      string
	Name    string
	Healthy bool
}

func (b BoxListResponseBody) method() MethodName {
	return MethodBoxList
}

func NewBoxListRequest(origin string) *Message[BoxListRequestBody] {
	return newMessage[BoxListRequestBody](origin, BoxListRequestBody{})
}

func NewBoxListResponse(origin string, items []BoxListItem) *Message[BoxListResponseBody] {
	return newMessage[BoxListResponseBody](origin, BoxListResponseBody{items})
}
