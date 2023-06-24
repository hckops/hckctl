package v1

type BoxDeleteBody struct {
	Name string `json:"name"`
}

func (b BoxDeleteBody) name() requestName {
	return requestBoxDelete
}

func NewBoxDeleteRequest(name string) *Request[BoxDeleteBody] {
	return newRequest[BoxDeleteBody](BoxDeleteBody{name})
}
