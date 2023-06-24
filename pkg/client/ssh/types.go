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

type TunnelOpts struct {
	LocalPort             string
	RemoteHost            string
	RemotePort            string
	OnTunnelErrorCallback func(error)
}

func (t *TunnelOpts) Network() string {
	return "tcp"
}

func (t *TunnelOpts) LocalAddress() string {
	return fmt.Sprintf("0.0.0.0:%s", t.LocalPort)
}
func (t *TunnelOpts) RemoteAddress() string {
	return fmt.Sprintf("%s:%s", t.RemoteHost, t.RemotePort)
}

type ExecOpts struct {
	Command               string
	OnStreamStartCallback func()
	OnStreamErrorCallback func(error)
}
