package cloud

import (
	"github.com/hckops/hckctl/pkg/client/ssh"
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/model"
)

func newCloudBox(internalOpts *model.BoxInternalOpts, sshConfig *ssh.SshClientConfig) (*CloudBox, error) {
	internalOpts.EventBus.Publish(newClientInitCloudEvent())

	sshClient, err := ssh.NewSshClient(sshConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error cloud box")
	}

	return &CloudBox{
		client:       sshClient,
		clientConfig: sshConfig,
		streams:      internalOpts.Streams,
		eventBus:     internalOpts.EventBus,
	}, nil
}

func (box *CloudBox) close() error {
	box.eventBus.Publish(newClientCloseCloudEvent())
	box.eventBus.Close()
	return box.client.Close()
}
