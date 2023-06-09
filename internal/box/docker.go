package box

import (
	"fmt"
	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/old/client"
	"github.com/hckops/hckctl/pkg/old/model"
	"github.com/hckops/hckctl/pkg/old/schema"

	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
)

type LocalDockerBox struct {
	// TODO dockerConfig
	log     zerolog.Logger
	loader  *common.Loader
	box     *client.DockerBox
	streams *model.BoxStreams
}

func NewDockerBox(template *schema.BoxV1) *LocalDockerBox {
	l := logger.With().Str("provider", "docker").Logger()

	box, err := client.NewDockerBox(template)
	if err != nil {
		l.Fatal().Err(err).Msg("error docker box")
	}

	return &LocalDockerBox{
		log:     l,
		loader:  common.NewLoader(),
		box:     box,
		streams: model.NewDefaultStreams(true), // TODO tty
	}
}

func (local *LocalDockerBox) Open() {
	defer local.box.Close()

	local.log.Debug().Msgf("init docker box:\n%v\n", local.box.Template.Pretty())
	local.loader.Start(fmt.Sprintf("loading %s", local.box.Template.Name))

	imageName := local.box.Template.ImageName()

	local.box.OnSetupCallback = func() {
		local.loader.Refresh(fmt.Sprintf("pulling %s", imageName))
	}
	if err := local.box.Setup(); err != nil {
		local.loader.Halt(err, "error docker box setup")
	}

	containerName := local.box.Template.GenerateName()
	local.loader.Refresh(fmt.Sprintf("creating %s", containerName))

	local.box.OnCreateCallback = func(port schema.PortV1) {
		local.log.Info().Msgf("[%s][%s] exposing %s (local) -> %s (container)", containerName, port.Alias, port.Local, port.Remote)
	}
	containerId, err := local.box.Create(containerName)
	if err != nil {
		local.loader.Halt(err, "error docker box create")
	}

	local.log.Info().Msgf("opening new box: image=%s, containerName=%s, containerId=%s", imageName, containerName, containerId)

	local.box.OnExecCallback = func() {
		local.loader.Stop()
	}
	local.box.OnCloseCallback = func() {
		local.log.Debug().Msgf("removing container: %s", containerId)
	}
	local.box.OnCloseErrorCallback = func(err error, message string) {
		local.log.Warn().Err(err).Msg(message)
	}
	local.box.OnStreamErrorCallback = func(err error, message string) {
		local.log.Warn().Err(err).Msg(message)
	}
	if err := local.box.Exec(containerId, local.streams); err != nil {
		local.loader.Halt(err, "error docker box exec")
	}
}
