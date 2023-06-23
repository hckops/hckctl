package ssh

import (
	"context"

	"golang.org/x/crypto/ssh"
)

type SshClientConfig struct {
	Address  string
	Username string
	Token    string
}

type SshClient struct {
	ctx context.Context
	ssh *ssh.Client
}
