package kubernetes

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/common"
	"github.com/hckops/hckctl/pkg/client/kubernetes"
	"github.com/hckops/hckctl/pkg/schema"
	"github.com/hckops/hckctl/pkg/util"
)

func newKubeBoxClient(commonOpts *model.CommonBoxOptions, kubeOpts *model.KubeBoxOptions) (*KubeBoxClient, error) {
	commonOpts.EventBus.Publish(newInitKubeClientEvent())

	kubeClient, err := kubernetes.NewKubeClient(kubeOpts.InCluster, kubeOpts.ConfigPath)
	if err != nil {
		return nil, errors.Wrap(err, "error kube box")
	}

	return &KubeBoxClient{
		client:     kubeClient,
		clientOpts: kubeOpts,
		eventBus:   commonOpts.EventBus,
	}, nil
}

func (box *KubeBoxClient) close() error {
	box.eventBus.Publish(newCloseKubeClientEvent())
	box.eventBus.Close()
	return box.client.Close()
}

func (box *KubeBoxClient) createBox(opts *model.CreateOptions) (*model.BoxInfo, error) {
	namespace := box.clientOpts.Namespace

	// TODO add env var container override
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

	podInfo, err := box.client.PodDescribe(deployment)
	if err != nil {
		return nil, err
	}
	box.eventBus.Publish(newPodNameKubeEvent(namespace, podInfo.PodName, podInfo.ContainerName))

	// TODO always healthy unused? otherwise use DeploymentDescribe instead of PodDescribe
	return &model.BoxInfo{Id: podInfo.PodName, Name: boxName, Healthy: true}, nil
}

func newResources(namespace string, name string, opts *model.CreateOptions) *kubernetes.ResourcesOpts {

	var ports []kubernetes.KubePort
	for _, p := range opts.Template.NetworkPortValues(false) {
		ports = append(ports, kubernetes.KubePort{Name: p.Alias, Port: p.Remote})
	}

	return &kubernetes.ResourcesOpts{
		Namespace:   namespace,
		Name:        name,
		Annotations: opts.Labels,
		Labels: kubernetes.BuildLabels(name, opts.Template.Image.Repository, opts.Template.ImageVersion(),
			map[string]string{model.LabelSchemaKind: common.ToKebabCase(schema.KindBoxV1.String())}),
		Ports: ports,
		PodInfo: &kubernetes.PodInfo{
			Namespace:     namespace,
			PodName:       "INVALID_POD_NAME", // not used, generated suffix by kube
			ContainerName: opts.Template.Image.Repository,
			ImageName:     opts.Template.ImageName(),
			Env:           nil, // TODO not used
			Resource:      opts.Size.ToKubeResource(),
		},
	}
}

