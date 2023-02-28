package box

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	model "github.com/hckops/hckctl/internal/model"
	"github.com/hckops/hckctl/internal/terminal"
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

	kubeconfig := filepath.Join(homedir.HomeDir(), strings.Replace(config.ConfigPath, "~/", "", 1))
	log.Debug().Msgf("read config: configPath=%s, kubeconfig=%s", config.ConfigPath, kubeconfig)

	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal().Err(err).Msg("error restConfig")
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("error clientSet")
	}

	// https://github.com/kubernetes/client-go/issues/1036
	namespace, err := clientSet.CoreV1().Namespaces().Apply(b.ctx, v1.Namespace(config.Namespace), metav1.ApplyOptions{FieldManager: "application/apply-patch"})
	if err != nil {
		log.Fatal().Err(err).Msg("error apply namespace")
	}
	log.Debug().Msgf("namespace %s successfully applied", namespace.Name)

	b.loader.Stop()
}
