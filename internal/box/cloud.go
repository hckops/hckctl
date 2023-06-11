package box

import (
	"encoding/hex"
	"fmt"
	"github.com/hckops/hckctl/internal/old/common"
	"github.com/hckops/hckctl/internal/old/schema"
	common2 "github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/util"
	"io"
	"net"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"

	"github.com/hckops/hckctl/internal/config"
)

type RemoteSshBox struct {
	log      zerolog.Logger
	loader   *common2.Loader
	config   *config.CloudConfig
	revision string
	template *schema.BoxV1 // only name is actually needed
	client   *ssh.Client
}

func NewRemoteSshBox(template *schema.BoxV1, revision string, config *config.CloudConfig) *RemoteSshBox {
	l := logger.With().Str("provider", "cloud").Logger()

	return &RemoteSshBox{
		log:      l,
		loader:   common2.NewLoader(),
		config:   config,
		revision: revision,
		template: template,
	}
}

func (remote *RemoteSshBox) Open() {
	defer remote.loader.Stop()

	remote.log.Debug().Msgf("init cloud box:\n%v\n", remote.template.Pretty())
	remote.loader.Start(fmt.Sprintf("loading to %s/%s", remote.config.Address(), remote.template.Name))
	remote.loader.Sleep(1)

	sshConfig := sshClientConfig(remote.config)

	client, err := ssh.Dial("tcp", remote.config.Address(), sshConfig)
	if err != nil {
		remote.loader.Halt(err, "connection error")
	}
	defer client.Close()
	remote.client = client

	remote.log.Info().
		Str("User", client.User()).
		Str("ClientVersion", string(client.ClientVersion())).
		Str("ServerVersion", string(client.ServerVersion())).
		Str("RemoteAddress", client.RemoteAddr().String()).
		Str("LocalAddress", client.LocalAddr().String()).
		Str("ConnectionId", hex.EncodeToString(client.SessionID())).
		Msg("ssh connection established")

	boxId := remote.create()
	defer remote.delete(boxId)

	remote.loader.Refresh(fmt.Sprintf("tunneling %s", boxId))
	remote.tunnelBox(boxId)
	remote.exec(boxId)
}

func (remote *RemoteSshBox) create() string {

	boxId := remote.sendRequest(common.NewCommandCreateBox(remote.template.Name, remote.revision))
	remote.log.Info().Msgf("create cloud box: %s", boxId)
	return boxId
}

func (remote *RemoteSshBox) delete(boxId string) {

	_ = remote.sendRequest(common.NewCommandDeleteBox(remote.template.Name, remote.revision, boxId))
	remote.log.Info().Msgf("delete cloud box: %s", boxId)
}

func (remote *RemoteSshBox) sendRequest(payload string) string {
	remote.log.Debug().Msgf("send request [%s]", payload)

	_, response, err := remote.client.SendRequest(common.CommandRequestType, true, []byte(payload))
	if err != nil || string(response) == common.CommandResponseError {
		remote.loader.Halt(err, "error cloud: send request")
	}

	remote.log.Debug().Msgf("response [%s]", response)
	return string(response)
}

func (remote *RemoteSshBox) tunnelBox(boxId string) {

	for _, port := range remote.template.NetworkPorts() {
		localPort, _ := util.FindOpenPort(port.Local)

		openPort := schema.PortV1{
			Alias:  port.Alias,
			Local:  localPort,
			Remote: port.Remote,
		}

		message := fmt.Sprintf("[%s][%s] tunnel %s (local) -> %s (remote)", boxId, port.Alias, port.Local, port.Remote)
		remote.log.Info().Msgf(message)
		// prints to terminal
		fmt.Println(message)
		go remote.tunnel(boxId, openPort)
	}
}

func (remote *RemoteSshBox) tunnel(boxId string, port schema.PortV1) {

	// starts local server to forward traffic to remote connection
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port.Local))
	if err != nil {
		remote.loader.Halt(err, "error cloud: tunnel listen")
	}
	defer listener.Close()

	copyConnection := func(writer, reader net.Conn) {
		defer writer.Close()
		defer reader.Close()

		_, err := io.Copy(writer, reader)
		if err != nil {
			remote.log.Warn().Err(err).Msg("error copy connection")
		}
	}

	for {
		localConnection, err := listener.Accept()
		if err != nil {
			remote.loader.Halt(err, "error cloud: local tunnel")
		}
		// forward
		go func() {
			remoteConnection, err := remote.client.Dial("tcp", fmt.Sprintf("%s:%s", boxId, port.Remote))
			if err != nil {
				remote.loader.Halt(err, "error cloud: remote tunnel")
			}

			go copyConnection(localConnection, remoteConnection)
			go copyConnection(remoteConnection, localConnection)
		}()
	}
}

func (remote *RemoteSshBox) exec(boxId string) {

	session, err := remote.client.NewSession()
	if err != nil {
		remote.loader.Halt(err, "ssh session error")
	}
	defer session.Close()

	onStreamErrorCallback := func(err error, message string) {
		remote.log.Warn().Err(err).Msg(message)
	}
	if err := handleStreams(session, onStreamErrorCallback); err != nil {
		remote.loader.Halt(err, "error streams")
	}

	if rawTerminal, err := util.NewRawTerminal(os.Stdin); err == nil {
		defer rawTerminal.Restore()
	} else {
		remote.log.Warn().Err(err).Msg("error terminal")
	}

	remote.loader.Stop()

	payload := common.NewCommandExecBox(remote.template.Name, remote.revision, boxId)
	if err := session.Run(payload); err != nil && err != io.EOF {
		remote.loader.Halt(err, "error cloud box exec")
	}
}

// TODO ssh agent auth
func sshClientConfig(config *config.CloudConfig) *ssh.ClientConfig {
	sshConfig := &ssh.ClientConfig{
		User: config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Token),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO remove
	}
	return sshConfig
}

func handleStreams(session *ssh.Session, onStreamErrorCallback func(error, string)) error {

	stdin, err := session.StdinPipe()
	if err != nil {
		return errors.Wrap(err, "unable to setup stdin for session")
	}
	go func() {
		if _, err := io.Copy(stdin, os.Stdin); err != nil {
			onStreamErrorCallback(err, "error copy stdin local->remote")
		}
	}()

	stdout, err := session.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "unable to setup stdout for session")
	}
	go func() {
		if _, err := io.Copy(os.Stdout, stdout); err != nil {
			onStreamErrorCallback(err, "error copy stdout remote->local")
		}
	}()

	stderr, err := session.StderrPipe()
	if err != nil {
		return errors.Wrap(err, "unable to setup stderr for session")
	}
	go func() {
		if _, err := io.Copy(os.Stderr, stderr); err != nil {
			onStreamErrorCallback(err, "error copy stderr remote->local")
		}
	}()

	return nil
}