func (box *KubeBoxClient) connectBox(opts *model.ConnectOptions) error {
	if info, err := box.searchBox(opts.Name); err != nil {
		return err
	} else {
		if opts.DisableExec && opts.DisableTunnel {
			return errors.New("invalid connection options")
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

		return box.execBox(opts.Template, info, opts.Streams, opts.DeleteOnExit)
	}
}

func boxLabelSelector() string {
	// value must be sanitized
	return fmt.Sprintf("%s=%s", model.LabelSchemaKind, common.ToKebabCase(schema.KindBoxV1.String()))
}

func boxNameLabelSelector(name string) string {
	return fmt.Sprintf("%s,%s=%s", boxLabelSelector(), kubernetes.LabelKubeName, name)
}

func (box *KubeBoxClient) searchBox(name string) (*model.BoxInfo, error) {
	namespace := box.clientOpts.Namespace
	box.eventBus.Publish(newDeploymentSearchKubeEvent(namespace, name))

	deployments, err := box.client.DeploymentList(namespace, model.BoxPrefixName, boxNameLabelSelector(name))
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

func (box *KubeBoxClient) execBox(template *model.BoxV1, info *model.BoxInfo, streams *model.BoxStreams, deleteOnExit bool) error {
	box.eventBus.Publish(newPodExecKubeEvent(template.Name, box.clientOpts.Namespace, info.Id, template.Shell))

	// TODO if BoxInfo not Healthy attempt scale 1
	// TODO model.BoxShellNone
	// TODO print environment variables

	// exec
	opts := &kubernetes.PodExecOpts{
		Namespace: box.clientOpts.Namespace,
		PodName:   common.ToKebabCase(template.Image.Repository), // pod.Spec.Containers[0].Name
		PodId:     info.Id,
		Shell:     template.Shell,
		InStream:  streams.In,
		OutStream: streams.Out,
		ErrStream: streams.Err,
		IsTty:     streams.IsTty,
		OnExecCallback: func() {
			// stop loader
			box.eventBus.Publish(newPodExecKubeLoaderEvent())
		},
	}

	if deleteOnExit {
		defer box.deleteBox(info.Name)
	}

	return box.client.PodExec(opts)
}

func (box *KubeBoxClient) podPortForward(template *model.BoxV1, boxInfo *model.BoxInfo, isWait bool) error {
	namespace := box.clientOpts.Namespace

	if !template.HasPorts() {
		box.eventBus.Publish(newPodPortForwardIgnoreKubeEvent(namespace, boxInfo.Id))
		// exit, no service/port available to bind
		return nil
	}

	networkPorts := template.NetworkPortValues(false)
	portPadding := model.PortFormatPadding(networkPorts)
	ports, err := toPortBindings(networkPorts, func(port model.BoxPort) {
		box.eventBus.Publish(newPodPortForwardBindingKubeEvent(namespace, boxInfo.Id, port))
		box.eventBus.Publish(newPodPortForwardBindingKubeConsoleEvent(namespace, boxInfo.Name, port, portPadding))
	})
	if err != nil {
		return err
	}

	opts := &kubernetes.PodPortForwardOpts{
		Namespace: namespace,
		PodId:     boxInfo.Id,
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

func toPortBindings(ports []model.BoxPort, onPortBindCallback func(port model.BoxPort)) ([]string, error) {

	var portBindings []string
	for _, port := range ports {

		localPort, err := util.FindOpenPort(port.Local)
		if err != nil {
			return nil, errors.Wrapf(err, "error kube local port %s", port.Local)
		}

		// actual bound port
		onPortBindCallback(model.BoxPort{
			Alias:  port.Alias,
			Local:  localPort,
			Remote: port.Remote,
		})

		portBindings = append(portBindings, fmt.Sprintf("%s:%s", localPort, port.Remote))
	}
	return portBindings, nil
}

func (box *KubeBoxClient) describe(name string) (*model.BoxDetails, error) {
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

	return toBoxDetails(deployment, service)
}

func toBoxDetails(deployment *kubernetes.DeploymentDetails, serviceInfo *kubernetes.ServiceInfo) (*model.BoxDetails, error) {

	labels := model.BoxLabels(deployment.Annotations)

	size, err := labels.ToSize()
	if err != nil {
		return nil, err
	}

	var env []model.BoxEnv
	for key, value := range deployment.Info.PodInfo.Env {
		env = append(env, model.BoxEnv{
			Key:   key,
			Value: value,
		})
	}

	var ports []model.BoxPort
	for _, p := range serviceInfo.Ports {
		ports = append(ports, model.BoxPort{
			Alias:  p.Name,
			Local:  model.BoxPortNone, // runtime only
			Remote: p.Port,
			Public: false,
		})
	}

	return &model.BoxDetails{
		Info: newBoxInfo(*deployment.Info),
		TemplateInfo: &model.BoxTemplateInfo{
			CachedTemplate: labels.ToCachedTemplateInfo(),
			GitTemplate:    labels.ToGitTemplateInfo(),
		},
		ProviderInfo: &model.BoxProviderInfo{
			Provider: model.Kubernetes,
			KubeProvider: &model.KubeProviderInfo{
				Namespace: deployment.Info.Namespace,
			},
		},
		Size:    size,
		Env:     model.SortEnv(env),
		Ports:   model.SortPorts(ports),
		Created: deployment.Created,
	}, nil
}

func (box *KubeBoxClient) listBoxes() ([]model.BoxInfo, error) {
	namespace := box.clientOpts.Namespace

	deployments, err := box.client.DeploymentList(namespace, model.BoxPrefixName, boxLabelSelector())
	if err != nil {
		return nil, err
	}
	var result []model.BoxInfo
	for index, d := range deployments {
		result = append(result, newBoxInfo(d))
		box.eventBus.Publish(newDeploymentListKubeEvent(index, namespace, d.Name, d.PodInfo.PodName, d.Healthy))
	}

	return result, nil
}

func newBoxInfo(deployment kubernetes.DeploymentInfo) model.BoxInfo {
	return model.BoxInfo{
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

	return nil
}

func (box *KubeBoxClient) clean() error {
	namespace := box.clientOpts.Namespace

	box.eventBus.Publish(newNamespaceDeleteKubeEvent(namespace))
	return box.client.NamespaceDelete(namespace)
}
