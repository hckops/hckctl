package kubernetes

import (
	"github.com/pkg/errors"

	boxModel "github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/kubernetes"
	commonKube "github.com/hckops/hckctl/pkg/common/kubernetes"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/event"
)

type KubeBoxClient struct {
	client     *kubernetes.KubeClient
	clientOpts *commonModel.KubeOptions
	kubeCommon *commonKube.KubeCommonClient
	eventBus   *event.EventBus
}

func NewKubeBoxClient(commonOpts *boxModel.CommonBoxOptions, kubeOpts *commonModel.KubeOptions) (*KubeBoxClient, error) {
	return newKubeBoxClient(commonOpts, kubeOpts)
}

func (box *KubeBoxClient) Provider() boxModel.BoxProvider {
	return boxModel.Kubernetes
}

func (box *KubeBoxClient) Events() *event.EventBus {
	return box.eventBus
}

func (box *KubeBoxClient) Create(opts *boxModel.CreateOptions) (*boxModel.BoxInfo, error) {
	defer box.close()
	return box.createBox(opts)
}

func (box *KubeBoxClient) Connect(opts *boxModel.ConnectOptions) error {
	defer box.close()
	return box.connectBox(opts)
}

func (box *KubeBoxClient) Copy(string, string, string) error {
	defer box.close()
	return errors.New("not implemented")
}

func (box *KubeBoxClient) Describe(name string) (*boxModel.BoxDetails, error) {
	defer box.close()
	return box.describeBox(name)
}

func (box *KubeBoxClient) List() ([]boxModel.BoxInfo, error) {
	defer box.close()
	return box.listBoxes()
}

func (box *KubeBoxClient) Delete(names []string) ([]string, error) {
	defer box.close()
	return box.deleteBoxes(names)
}

func (box *KubeBoxClient) Clean() error {
	defer box.close()
	return box.clean()
}

func (box *KubeBoxClient) Version() (string, error) {
	return "", errors.New("not implemented")
}
