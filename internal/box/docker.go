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
	// TODO is this the right place?
	defer box.Close()

	return &DockerBoxCli{
		loader:  terminal.NewLoader(),
		log:     l,
		box:     box,
		streams: model.NewDefaultStreams(true), // TODO tty
	}
}

func (cli *DockerBoxCli) Open() {
	cli.log.Debug().Msgf("init docker box:\n%v\n", cli.box.Template.Pretty())
	cli.loader.Start(fmt.Sprintf("loading %s", cli.box.Template.Name))

	imageName := cli.box.Template.ImageName()

	// TODO add sleep / move loader outside
	if err := cli.box.Setup(func() {
		cli.loader.Refresh(fmt.Sprintf("pulling %s", imageName))
	}); err != nil {
		cli.log.Fatal().Err(err).Msg("error docker box setup")
	}

	containerName := cli.box.Template.GenerateName()

	cli.loader.Refresh(fmt.Sprintf("creating %s", containerName))

	containerId, err := cli.box.Create(containerName, func(port model.PortV1) {
		cli.log.Info().Msgf("[%s] exposing %s (local) -> %s (container)", port.Alias, port.Local, port.Remote)
	})
	if err != nil {
		cli.log.Fatal().Err(err).Msg("error docker box create")
	}

	cli.log.Info().Msgf("opening new box: image=%s, containerName=%s, containerId=%s", imageName, containerName, containerId)

	cli.box.Exec(containerId, cli.streams, cli.loader.Stop)
}
