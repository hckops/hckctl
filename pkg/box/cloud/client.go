package cloud

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/cloud"
	"github.com/hckops/hckctl/pkg/event"
)

type CloudBox struct {
	client       *cloud.CloudClient
	clientConfig *cloud.CloudClientConfig
	streams      *model.BoxStreams
	eventBus     *event.EventBus
}

func NewCloudBox(internalOpts *model.BoxInternalOpts, cloudConfig *cloud.CloudClientConfig) (*CloudBox, error) {
	return newCloudBox(internalOpts, cloudConfig)
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
