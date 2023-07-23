package v1

type BoxDescribeRequestBody struct {
	Name string `json:"name"`
}

func (b BoxDescribeRequestBody) method() MethodName {
	return MethodBoxDescribe
}

// TODO
type BoxDescribeResponseBody struct {
	Name string `json:"name"`
}

func (b BoxDescribeResponseBody) method() MethodName {
	return MethodBoxDescribe
}

func NewBoxDescribeRequest(origin string, name string) *Message[BoxDescribeRequestBody] {
	return newMessage[BoxDescribeRequestBody](origin, BoxDescribeRequestBody{Name: name})
}

func NewBoxDescribeResponse(origin string, name string) *Message[BoxDescribeResponseBody] {
	return newMessage[BoxDescribeResponseBody](origin, BoxDescribeResponseBody{Name: name})
}
