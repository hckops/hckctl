package box

import (
	"context"

	"github.com/hckops/hckctl/pkg/template/model"
)

type Connection struct {
	Out chan string // TODO
}

type Client interface {
	Open() (*Connection, error)
}

// ******************************

type DockerClientOpts struct {
	ctx      context.Context
	template model.BoxV1
}

type DockerClient struct {
	opts *DockerClientOpts
}

func (client *DockerClient) Open() (*Connection, error) {
	return nil, nil
}

// ******************************

type KubeClientOpts struct {
	ctx      context.Context
	template model.BoxV1
}

type KubeClient struct {
	opts *KubeClientOpts
}

func (client *KubeClient) Open() (*Connection, error) {
	return nil, nil
}
