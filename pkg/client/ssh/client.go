package ssh

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/pkg/errors"
	gossh "golang.org/x/crypto/ssh"

	"github.com/hckops/hckctl/pkg/client/common"
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

// TODO ssh agent auth
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
	// "wantReply" must be true to get a response
	ok, response, err := client.ssh.SendRequest(protocol, true, []byte(payload))
	if !ok {
		if err != nil {
			return "", errors.Wrapf(err, "error ssh send request")
		} else if strings.TrimSpace(string(response)) != "" {
			return "", fmt.Errorf("error ssh server response %s", response)
		} else {
			return "", errors.New("error ssh invalid request")
		}
	}
	return string(response), nil
}

// Tunnel starts a local server and forwards traffic to a remote connection
func (client *SshClient) Tunnel(opts *SshTunnelOpts) {

	listener, err := net.Listen(opts.Network(), opts.LocalAddress())
	if err != nil {
		opts.OnTunnelErrorCallback(errors.Wrapf(err, "error ssh creating local tunnel: address=%s", opts.LocalAddress()))
	}
	defer listener.Close()

	copyStream := func(writer, reader net.Conn, label string) {
		defer writer.Close()
		defer reader.Close()

		_, err := io.Copy(writer, reader)
		if err != nil {
			opts.OnTunnelErrorCallback(errors.Wrapf(err, "error ssh copying stream: %s", label))
		}
	}

	for {
		localConnection, err := listener.Accept()
		if err != nil {
			opts.OnTunnelErrorCallback(errors.Wrapf(err, "error ssh opening local tunnel: address=%s", opts.LocalAddress()))
		}
		// forward connections
		go func() {
			remoteConnection, err := client.ssh.Dial(opts.Network(), opts.RemoteAddress())
			if err != nil {
				opts.OnTunnelErrorCallback(errors.Wrapf(err, "error ssh opening remote tunnel: address=%s", opts.RemoteAddress()))
			}

			go copyStream(localConnection, remoteConnection, "remote->local")
			go copyStream(remoteConnection, localConnection, "local->remote")
		}()
	}
}

func (client *SshClient) Exec(opts *SshExecOpts) error {

	session, err := client.ssh.NewSession()
	if err != nil {
		return errors.Wrapf(err, "error ssh new session")
	}
	defer session.Close()

	if err := handleStreams(session, opts.OnStreamErrorCallback); err != nil {
		return errors.Wrapf(err, "error ssh stream")
	}

	terminal, err := common.NewRawTerminal(os.Stdin)
	if err != nil {
		return errors.Wrap(err, "error ssh terminal")
	}
	defer terminal.Restore()

	opts.OnStreamStartCallback()

	if err := session.Run(opts.Payload); err != nil && err != io.EOF {
		return errors.Wrapf(err, "error ssh exec session")
	}
	return nil
}

func handleStreams(session *gossh.Session, onStreamErrorCallback func(error)) error {

	stdin, err := session.StdinPipe()
	if err != nil {
		return errors.Wrap(err, "error opening stdin pipe")
	}
	go func() {
		if _, err := io.Copy(stdin, os.Stdin); err != nil {
			onStreamErrorCallback(errors.Wrap(err, "error copy stdin local->remote"))
		}
	}()

	stdout, err := session.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "error opening stdout pipe")
	}
	go func() {
		if _, err := io.Copy(os.Stdout, stdout); err != nil {
			onStreamErrorCallback(errors.Wrap(err, "error copy stdout remote->local"))
		}
	}()

	stderr, err := session.StderrPipe()
	if err != nil {
		return errors.Wrap(err, "error opening stderr pipe")
	}
	go func() {
		if _, err := io.Copy(os.Stderr, stderr); err != nil {
			onStreamErrorCallback(errors.Wrap(err, "error copy stderr remote->local"))
		}
	}()

	return nil
}
