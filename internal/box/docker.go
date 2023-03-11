package box

import (
	"fmt"

	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"

	"github.com/hckops/hckctl/internal/terminal"
	"github.com/hckops/hckctl/pkg/client"
	"github.com/hckops/hckctl/pkg/model"
)

type DockerBoxCli struct {
	// TODO dockerConfig
	log     zerolog.Logger
	loader  *terminal.Loader
	box     *client.DockerBox
	streams *model.BoxStreams
}

func NewDockerBox(template *model.BoxV1) *DockerBoxCli {
	l := logger.With().Str("cmd", "docker").Logger()

	box, err := client.NewDockerBox(template)
	if err != nil {
		l.Fatal().Err(err).Msg("error docker box")
	}

	return &DockerBoxCli{
		loader:  terminal.NewLoader(),
		log:     l,
		box:     box,
		streams: model.NewDefaultStreams(true), // TODO tty
	}
}

func (cli *DockerBoxCli) Open() {
	defer cli.box.Close()

	cli.log.Debug().Msgf("init docker box:\n%v\n", cli.box.Template.Pretty())
	cli.loader.Start(fmt.Sprintf("loading %s", cli.box.Template.Name))

	imageName := cli.box.Template.ImageName()

	cli.box.OnSetupCallback = func() {
		cli.loader.Refresh(fmt.Sprintf("pulling %s", imageName))
	}
	if err := cli.box.Setup(); err != nil {
		cli.shutDown(err, "error docker box setup")
	}

	containerName := cli.box.Template.GenerateName()
	cli.loader.Refresh(fmt.Sprintf("creating %s", containerName))

	cli.box.OnCreateCallback = func(port model.PortV1) {
		cli.log.Info().Msgf("[%s] exposing %s (local) -> %s (container)", port.Alias, port.Local, port.Remote)
	}
	containerId, err := cli.box.Create(containerName)
	if err != nil {
		cli.shutDown(err, "error docker box create")
	}

	cli.log.Info().Msgf("opening new box: image=%s, containerName=%s, containerId=%s", imageName, containerName, containerId)

	cli.box.OnCloseCallback = func() {
		cli.log.Debug().Msgf("removing container: %s", containerId)
	}
	cli.box.OnCloseErrorCallback = func(err error, message string) {
		cli.shutDown(err, message)
	}
	cli.box.OnStreamErrorCallback = func(err error, message string) {
		cli.log.Warn().Err(err).Msg(message)
	}

	cli.box.Exec(containerId, cli.streams, cli.loader.Stop)
}

func (cli *DockerBoxCli) shutDown(err error, message string) {
	cli.loader.Stop()
	fmt.Println(message)
	cli.log.Fatal().Err(err).Msg(message)
}
