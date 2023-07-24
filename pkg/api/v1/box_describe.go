package v1

type BoxDescribeRequestBody struct {
	Name string `json:"name"`
}

func (b BoxDescribeRequestBody) method() MethodName {
	return MethodBoxDescribe
}

type BoxDescribeResponseBody struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Created string   `json:"created"`
	Healthy bool     `json:"healthy"`
	Size    string   `json:"size"`
	Env     []string `json:"env"`
	Ports   []string `json:"ports"`
}

func (b BoxDescribeResponseBody) method() MethodName {
	return MethodBoxDescribe
}

func NewBoxDescribeRequest(origin string, name string) *Message[BoxDescribeRequestBody] {
	return newMessage[BoxDescribeRequestBody](origin, BoxDescribeRequestBody{Name: name})
}

func NewBoxDescribeResponse(
	origin string,
	id string,
	name string,
	created string,
	healthy bool,
	size string,
	env []string,
	ports []string,
) *Message[BoxDescribeResponseBody] {
	return newMessage[BoxDescribeResponseBody](origin, BoxDescribeResponseBody{
		Id:      id,
		Name:    name,
		Created: created,
		Healthy: healthy,
		Size:    size,
		Env:     env,
		Ports:   ports,
	})
}
