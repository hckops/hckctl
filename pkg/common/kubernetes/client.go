package kubernetes

import (
	"io"
	"os"

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

func (common *KubeCommonClient) SidecarVpnInject(namespace string, opts *commonModel.SidecarVpnInjectOpts, podSpec *corev1.PodSpec) error {

	// create secret
	secret := buildSidecarVpnSecret(namespace, opts.MainContainerId, opts.NetworkVpn.ConfigValue)
	common.eventBus.Publish(newSecretCreateKubeEvent(namespace, secret.Name))
	if err := common.client.SecretCreate(namespace, secret); err != nil {
		return err
	}

	// update pod
	injectSidecarVpn(podSpec, opts.MainContainerId)
	common.eventBus.Publish(newSidecarVpnConnectKubeEvent(opts.NetworkVpn.Name))

	return nil
}

func (common *KubeCommonClient) SidecarShareInject(opts *commonModel.SidecarShareInjectOpts, podSpec *corev1.PodSpec) error {
	// update pod
	injectSidecarShare(podSpec, opts.MainContainerName, opts.ShareDir)
	common.eventBus.Publish(newSidecarShareMountKubeEvent(opts.ShareDir.RemotePath))
	return nil
}

func (common *KubeCommonClient) SidecarShareUpload(opts *commonModel.SidecarShareUploadOpts) error {
	common.eventBus.Publish(newSidecarShareUploadKubeEvent(opts.ShareDir.LocalPath, opts.ShareDir.RemotePath))
	common.eventBus.Publish(newSidecarShareUploadKubeLoaderEvent())

	sidecarContainerName := buildSidecarShareContainerName()
	copyOpts := &kubernetes.CopyPodOpts{
		Namespace:     opts.Namespace,
		PodName:       opts.PodName,
		ContainerName: sidecarContainerName,
		LocalPath:     opts.ShareDir.LocalPath,
		RemotePath:    opts.ShareDir.RemotePath,
	}
	if err := common.client.CopyToPod(copyOpts); err != nil {
		return err
	}

	if opts.ShareDir.LockDir {
		// remove lock
		execDeleteLock := &kubernetes.PodExecOpts{
			Namespace:      opts.Namespace,
			PodId:          opts.PodName,
			PodName:        sidecarContainerName,
			Commands:       []string{"rm", buildSidecarShareLock(opts.ShareDir.RemotePath)},
			InStream:       os.Stdin,
			OutStream:      io.Discard,
			ErrStream:      io.Discard,
			IsTty:          false,
			OnExecCallback: func() {},
		}
		if err := common.client.PodExecCommand(execDeleteLock); err != nil {
			return errors.Wrapf(err, "error delete lock")
		}
	}

	return nil
}
