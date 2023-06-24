package v1

type BoxListBody struct{}

func (b BoxListBody) name() requestName {
	return requestBoxList
}

func NewBoxListRequest() *Request[BoxListBody] {
	return newRequest[BoxListBody](BoxListBody{})
}
