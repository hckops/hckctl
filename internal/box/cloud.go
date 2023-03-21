package box

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"

	"github.com/hckops/hckctl/internal/config"
	"github.com/hckops/hckctl/internal/terminal"
	"github.com/hckops/hckctl/pkg/common"
	"github.com/hckops/hckctl/pkg/schema"
	"github.com/hckops/hckctl/pkg/util"
)

type RemoteSshBox struct {
	log      zerolog.Logger
	loader   *terminal.Loader
	config   *config.CloudConfig
	revision string
	template *schema.BoxV1 // only name is actually needed
}

func NewRemoteSshBox(template *schema.BoxV1, revision string, config *config.CloudConfig) *RemoteSshBox {
	l := logger.With().Str("provider", "cloud").Logger()

	return &RemoteSshBox{
		log:      l,
		loader:   terminal.NewLoader(),
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

	remote.create(client, remote.template.Name, remote.revision)
}

func (remote *RemoteSshBox) create(client *ssh.Client, name, revision string) string {
	wantReply, response, err := client.SendRequest(string(common.CommandBoxCreate), true, []byte(common.NewCommandCreateBox(name, revision)))

	remote.log.Debug().Msgf("wantReply=%v", wantReply)
	remote.log.Debug().Msgf("response=%s", response)
	remote.log.Debug().Msgf("err=%v", err)
	return ""
}

// TODO refactor cloud in pkg/client
func (remote *RemoteSshBox) OpenOld() {
	remote.log.Debug().Msgf("init cloud box:\n%v\n", remote.template.Pretty())
	remote.loader.Start(fmt.Sprintf("loading to %s/%s", remote.config.Address(), remote.template.Name))
	remote.loader.Sleep(1)

	sshConfig := sshClientConfig(remote.config)

	client, err := ssh.Dial("tcp", remote.config.Address(), sshConfig)
	if err != nil {
		remote.loader.Halt(err, "connection error")
	}

	// TODO timeout if no response

	// TODO sent CREATE requests, wait for BOX-ID/OK
	//client.SendRequest()

	// TODO init direct-tcpip
	// client.Dial()
	// TODO start local server to forward traffic to remote connection
	// net.Listen

	session, err := client.NewSession()
	if err != nil {
		remote.loader.Halt(err, "ssh session error")
	}
	defer session.Close()

	remote.log.Info().
		Str("User", client.User()).
		Str("ClientVersion", string(client.ClientVersion())).
		Str("ServerVersion", string(client.ServerVersion())).
		Str("RemoteAddress", client.RemoteAddr().String()).
		Str("LocalAddress", client.LocalAddr().String()).
		Str("ConnectionId", string(client.SessionID())).
		Msg("ssh connection established")

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

	// TODO split channel requests to show progress
	remote.loader.Stop()

	// this is session EXEC
	if err := session.Run(common.NewCommandOpenBox(remote.template.Name, remote.revision)); err != nil && err != io.EOF {
		remote.loader.Halt(err, "error cloud box open")
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

// TODO rename handleExec
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
