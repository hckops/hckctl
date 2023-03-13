package box

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"

	cli "github.com/hckops/hckctl/internal/model"
	"github.com/hckops/hckctl/internal/terminal"
	"github.com/hckops/hckctl/pkg/schema"
	"github.com/hckops/hckctl/pkg/util"
)

type CloudBoxCli struct {
	log      zerolog.Logger
	loader   *terminal.Loader
	config   *cli.CloudConfig
	template *schema.BoxV1 // only name is actually needed
}

func NewCloudBox(template *schema.BoxV1, config *cli.CloudConfig) *CloudBoxCli {
	l := logger.With().Str("cmd", "cloud").Logger()

	return &CloudBoxCli{
		log:      l,
		loader:   terminal.NewLoader(),
		config:   config,
		template: template,
	}
}

// TODO refactor in pkg/client
func (cli *CloudBoxCli) Open() {
	cli.log.Debug().Msgf("init cloud box:\n%v\n", cli.template.Pretty())
	cli.loader.Start(fmt.Sprintf("loading to %s/%s", cli.config.Address(), cli.template.Name))

	sshConfig := sshClientConfig(cli.config)

	client, err := ssh.Dial("tcp", cli.config.Address(), sshConfig)
	if err != nil {
		cli.loader.Halt(err, "connection failed")
	}

	session, err := client.NewSession()
	if err != nil {
		cli.loader.Halt(err, "ssh session failed")
	}
	defer session.Close()

	cli.log.Debug().Msgf("[%s] ssh connection established (%s)", client.RemoteAddr(), client.ClientVersion())

	onStreamErrorCallback := func(err error, message string) {
		cli.log.Warn().Err(err).Msg(message)
	}

	if err := handleStreams(session, onStreamErrorCallback); err != nil {
		cli.loader.Halt(err, "error streams")
	}

	terminal, err := util.NewRawTerminal(os.Stdin)
	if err != nil {
		cli.log.Warn().Err(err).Msg("error terminal")
	}
	defer terminal.Restore()

	// TODO split channel requests to show progress
	cli.loader.Stop()

	// TODO schema "{"kind":"action/v1","name":"hck-box-open","body":{"template":"alpine","revision":"main"}}"
	if err := session.Run(fmt.Sprintf("hck-box-open::%s", cli.template.Name)); err != nil && err != io.EOF {
		cli.loader.Halt(err, "error cloud box open")
	}
}

// TODO ssh agent auth
func sshClientConfig(config *cli.CloudConfig) *ssh.ClientConfig {
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
