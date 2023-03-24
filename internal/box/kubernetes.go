package box

import (
	"fmt"

	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"

	"github.com/hckops/hckctl/internal/config"
	"github.com/hckops/hckctl/internal/terminal"
	"github.com/hckops/hckctl/pkg/client"
	"github.com/hckops/hckctl/pkg/model"
	"github.com/hckops/hckctl/pkg/schema"
)

type LocalKubeBox struct {
	log     zerolog.Logger
	loader  *terminal.Loader
	box     *client.KubeBox
	streams *model.BoxStreams
}

func NewKubeBox(template *schema.BoxV1, config *config.KubeConfig) *LocalKubeBox {
	l := logger.With().Str("provider", "kube").Logger()

	box, err := client.NewOutOfClusterKubeBox(
		template,
		&client.ResourceOptions{
			Namespace: config.Namespace,
			Memory:    config.Resources.Memory,
			Cpu:       config.Resources.Cpu,
		},
		config.ConfigPath,
	)
	if err != nil {
		l.Fatal().Err(err).Msg("error kube box")
	}

	return &LocalKubeBox{
		log:     l,
		loader:  terminal.NewLoader(),
		box:     box,
		streams: model.NewDefaultStreams(true), // TODO tty
	}
}

func (local *LocalKubeBox) Open() {
	local.log.Debug().Msgf("init kube box:\n%v\n", local.box.Template.Pretty())
	local.loader.Start(fmt.Sprintf("loading %s", local.box.Template.Name))
	local.loader.Sleep(1)

	containerName := local.box.Template.GenerateName()
	deployment, service, err := local.box.BuildSpec(containerName)
	if err != nil {
		local.loader.Halt(err, "error kube: invalid template")
	}

	local.box.OnSetupCallback = func(message string) {
		local.log.Debug().Msg(message)
	}
	local.loader.Refresh(fmt.Sprintf("creating %s/%s", local.box.ResourceOptions.Namespace, containerName))
	err = local.box.ApplyTemplate(deployment, service)
	if err != nil {
		local.loader.Halt(err, "error kube: apply template")
	}

	local.box.OnCloseCallback = func(message string) {
		local.log.Debug().Msg(message)
	}
	local.box.OnCloseErrorCallback = func(err error, message string) {
		local.log.Warn().Err(err).Msg(message)
	}
	defer local.box.RemoveTemplate(deployment, service)

	pod, err := local.box.GetPod(deployment)
	if err != nil {
		local.loader.Halt(err, "error kube: invalid pod")
	}
	local.log.Debug().Msgf("found matching pod %s", pod.Name)

	local.log.Info().Msgf("opening new box: image=%s, namespace=%s, podName=%s", local.box.Template.ImageName(), pod.Namespace, pod.Name)

	local.box.OnTunnelCallback = func(port schema.PortV1) {
		local.log.Info().Msgf("[%s][%s] forwarding %s (local) -> %s (remote)", pod.Name, port.Alias, port.Local, port.Remote)
	}
	local.box.OnTunnelErrorCallback = func(err error, message string) {
		local.loader.Halt(err, "error kube: port forward")
	}
	local.box.PortForward(pod.Name, pod.Namespace)

	local.box.OnExecCallback = func() {
		local.log.Debug().Msgf("exec into pod: %s", pod.Name)
		local.loader.Stop()
	}
	if err := local.box.Exec(pod, local.streams); err != nil {
		local.loader.Halt(err, "error kube: exec pod")
	}
}
