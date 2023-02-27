package box

import (
	"context"
	"fmt"

	model "github.com/hckops/hckctl/internal/model"
	"github.com/hckops/hckctl/internal/terminal"
	"github.com/rs/zerolog/log"
)

type KubeBox struct {
	ctx         context.Context
	loader      *terminal.Loader
	boxTemplate *model.BoxV1
}

// TODO convert box to kube, create vs apply (e.g. no failures for existing namespace), exec + forward
func NewKubeBox(box *model.BoxV1) *KubeBox {
	return &KubeBox{
		ctx:         context.Background(),
		loader:      terminal.NewLoader(),
		boxTemplate: box,
	}
}

func (b *KubeBox) InitBox(config model.KubeConfig) {
	log.Debug().Msgf("init kube box: \n%v\n", b.boxTemplate.Pretty())
	b.loader.Start(fmt.Sprintf("loading %s", b.boxTemplate.Name))

	b.loader.Sleep(2)

	// config.ConfigPath
	// config.Namespace

	b.loader.Stop()
}
