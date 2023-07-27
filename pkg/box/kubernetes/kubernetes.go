package kubernetes

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/common"
	"github.com/hckops/hckctl/pkg/client/kubernetes"
	"github.com/hckops/hckctl/pkg/schema"
	"github.com/hckops/hckctl/pkg/util"
)

func newKubeBoxClient(commonOpts *model.CommonBoxOptions, kubeOpts *model.KubeBoxOptions) (*KubeBoxClient, error) {
	commonOpts.EventBus.Publish(newClientInitKubeEvent())

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
	box.eventBus.Publish(newClientCloseKubeEvent())
	box.eventBus.Close()
	return box.client.Close()
}

func (box *KubeBoxClient) createBox(opts *model.TemplateOptions) (*model.BoxInfo, error) {
	namespace := box.clientOpts.Namespace

	// TODO add env var container override

	boxName := opts.Template.GenerateName()
	deployment, service, err := buildSpec(boxName, namespace, opts)
	if err != nil {
		return nil, err
	}

	if err := box.client.NamespaceApply(namespace); err != nil {
		return nil, err
	}
	box.eventBus.Publish(newNamespaceApplyKubeEvent(namespace))

	if opts.Template.HasPorts() {
		if err := box.client.ServiceCreate(namespace, service); err != nil {
			return nil, err
		}
		box.eventBus.Publish(newServiceCreateKubeEvent(namespace, service.Name))
	} else {
		box.eventBus.Publish(newServiceCreateSkippedKubeEvent(namespace, service.Name))
	}

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

	return &model.BoxInfo{Id: podInfo.PodName, Name: boxName, Healthy: true}, nil
}

func buildSpec(name string, namespace string, templateOpts *model.TemplateOptions) (*appsv1.Deployment, *corev1.Service, error) {

	labels := buildLabels(name, templateOpts.Template.Image.Repository, templateOpts.Template.ImageVersion())

	objectMeta := metav1.ObjectMeta{
		Name:        name,
		Namespace:   namespace,
		Annotations: templateOpts.Labels,
		Labels:      labels,
	}
	resourceOptions := templateOpts.Size.ToKubeResource()
	pod, err := buildPod(objectMeta, templateOpts.Template, resourceOptions.Memory, resourceOptions.Cpu)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error kube pod spec")
	}

	deployment := buildDeployment(objectMeta, pod)

	service, err := buildService(objectMeta, templateOpts.Template)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error kube service spec")
	}

	return deployment, service, nil
}

func buildLabels(name, instance, version string) model.BoxLabels {
	// default
	labels := map[string]string{
		"app.kubernetes.io/name":       name,
		"app.kubernetes.io/instance":   common.ToKebabCase(instance),
		"app.kubernetes.io/version":    version,
		"app.kubernetes.io/managed-by": "hckops", // TODO common?
		model.LabelSchemaKind:          common.ToKebabCase(schema.KindBoxV1.String()),
	}
	return labels
}

func buildContainerPorts(ports []model.BoxPort) ([]corev1.ContainerPort, error) {

	containerPorts := make([]corev1.ContainerPort, 0)
	for _, port := range ports {

		portNumber, err := strconv.Atoi(port.Remote)
		if err != nil {
			return nil, errors.Wrap(err, "error kube container port")
		}

		containerPort := corev1.ContainerPort{
			Name:          fmt.Sprintf("%s-svc", port.Alias),
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: int32(portNumber),
		}
		containerPorts = append(containerPorts, containerPort)
	}
	return containerPorts, nil
}

func buildPod(objectMeta metav1.ObjectMeta, template *model.BoxV1, memory string, cpu string) (*corev1.Pod, error) {

	networkPorts := template.NetworkPorts(false)
	containerPorts, err := buildContainerPorts(networkPorts)
	if err != nil {
		return nil, err
	}

	return &corev1.Pod{
		ObjectMeta: objectMeta,
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            common.ToKebabCase(template.Image.Repository),
					Image:           template.ImageName(),
					ImagePullPolicy: corev1.PullIfNotPresent,
					TTY:             true,
					Stdin:           true,
					Ports:           containerPorts,
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							"memory": resource.MustParse(memory),
						},
						Requests: corev1.ResourceList{
							"cpu":    resource.MustParse(cpu),
							"memory": resource.MustParse(memory),
						},
					},
				},
			},
		},
	}, nil
}

func int32Ptr(i int32) *int32 { return &i }

func buildDeployment(objectMeta metav1.ObjectMeta, pod *corev1.Pod) *appsv1.Deployment {

	return &appsv1.Deployment{
		ObjectMeta: objectMeta,
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(kubernetes.SingleReplica), // only 1 replica
			Selector: &metav1.LabelSelector{
				MatchLabels: objectMeta.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: pod.ObjectMeta,
				Spec:       pod.Spec,
			},
		},
	}
}

func buildServicePorts(ports []model.BoxPort) ([]corev1.ServicePort, error) {

	servicePorts := make([]corev1.ServicePort, 0)
	for _, port := range ports {

		portNumber, err := strconv.Atoi(port.Remote)
		if err != nil {
			return nil, errors.Wrap(err, "error kube service port")
		}

		containerPort := corev1.ServicePort{
			Name:       port.Alias,
			Protocol:   corev1.ProtocolTCP,
			Port:       int32(portNumber),
			TargetPort: intstr.FromString(fmt.Sprintf("%s-svc", port.Alias)),
		}
		servicePorts = append(servicePorts, containerPort)
	}
	return servicePorts, nil
}

