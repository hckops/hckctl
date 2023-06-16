package kubernetes

import (
	"github.com/pkg/errors"
)

func NewKubeClient() (*KubeClient, error) {
	return &KubeClient{}, nil
}

func (client *KubeClient) Close() error {
	return errors.New("not implemented")
}
