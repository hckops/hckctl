package kubernetes

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"

	boxModel "github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/kubernetes"
	"github.com/hckops/hckctl/pkg/client/terminal"
	commonKube "github.com/hckops/hckctl/pkg/common/kubernetes"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/schema"
	"github.com/hckops/hckctl/pkg/util"
)

func newKubeBoxClient(commonOpts *boxModel.CommonBoxOptions, kubeOpts *commonModel.KubeOptions) (*KubeBoxClient, error) {

	kubeCommonClient, err := commonKube.NewKubeCommonClient(kubeOpts, commonOpts.EventBus)
	if err != nil {
		return nil, errors.Wrap(err, "error kube box client")
	}

	return &KubeBoxClient{
		client:     kubeCommonClient.GetClient(),
		clientOpts: kubeOpts,
		kubeCommon: kubeCommonClient,
		eventBus:   commonOpts.EventBus,
	}, nil
}

func (box *KubeBoxClient) close() error {
	return box.kubeCommon.Close()
}

func (box *KubeBoxClient) createBox(opts *boxModel.CreateOptions) (*boxModel.BoxInfo, error) {
	namespace := box.clientOpts.Namespace

	boxName := opts.Template.GenerateName()
	deployment, service, err := kubernetes.BuildResources(newResources(namespace, boxName, opts))
	if err != nil {
		return nil, err
	}

	// create namespace
	if err := box.client.NamespaceApply(namespace); err != nil {
		return nil, err
	}
	box.eventBus.Publish(newNamespaceApplyKubeEvent(namespace))

	// create service
	if opts.Template.HasPorts() {
		if err := box.client.ServiceCreate(namespace, service); err != nil {
			return nil, err
		}
		box.eventBus.Publish(newServiceCreateKubeEvent(namespace, service.Name))
	} else {
		box.eventBus.Publish(newServiceCreateIgnoreKubeEvent(namespace, service.Name))
	}

	// inject sidecar-volume
	if opts.CommonInfo.ShareDir != nil {
		sidecarOpts := &commonModel.SidecarShareInjectOpts{
			MainContainerName: opts.Template.MainContainerName(),
			ShareDir:          opts.CommonInfo.ShareDir,
		}
		if err := box.kubeCommon.SidecarShareInject(sidecarOpts, &deployment.Spec.Template.Spec); err != nil {
			return nil, err
		}
	}

	// create secret and inject sidecar-vpn
	if opts.CommonInfo.NetworkVpn != nil {
		sidecarOpts := &commonModel.SidecarVpnInjectOpts{
			Name:       boxName,
			NetworkVpn: opts.CommonInfo.NetworkVpn,
		}
		if err := box.kubeCommon.SidecarVpnInject(namespace, sidecarOpts, &deployment.Spec.Template.Spec); err != nil {
			return nil, err
		}
	}

	// create deployment
	box.eventBus.Publish(newResourcesDeployKubeLoaderEvent(namespace, opts.Template.Name))
	deploymentOpts := &kubernetes.DeploymentCreateOpts{
		Namespace: namespace,
		Spec:      deployment,
		OnStatusEventCallback: func(event string) {
			box.eventBus.Publish(newDeploymentCreateStatusKubeEvent(event))
		},
	}
	if err := box.client.DeploymentCreate(deploymentOpts); err != nil {
		return nil, err
	}
	box.eventBus.Publish(newDeploymentCreateKubeEvent(namespace, deployment.Name))

	podInfo, err := box.client.PodDescribeFromDeployment(deployment)
	if err != nil {
		return nil, err
	}
	box.eventBus.Publish(newPodNameKubeEvent(namespace, podInfo.PodName, podInfo.ContainerName))

	// upload shared directory
	if opts.CommonInfo.ShareDir != nil {
		sidecarOpts := &commonModel.SidecarShareUploadOpts{
			Namespace: namespace,
			PodName:   podInfo.PodName,
			ShareDir:  opts.CommonInfo.ShareDir,
		}
		if err := box.kubeCommon.SidecarShareUpload(sidecarOpts); err != nil {
			return nil, err
		}
	}

	// TODO always healthy unused? otherwise use DeploymentDescribe instead of PodDescribe
	return &boxModel.BoxInfo{Id: podInfo.PodName, Name: boxName, Healthy: true}, nil
}

