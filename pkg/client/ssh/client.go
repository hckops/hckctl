package ssh

import (
	"context"

	"github.com/pkg/errors"

	gossh "golang.org/x/crypto/ssh"
)

func NewSshClient(config *SshClientConfig) (*SshClient, error) {

	sshConfig := sshClientConfig(config)
	client, err := gossh.Dial("tcp", config.Address, sshConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error ssh client")
	}

	return &SshClient{
		ctx: context.Background(),
		ssh: client,
	}, nil
}

func sshClientConfig(config *SshClientConfig) *gossh.ClientConfig {
	sshConfig := &gossh.ClientConfig{
		User: config.Username,
		Auth: []gossh.AuthMethod{
			gossh.Password(config.Token),
		},
		HostKeyCallback: gossh.InsecureIgnoreHostKey(), // TODO remove
	}
	return sshConfig
}

func (client *SshClient) Close() error {
	return client.ssh.Close()
}

func (client *SshClient) SendRequest(protocol string, payload string) (string, error) {
	_, response, err := client.ssh.SendRequest(protocol, true, []byte(payload))
	if err != nil {
		return "", errors.Wrapf(err, "error send request")
	}
	return string(response), nil
}
