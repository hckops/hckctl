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

type KubeBoxCli struct {
	log     zerolog.Logger
	loader  *terminal.Loader
	box     *client.KubeBox
	streams *model.BoxStreams
}

func NewKubeBox(template *schema.BoxV1, config *config.KubeConfig) *KubeBoxCli {
	l := logger.With().Str("cmd", "kube").Logger()

	box, err := client.NewOutOfClusterKubeBox(
		template,
		config.ConfigPath,
		&client.ResourceOptions{
			Namespace: config.Namespace,
			Memory:    config.Resources.Memory,
			Cpu:       config.Resources.Cpu,
		},
	)
	if err != nil {
		l.Fatal().Err(err).Msg("error kube box")
	}

	return &KubeBoxCli{
		log:     l,
		loader:  terminal.NewLoader(),
		box:     box,
		streams: model.NewDefaultStreams(true), // TODO tty
	}
}

func (cli *KubeBoxCli) Open() {
	cli.log.Debug().Msgf("init kube box:\n%v\n", cli.box.Template.Pretty())
	cli.loader.Start(fmt.Sprintf("loading %s", cli.box.Template.Name))
	// TODO remove ???
	cli.loader.Sleep(1)

	containerName := cli.box.Template.GenerateName()
	deployment, service, err := cli.box.BuildSpec(containerName)
	if err != nil {
		cli.loader.Halt(err, "error kube: invalid template")
	}

	cli.box.OnSetupCallback = func(message string) {
		cli.log.Debug().Msg(message)
	}
	cli.loader.Refresh(fmt.Sprintf("creating %s/%s", cli.box.ResourceOptions.Namespace, containerName))
	err = cli.box.ApplyTemplate(deployment, service)
	if err != nil {
		cli.loader.Halt(err, "error kube: apply template")
	}

	cli.box.OnCloseCallback = func(message string) {
		cli.log.Debug().Msg(message)
	}
	cli.box.OnCloseErrorCallback = func(err error, message string) {
		cli.log.Warn().Err(err).Msg(message)
	}
	defer cli.box.RemoveTemplate(deployment, service)

	pod, err := cli.box.GetPod(deployment)
	if err != nil {
		cli.loader.Halt(err, "error kube: invalid pod")
	}
	cli.log.Debug().Msgf("found matching pod %s", pod.Name)

	cli.log.Info().Msgf("opening new box: image=%s, namespace=%s, podName=%s", cli.box.Template.ImageName(), pod.Namespace, pod.Name)

	cli.box.OnTunnelCallback = func(port schema.PortV1) {
		cli.log.Info().Msgf("[%s][%s] forwarding %s (local) -> %s (remote)", pod.Name, port.Alias, port.Local, port.Remote)
	}
	cli.box.OnTunnelErrorCallback = func(err error, message string) {
		cli.loader.Halt(err, "error kube: port forward")
	}
	cli.box.PortForward(pod)

	cli.box.OnExecCallback = func() {
		cli.log.Debug().Msgf("exec into pod %s", pod.Name)
		cli.loader.Stop()
	}
	cli.box.Exec(pod, cli.streams)
	if err != nil {
		cli.loader.Halt(err, "error kube: exec pod")
	}
}
