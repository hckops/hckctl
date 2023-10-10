package kubernetes

import (
	"github.com/hckops/hckctl/pkg/util"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"

	"github.com/hckops/hckctl/pkg/client/kubernetes"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/event"
)

type KubeCommonClient struct {
	client     *kubernetes.KubeClient
	clientOpts *commonModel.KubeOptions
	eventBus   *event.EventBus
}

func NewKubeCommonClient(kubeOpts *commonModel.KubeOptions, eventBus *event.EventBus) (*KubeCommonClient, error) {
	eventBus.Publish(newInitKubeClientEvent())

	kubeClient, err := kubernetes.NewKubeClient(kubeOpts.InCluster, kubeOpts.ConfigPath)
	if err != nil {
		return nil, errors.Wrap(err, "error kube common client")
	}

	return &KubeCommonClient{
		client:     kubeClient,
		clientOpts: kubeOpts,
		eventBus:   eventBus,
	}, nil
}

func (common *KubeCommonClient) GetClient() *kubernetes.KubeClient {
	return common.client
}

func (common *KubeCommonClient) Close() error {
	common.eventBus.Publish(newCloseKubeClientEvent())
	common.eventBus.Close()
	return common.client.Close()
}

func (common *KubeCommonClient) SidecarVpnDelete(namespace string, mainContainerName string) error {
	// delete secret
	name := buildSidecarVpnSecretName(mainContainerName)
	if ok, err := common.client.SecretDelete(namespace, name); err != nil {
		return err
	} else if ok {
		common.eventBus.Publish(newSecretDeleteKubeEvent(namespace, name))
	}
	return nil
}

func boolPtr(b bool) *bool { return &b }

func (common *KubeCommonClient) SidecarVpnInject(namespace string, opts *commonModel.SidecarVpnInjectOpts, podSpec *corev1.PodSpec) error {

	// create secret
	secret := buildSidecarVpnSecret(namespace, opts.MainContainerName, opts.VpnInfo.ConfigValue)
	common.eventBus.Publish(newSecretCreateKubeEvent(namespace, secret.Name))
	if err := common.client.SecretCreate(namespace, secret); err != nil {
		return err
	}

	// https://kubernetes.io/docs/tasks/configure-pod-container/share-process-namespace
	//podSpec.ShareProcessNamespace = boolPtr(true)

	// disable ipv6, see https://kubernetes.io/docs/tasks/administer-cluster/sysctl-cluster
	podSpec.SecurityContext = &corev1.PodSecurityContext{
		Sysctls: []corev1.Sysctl{
			{Name: "net.ipv6.conf.all.disable_ipv6", Value: "0"},
		},
	}

	// inject containers
	podSpec.Containers = append(
		// order matters
		[]corev1.Container{
			buildSidecarVpnContainer(),
			// add fake sleep to allow sidecar-vpn to connect properly before starting the main container
			{
				Name:  "sidecar-sleep",
				Image: "busybox",
				Lifecycle: &corev1.Lifecycle{
					PostStart: &corev1.LifecycleHandler{
						Exec: &corev1.ExecAction{
							Command: []string{"sleep", "1s"},
						},
					},
				},
			},
		},
		podSpec.Containers..., // current containers
	)

	// inject volumes
	podSpec.Volumes = append(
		podSpec.Volumes, // current volumes
		buildSidecarVpnVolumes(util.ToLowerKebabCase(opts.MainContainerName))..., // join slices
	)
	return nil
}
