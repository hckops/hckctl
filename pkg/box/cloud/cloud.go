package cloud

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/cloud"
)

func newCloudBox(internalOpts *model.BoxInternalOpts, cloudConfig *cloud.CloudClientConfig) (*CloudBox, error) {
	internalOpts.EventBus.Publish(newClientInitCloudEvent())

	cloudClient, err := cloud.NewCloudClient(cloudConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error cloud box")
	}

	return &CloudBox{
		client:       cloudClient,
		clientConfig: cloudConfig,
		streams:      internalOpts.Streams,
		eventBus:     internalOpts.EventBus,
	}, nil
}

func (box *CloudBox) close() error {
	box.eventBus.Publish(newClientCloseCloudEvent())
	box.eventBus.Close()
	return box.client.Close()
}
