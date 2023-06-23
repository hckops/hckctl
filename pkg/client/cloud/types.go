package cloud

import (
	"context"

	"golang.org/x/crypto/ssh"
)

type CloudClientConfig struct {
	Address  string
	Username string
	Token    string
}

type CloudClient struct {
	ctx   context.Context
	cloud *ssh.Client
}
