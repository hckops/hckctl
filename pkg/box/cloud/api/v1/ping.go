package v1

type PingBody struct{}

func (b PingBody) name() requestName {
	return requestPing
}

func NewPingRequest() *Request[PingBody] {
	return newRequest[PingBody](PingBody{})
}
