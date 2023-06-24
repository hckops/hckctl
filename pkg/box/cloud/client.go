package cloud

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/ssh"
	"github.com/hckops/hckctl/pkg/event"
)

type CloudBox struct {
	clientVersion string
	clientConfig  *ssh.SshClientConfig
	client        *ssh.SshClient
	streams       *model.BoxStreams
	eventBus      *event.EventBus
}

func NewCloudBox(internalOpts *model.BoxInternalOpts, sshConfig *ssh.SshClientConfig) (*CloudBox, error) {
	return newCloudBox(internalOpts, sshConfig)
}

func (box *CloudBox) Provider() model.BoxProvider {
	return model.Cloud
}

func (box *CloudBox) Events() *event.EventBus {
	return box.eventBus
}

func (box *CloudBox) Create(template *model.BoxV1) (*model.BoxInfo, error) {
	defer box.close()
	return nil, errors.New("not implemented")
}

func (box *CloudBox) Exec(template *model.BoxV1, name string) error {
	defer box.close()
	return errors.New("not implemented")
}

func (box *CloudBox) Open(template *model.BoxV1) error {
	defer box.close()
	return errors.New("not implemented")
}

func (box *CloudBox) List() ([]model.BoxInfo, error) {
	defer box.close()
	return nil, errors.New("not implemented")
}

func (box *CloudBox) Copy(string, string, string) error {
	defer box.close()
	return errors.New("not implemented")
}

func (box *CloudBox) Tunnel(string) error {
	defer box.close()
	return errors.New("not implemented")
}

func (box *CloudBox) Delete(name string) error {
	defer box.close()
	return errors.New("not implemented")
}

func (box *CloudBox) DeleteAll() ([]model.BoxInfo, error) {
	defer box.close()
	return nil, errors.New("not implemented")
}
