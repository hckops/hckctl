package lab

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/event"
	"github.com/hckops/hckctl/pkg/lab/cloud"
	"github.com/hckops/hckctl/pkg/lab/model"
)

type LabClient interface {
	Provider() model.LabProvider
	Events() *event.EventBus
	Create(opts *model.CreateOptions) (*model.LabInfo, error)
}

// TODO generics Box/Lab
func NewLabClient(opts *model.LabClientOptions) (LabClient, error) {
	commonOpts := model.NewCommonLabOpts()
	switch opts.Provider {
	case model.Cloud:
		return cloud.NewCloudLabClient(commonOpts, opts.CloudOpts)
	default:
		return nil, errors.New("invalid provider")
	}
}