func newResources(namespace string, name string, opts *boxModel.CreateOptions) *kubernetes.ResourcesOpts {

	var envs []kubernetes.KubeEnv
	for _, e := range opts.Template.EnvironmentVariables() {
		envs = append(envs, kubernetes.KubeEnv{Key: e.Key, Value: e.Value})
	}
	var ports []kubernetes.KubePort
	for _, p := range opts.Template.NetworkPortValues(false) {
		ports = append(ports, kubernetes.KubePort{Name: p.Alias, Port: p.Remote})
	}

	return &kubernetes.ResourcesOpts{
		Namespace:   namespace,
		Name:        name,
		Annotations: opts.Labels,
		Labels: kubernetes.BuildLabels(name, opts.Template.Image.Repository, opts.Template.Image.ResolveVersion(),
			map[string]string{commonModel.LabelSchemaKind: util.ToLowerKebabCase(schema.KindBoxV1.String())}),
		Ports: ports,
		PodInfo: &kubernetes.PodInfo{
			Namespace:     namespace,
			PodName:       "INVALID_POD_NAME", // not used, generated suffix by kube
			ContainerName: opts.Template.Image.Repository,
			ImageName:     opts.Template.Image.Name(),
			Arguments:     []string{},
			Env:           envs,
			Resource:      opts.Size.ToKubeResource(),
		},
	}
}

func (box *KubeBoxClient) connectBox(opts *boxModel.ConnectOptions) error {
	if info, err := box.searchBox(opts.Name); err != nil {
		return err
	} else {
		if opts.DisableExec && opts.DisableTunnel {
			return errors.New("invalid connection options")
		}

		namespace := box.clientOpts.Namespace
		for _, e := range opts.Template.EnvironmentVariables() {
			box.eventBus.Publish(newPodEnvKubeEvent(namespace, info.Id, e))
			box.eventBus.Publish(newPodEnvKubeConsoleEvent(namespace, info.Name, e))
		}

		// tunnel only
		if opts.DisableExec {
			// tunnel and block to exit, wait until killed
			return box.podPortForward(opts.Template, info, true)
		}

		if !opts.DisableTunnel {
			// tunnel and exec after, do not block
			if err := box.podPortForward(opts.Template, info, false); err != nil {
				return err
			}
		}

		return box.execBox(opts, info)
	}
}

func boxNameLabelSelector(name string) string {
	return fmt.Sprintf("%s,%s=%s", boxModel.BoxLabelSelector(), kubernetes.LabelKubeName, name)
}

func (box *KubeBoxClient) searchBox(name string) (*boxModel.BoxInfo, error) {
	namespace := box.clientOpts.Namespace
	box.eventBus.Publish(newDeploymentSearchKubeEvent(namespace, name))

	deployments, err := box.client.DeploymentList(namespace, boxModel.BoxPrefixName, boxNameLabelSelector(name))
	if err != nil {
		return nil, err
	}

	switch len(deployments) {
	case 0:
		return nil, errors.New("box not found")
	case 1:
		info := newBoxInfo(deployments[0])
		return &info, nil
	default:
		// this should never happen
		return nil, errors.New("unexpected label selector match")
	}
}

func (box *KubeBoxClient) execBox(opts *boxModel.ConnectOptions, info *boxModel.BoxInfo) error {

	// TODO if BoxInfo not Healthy attempt scale 1

	if opts.Template.Shell == boxModel.BoxShellNone {
		// stop loader
		box.eventBus.Publish(newPodExecKubeLoaderEvent())

		return box.logsBox(opts, info)
	}

	if opts.DeleteOnExit {
		defer box.deleteBox(info.Name)
	}

	// exec
	execOpts := &kubernetes.PodExecOpts{
		Namespace:     box.clientOpts.Namespace,
		PodName:       info.Id,
		ContainerName: opts.Template.MainContainerName(),
		Commands:      terminal.DefaultShellCommand(opts.Template.Shell),
		InStream:      opts.StreamOpts.In,
		OutStream:     opts.StreamOpts.Out,
		ErrStream:     opts.StreamOpts.Err,
		IsTty:         opts.StreamOpts.IsTty,
		OnExecCallback: func() {
			// stop loader
			box.eventBus.Publish(newPodExecKubeLoaderEvent())
		},
	}
	box.eventBus.Publish(newPodExecKubeEvent(opts.Template.Name, box.clientOpts.Namespace, info.Id, opts.Template.Shell))
	return box.client.PodExecShell(execOpts)
}

