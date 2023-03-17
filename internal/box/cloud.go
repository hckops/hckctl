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

// TODO refactor cloud in pkg/client
func (remote *RemoteSshBox) Open() {
	remote.log.Debug().Msgf("init cloud box:\n%v\n", remote.template.Pretty())
	remote.loader.Start(fmt.Sprintf("loading to %s/%s", remote.config.Address(), remote.template.Name))

	sshConfig := sshClientConfig(remote.config)

	client, err := ssh.Dial("tcp", remote.config.Address(), sshConfig)
	if err != nil {
		remote.loader.Halt(err, "connection error")
	}

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
		Msg("ssh connection established")

	onStreamErrorCallback := func(err error, message string) {
		remote.log.Warn().Err(err).Msg(message)
	}
	if err := handleStreams(session, onStreamErrorCallback); err != nil {
		remote.loader.Halt(err, "error streams")
	}

	terminal, err := util.NewRawTerminal(os.Stdin)
	if err != nil {
		remote.log.Warn().Err(err).Msg("error terminal")
	}
	defer terminal.Restore()

	// TODO split channel requests to show progress
	remote.loader.Stop()

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
