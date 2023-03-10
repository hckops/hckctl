package box

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"

	"github.com/hckops/hckctl/internal/terminal"
	"github.com/hckops/hckctl/pkg/docker"
	"github.com/hckops/hckctl/pkg/model"
)

type DockerBoxCli struct {
	// TODO dockerConfig
	loader *terminal.Loader
	log    zerolog.Logger
	box    *docker.DockerBox
}

func NewDockerBox(template *model.BoxV1, streams *model.BoxStreams) *DockerBoxCli {
	l := logger.With().Str("cmd", "docker").Logger()

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		l.Fatal().Err(err).Msg("error docker client")
	}

	return &DockerBoxCli{
		loader: terminal.NewLoader(),
		log:    l,
		box: &docker.DockerBox{
			DockerClient: dockerClient,
			Context: &model.BoxContext{
				Ctx:      context.Background(),
				Template: template,
				Streams:  streams,
			},
		},
	}
}

func (cli *DockerBoxCli) Open() {
	defer cli.box.DockerClient.Close()

	cli.log.Debug().Msgf("init docker box: \n%v\n", cli.box.Context.Template.Pretty())
	cli.loader.Start(fmt.Sprintf("loading %s", cli.box.Context.Template.Name))

	imageName := cli.box.Context.Template.ImageName()

	cli.box.Setup(func() {
		cli.loader.Refresh(fmt.Sprintf("pulling %s", imageName))
	})

	containerName := cli.box.Context.Template.GenerateName()

	cli.loader.Refresh(fmt.Sprintf("creating %s", containerName))

	containerId := cli.box.Create(containerName)

	cli.log.Info().Msgf("open new box: image=%s, containerName=%s, containerId=%s", imageName, containerName, containerId)

	cli.box.Exec(containerId, cli.loader.Stop)
}
