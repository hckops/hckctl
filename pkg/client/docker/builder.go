package docker

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"

	"github.com/hckops/hckctl/pkg/util"
)

func BuildContainerConfig(opts *ContainerConfigOpts) (*container.Config, error) {

	var envs []string
	for _, env := range opts.Env {
		envs = append(envs, fmt.Sprintf("%s=%s", env.Key, env.Value))
	}
	exposedPorts := make(nat.PortSet)
	for _, port := range opts.Ports {
		p, err := nat.NewPort("tcp", port.Remote)
		if err != nil {
			return nil, errors.Wrap(err, "error docker port: containerConfig")
		}
		exposedPorts[p] = struct{}{}
	}

	return &container.Config{
		Hostname:     opts.Hostname,
		Image:        opts.ImageName,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		StdinOnce:    true,
		Tty:          opts.Tty,
		Cmd:          opts.Cmd,
		Env:          envs,
		ExposedPorts: exposedPorts,
		Labels:       opts.Labels,
	}, nil
}

func BuildHostConfig(opts *ContainerHostConfigOpts) (*container.HostConfig, error) {

	portBindings := make(nat.PortMap)
	for _, port := range opts.Ports {

		localPort, err := util.FindOpenPort(port.Local)
		if err != nil {
			return nil, errors.Wrap(err, "error docker local port: hostConfig")
		}

		remotePort, err := nat.NewPort("tcp", port.Remote)
		if err != nil {
			return nil, errors.Wrap(err, "error docker remote port: hostConfig")
		}

		// actual bound port
		opts.OnPortBindCallback(ContainerPort{
			Local:  localPort,
			Remote: port.Remote,
		})

		portBindings[remotePort] = []nat.PortBinding{{
			HostIP:   "0.0.0.0",
			HostPort: localPort,
		}}
	}

	return &container.HostConfig{
		NetworkMode:  container.NetworkMode(opts.NetworkMode),
		PortBindings: portBindings,
	}, nil
}

func DefaultNetworkMode() string {
	return string(container.IsolationDefault)
}

func ContainerNetworkMode(idOrName string) string {
	return strings.Join([]string{"container", idOrName}, ":")
}

func BuildNetworkingConfig(networkName, networkId string) *network.NetworkingConfig {
	return &network.NetworkingConfig{EndpointsConfig: map[string]*network.EndpointSettings{networkName: {NetworkID: networkId}}}
}
