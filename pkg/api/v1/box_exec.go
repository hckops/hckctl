package v1

type BoxExecSessionBody struct {
	Name       string `json:"name"`
	TunnelOnly bool   `json:"tunnelOnly"`
	NoTunnel   bool   `json:"noTunnel"`
}

func (b BoxExecSessionBody) method() MethodName {
	return MethodBoxExec
}

func NewBoxExecSession(origin string, name string, tunnelOnly bool, noTunnel bool) *Message[BoxExecSessionBody] {
	return newMessage[BoxExecSessionBody](origin, BoxExecSessionBody{Name: name, TunnelOnly: tunnelOnly, NoTunnel: noTunnel})
}
