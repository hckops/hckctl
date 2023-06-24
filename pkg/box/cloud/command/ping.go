package command

type PingBody struct{}

func (b PingBody) cmdName() commandName {
	return commandPing
}

func NewPingRequest() *Request[PingBody] {
	return newRequest[PingBody](PingBody{})
}
