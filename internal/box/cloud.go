package box

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"

	"github.com/hckops/hckctl/internal/model"
	"github.com/hckops/hckctl/internal/terminal"
)

type CloudBoxCli struct {
	log    zerolog.Logger
	loader *terminal.Loader
	config *model.CloudConfig
}

func NewCloudBox(name string, config *model.CloudConfig) *CloudBoxCli {
	l := logger.With().Str("cmd", "docker").Logger()

	return &CloudBoxCli{
		log:    l,
		loader: terminal.NewLoader(),
		config: config,
	}
}

func (cli *CloudBoxCli) Open() {

	address := net.JoinHostPort(cli.config.Host, strconv.Itoa(cli.config.Port))
	sshConfig := sshClientConfig(cli.config)
	connect(address, sshConfig)

}

func sshClientConfig(config *model.CloudConfig) *ssh.ClientConfig {
	sshConfig := &ssh.ClientConfig{
		User: config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Token),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return sshConfig
}

func connect(address string, sshConfig *ssh.ClientConfig) {
	client, err := ssh.Dial("tcp", address, sshConfig)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	log.Printf("[%s] new ssh connection (%s)", client.RemoteAddr(), client.ClientVersion())

	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("failed to create session: %v", err)
	}
	log.Printf("[%s] new ssh connection", client.RemoteAddr())

	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatalf("Unable to setup stdin for session: %v", err)
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("Unable to setup stdout for session: %v", err)
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		log.Fatalf("Unable to setup stderr for session: %v", err)
	}
	go io.Copy(os.Stderr, stderr)

	err = session.Run("docker")

	defer session.Close()
}
