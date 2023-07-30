package ssh

import (
	"context"
	"fmt"

	gossh "golang.org/x/crypto/ssh"
)

type SshClient struct {
	ctx context.Context
	ssh *gossh.Client
}

type SshClientConfig struct {
	Address  string
	Username string
	Token    string
}

type SshTunnelOpts struct {
	LocalPort             string
	RemoteHost            string
	RemotePort            string
	OnTunnelStartCallback func(string)
	OnTunnelStopCallback  func(string)
	OnTunnelErrorCallback func(error)
}

func (t *SshTunnelOpts) Network() string {
	return "tcp"
}

func (t *SshTunnelOpts) LocalAddress() string {
	return fmt.Sprintf("0.0.0.0:%s", t.LocalPort)
}
func (t *SshTunnelOpts) RemoteAddress() string {
	return fmt.Sprintf("%s:%s", t.RemoteHost, t.RemotePort)
}

type SshExecOpts struct {
	Payload               string
	OnStreamStartCallback func()
	OnStreamErrorCallback func(error)
}