func buildService(objectMeta metav1.ObjectMeta, template *model.BoxV1) (*corev1.Service, error) {

	networkPorts := template.NetworkPorts(false)
	servicePorts, err := buildServicePorts(networkPorts)
	if err != nil {
		return nil, err
	}

	return &corev1.Service{
		ObjectMeta: objectMeta,
		Spec: corev1.ServiceSpec{
			Selector: objectMeta.Labels,
			Type:     corev1.ServiceTypeClusterIP,
			Ports:    servicePorts,
		},
	}, nil
}

// TODO common
func (box *KubeBoxClient) connectBox(template *model.BoxV1, tunnelOpts *model.TunnelOptions, name string) error {
	if info, err := box.findBox(name); err != nil {
		return err
	} else {
		return box.execBox(template, info, tunnelOpts, false)
	}
}

// TODO common
func (box *KubeBoxClient) openBox(templateOpts *model.TemplateOptions, tunnelOpts *model.TunnelOptions) error {
	if info, err := box.createBox(templateOpts); err != nil {
		return err
	} else {
		return box.execBox(templateOpts.Template, info, tunnelOpts, true)
	}
}

// TODO common
func (box *KubeBoxClient) findBox(name string) (*model.BoxInfo, error) {
	boxes, err := box.listBoxes()
	if err != nil {
		return nil, err
	}
	for _, boxInfo := range boxes {
		if boxInfo.Name == name {
			return &boxInfo, nil
		}
	}
	return nil, errors.New("box not found")
}

func (box *KubeBoxClient) execBox(template *model.BoxV1, info *model.BoxInfo, tunnelOpts *model.TunnelOptions, removeOnExit bool) error {
	box.eventBus.Publish(newPodExecKubeEvent(template.Name, box.clientOpts.Namespace, info.Id, template.Shell))

	if tunnelOpts.TunnelOnly {
		// tunnel and block to exit, wait until killed
		return box.podPortForward(template, info, true)
	}

	if !tunnelOpts.NoTunnel {
		// tunnel and exec after, do not block
		if err := box.podPortForward(template, info, false); err != nil {
			return err
		}
	}

	// TODO if BoxInfo not Healthy attempt scale 1
	// TODO model.BoxShellNone
	// TODO print environment variables

	// exec
	opts := &kubernetes.PodExecOpts{
		Namespace: box.clientOpts.Namespace,
		PodName:   common.ToKebabCase(template.Image.Repository), // pod.Spec.Containers[0].Name
		PodId:     info.Id,
		Shell:     template.Shell,
		InStream:  tunnelOpts.Streams.In,
		OutStream: tunnelOpts.Streams.Out,
		ErrStream: tunnelOpts.Streams.Err,
		IsTty:     tunnelOpts.Streams.IsTty,
		OnExecCallback: func() {
			// stop loader
			box.eventBus.Publish(newPodExecKubeLoaderEvent())
		},
	}

	if removeOnExit {
		defer box.deleteBox(info.Name)
	}

	return box.client.PodExec(opts)
}

func (box *KubeBoxClient) podPortForward(template *model.BoxV1, boxInfo *model.BoxInfo, isWait bool) error {
	namespace := box.clientOpts.Namespace

	if !template.HasPorts() {
		box.eventBus.Publish(newPodPortForwardSkippedKubeEvent(namespace, boxInfo.Id))
		// exit, no service/port available to bind
		return nil
	}

	networkPorts := template.NetworkPorts(false)
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

	info, err := box.findBox(name)
	if err != nil {
		return nil, err
	}

	deployment, err := box.client.DeploymentDescribe(namespace, info.Name)
	if err != nil {
		return nil, err
	}

	service, err := box.client.ServiceDescribe(namespace, info.Name)
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

	// TODO filter by prefix e.g. "HCK_"
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
			Local:  "TODO", // TODO runtime only
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
		Env:     env,
		Ports:   ports,
		Created: deployment.Created,
	}, nil
}

func boxLabel() string {
	// value must be sanitized
	return fmt.Sprintf("%s=%s", model.LabelSchemaKind, common.ToKebabCase(schema.KindBoxV1.String()))
}

func (box *KubeBoxClient) listBoxes() ([]model.BoxInfo, error) {
	namespace := box.clientOpts.Namespace

	deployments, err := box.client.DeploymentList(namespace, model.BoxPrefixName, boxLabel())
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

	boxes, err := box.listBoxes()
	if err != nil {
		return nil, err
	}

	var deleted []string
	for _, boxInfo := range boxes {

		if len(names) == 0 || slices.Contains(names, boxInfo.Name) {

			if err := box.deleteBox(boxInfo.Name); err == nil {
				deleted = append(deleted, boxInfo.Name)
			} else {
				// silently ignore
				box.eventBus.Publish(newResourcesDeleteSkippedKubeEvent(namespace, boxInfo.Name))
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
