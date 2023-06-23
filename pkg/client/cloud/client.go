package cloud

import (
	"context"
	"golang.org/x/crypto/ssh"

	"github.com/pkg/errors"
)

func NewCloudClient(config *CloudClientConfig) (*CloudClient, error) {

	sshConfig := sshClientConfig(config)
	client, err := ssh.Dial("tcp", config.Address, sshConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error cloud client")
	}

	return &CloudClient{
		ctx:   context.Background(),
		cloud: client,
	}, nil
}

func sshClientConfig(config *CloudClientConfig) *ssh.ClientConfig {
	sshConfig := &ssh.ClientConfig{
		User: config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Token),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO remove
	}
	return sshConfig
}

func (client *CloudClient) Close() error {
	return client.cloud.Close()
}
