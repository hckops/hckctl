package box

import (
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/event"
)

type BoxClient interface {
	Events() *event.EventBus
	Create(template *model.BoxV1) (*model.BoxInfo, error)
	Exec(name string, command string) error
	Copy(name string, from string, to string) error
	List() ([]model.BoxInfo, error)
	Open(template *model.BoxV1) error
	Tunnel(name string) error
	Delete(name string) error
}
