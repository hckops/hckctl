package v1

const (
	PingValue = "ping"
	PongValue = "pong"
)

type PingBody struct {
	Value string `json:"value"`
}

func (b PingBody) method() MethodName {
	return MethodPing
}

type PongBody struct {
	Value string `json:"value"`
}

func (b PongBody) method() MethodName {
	return MethodPing
}

func NewPingMessage(origin string) *Message[PingBody] {
	return newMessage[PingBody](origin, PingBody{PingValue})
}

func NewPongMessage(origin string) *Message[PongBody] {
	return newMessage[PongBody](origin, PongBody{PongValue})
}