func (box *KubeBoxClient) logsBox(opts *boxModel.ConnectOptions, info *boxModel.BoxInfo) error {
	namespace := box.clientOpts.Namespace

	if opts.DeleteOnExit {
		opts.OnInterruptCallback(func() {
			box.eventBus.Publish(newPodLogsExitKubeEvent(namespace, info.Id))
			box.eventBus.Publish(newPodLogsExitKubeConsoleEvent())
			box.deleteBox(info.Name)
		})
	}

	logsOpts := &kubernetes.PodLogsOpts{
		Namespace:     namespace,
		PodName:       info.Id,
		ContainerName: opts.Template.MainContainerName(),
		OutStream:     opts.StreamOpts.Out,
	}
	box.eventBus.Publish(newPodLogsKubeEvent(namespace, info.Id))
	return box.client.PodLogs(logsOpts)
}

func (box *KubeBoxClient) podPortForward(template *boxModel.BoxV1, boxInfo *boxModel.BoxInfo, isWait bool) error {
	namespace := box.clientOpts.Namespace

	if !template.HasPorts() {
		box.eventBus.Publish(newPodPortForwardIgnoreKubeEvent(namespace, boxInfo.Id))
		// exit, no service/port available to bind
		return nil
	}

	networkPorts := template.NetworkPortValues(false)
	portPadding := boxModel.PortFormatPadding(networkPorts)
	ports, err := ToPortBindings(networkPorts, func(port boxModel.BoxPort) {
		box.eventBus.Publish(newPodPortForwardBindingKubeEvent(namespace, boxInfo.Id, port))
		box.eventBus.Publish(newPodPortForwardBindingKubeConsoleEvent(namespace, boxInfo.Name, port, portPadding))
	})
	if err != nil {
		return err
	}

	opts := &kubernetes.PodPortForwardOpts{
		Namespace: namespace,
		PodName:   boxInfo.Id,
		Ports:     ports,
		IsWait:    isWait,
		OnTunnelStartCallback: func() {
			// stop loader
			box.eventBus.Publish(newPodExecKubeLoaderEvent())
		},
		OnTunnelErrorCallback: func(err error) {
			box.eventBus.Publish(newPodPortForwardErrorKubeEvent(namespace, boxInfo.Id, err))
		},
	}
	if err := box.client.PodPortForward(opts); err != nil {
		return err
	}

	return nil
}

func ToPortBindings(ports []boxModel.BoxPort, onPortBindCallback func(port boxModel.BoxPort)) ([]string, error) {

	var portBindings []string
	for _, port := range ports {

		localPort, err := util.FindOpenPort(port.Local)
		if err != nil {
			return nil, errors.Wrapf(err, "error kube local port %s", port.Local)
		}

		// actual bound port
		onPortBindCallback(boxModel.BoxPort{
			Alias:  port.Alias,
			Local:  localPort,
			Remote: port.Remote,
		})

		portBindings = append(portBindings, fmt.Sprintf("%s:%s", localPort, port.Remote))
	}
	return portBindings, nil
}

func (box *KubeBoxClient) describeBox(name string) (*boxModel.BoxDetails, error) {
	namespace := box.clientOpts.Namespace

	boxInfo, err := box.searchBox(name)
	if err != nil {
		return nil, err
	}

	box.eventBus.Publish(newDeploymentDescribeKubeEvent(namespace, name))
	deployment, err := box.client.DeploymentDescribe(namespace, boxInfo.Name)
	if err != nil {
		return nil, err
	}

	box.eventBus.Publish(newServiceDescribeKubeEvent(namespace, name))
	service, err := box.client.ServiceDescribe(namespace, boxInfo.Name)
	if err != nil {
		return nil, err
	}

	return ToBoxDetails(deployment, service, box.Provider())
}

