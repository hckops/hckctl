package v1

type BoxDescribeRequestBody struct {
	Name string `json:"name"`
}

func (b BoxDescribeRequestBody) method() MethodName {
	return MethodBoxDescribe
}

type BoxDescribeResponseBody struct {
	Id       string                   `json:"id"`
	Name     string                   `json:"name"`
	Created  string                   `json:"created"`
	Healthy  bool                     `json:"healthy"`
	Size     string                   `json:"size"`
	Template *BoxDescribeTemplateInfo `json:"template"`
	Env      []string                 `json:"env"`
	Ports    []string                 `json:"ports"`
}

type BoxDescribeTemplateInfo struct {
	Public   bool   `json:"public"`
	Url      string `json:"url"`
	Revision string `json:"revision"`
	Commit   string `json:"commit"`
	Name     string `json:"name"`
}

func (b BoxDescribeResponseBody) method() MethodName {
	return MethodBoxDescribe
}

func NewBoxDescribeRequest(origin string, name string) *Message[BoxDescribeRequestBody] {
	return newMessage[BoxDescribeRequestBody](origin, BoxDescribeRequestBody{Name: name})
}

func NewBoxDescribeResponse(origin string, body BoxDescribeResponseBody) *Message[BoxDescribeResponseBody] {
	return newMessage[BoxDescribeResponseBody](origin, body)
}