func ToBoxDetails(deployment *kubernetes.DeploymentDetails, serviceInfo *kubernetes.ServiceInfo, provider boxModel.BoxProvider) (*boxModel.BoxDetails, error) {

	labels := commonModel.Labels(deployment.Annotations)

	size, err := boxModel.ToBoxSize(labels)
	if err != nil {
		return nil, err
	}

	var envs []boxModel.BoxEnv
	for _, env := range deployment.Info.PodInfo.Env {
		envs = append(envs, boxModel.BoxEnv{
			Key:   env.Key,
			Value: env.Value,
		})
	}

	var ports []boxModel.BoxPort
	for _, p := range serviceInfo.Ports {
		ports = append(ports, boxModel.BoxPort{
			Alias:  p.Name,
			Local:  boxModel.BoxPortNone, // runtime only
			Remote: p.Port,
			Public: false,
		})
	}

	return &boxModel.BoxDetails{
		Info: newBoxInfo(*deployment.Info),
		TemplateInfo: &boxModel.BoxTemplateInfo{
			CachedTemplate: labels.ToCachedTemplateInfo(),
			GitTemplate:    labels.ToGitTemplateInfo(),
		},
		ProviderInfo: &boxModel.BoxProviderInfo{
			Provider: provider,
			KubeProvider: &commonModel.KubeProviderInfo{
				Namespace: deployment.Info.Namespace,
			},
		},
		Size:    size,
		Env:     boxModel.SortEnv(envs),
		Ports:   boxModel.SortPorts(ports),
		Created: deployment.Created,
	}, nil
}

func (box *KubeBoxClient) listBoxes() ([]boxModel.BoxInfo, error) {
	namespace := box.clientOpts.Namespace

	deployments, err := box.client.DeploymentList(namespace, boxModel.BoxPrefixName, boxModel.BoxLabelSelector())
	if err != nil {
		return nil, err
	}
	var result []boxModel.BoxInfo
	for index, d := range deployments {
		result = append(result, newBoxInfo(d))
		box.eventBus.Publish(newDeploymentListKubeEvent(index, namespace, d.Name, d.Healthy))
	}

	return result, nil
}

func newBoxInfo(deployment kubernetes.DeploymentInfo) boxModel.BoxInfo {
	return boxModel.BoxInfo{
		Id:      deployment.PodInfo.PodName,
		Name:    deployment.Name,
		Healthy: deployment.Healthy,
	}
}

func (box *KubeBoxClient) deleteBoxes(names []string) ([]string, error) {
	namespace := box.clientOpts.Namespace

	// optimize delete
	if len(names) == 1 {
		boxInfo, err := box.searchBox(names[0])
		if err != nil {
			return nil, err
		}
		return []string{boxInfo.Name}, box.deleteBox(boxInfo.Name)
	}

	boxes, err := box.listBoxes()
	if err != nil {
		return nil, err
	}

	var deleted []string
	for _, boxInfo := range boxes {

		// all or filter
		if len(names) == 0 || slices.Contains(names, boxInfo.Name) {

			if err := box.deleteBox(boxInfo.Name); err == nil {
				deleted = append(deleted, boxInfo.Name)
			} else {
				// silently ignore
				box.eventBus.Publish(newResourcesDeleteIgnoreKubeEvent(namespace, boxInfo.Name))
			}
		}
	}
	return deleted, nil
}

func (box *KubeBoxClient) deleteBox(name string) error {
	namespace := box.clientOpts.Namespace

	box.eventBus.Publish(newDeploymentDeleteKubeEvent(namespace, name))
	if err := box.client.DeploymentDelete(namespace, name); err != nil {
		return err
	}

	box.eventBus.Publish(newServiceDeleteKubeEvent(namespace, name))
	if err := box.client.ServiceDelete(namespace, name); err != nil {
		return err
	}

	if err := box.kubeCommon.SidecarVpnDelete(namespace, name); err != nil {
		return err
	}
	return nil
}

func (box *KubeBoxClient) clean() error {
	namespace := box.clientOpts.Namespace

	box.eventBus.Publish(newNamespaceDeleteKubeEvent(namespace))
	return box.client.NamespaceDelete(namespace)
}
