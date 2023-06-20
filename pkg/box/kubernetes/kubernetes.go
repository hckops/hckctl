package kubernetes

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/kubernetes"
	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/util"
)

func newKubeBox(internalOpts *model.BoxInternalOpts, kubeConfig *kubernetes.KubeClientConfig) (*KubeBox, error) {
	internalOpts.EventBus.Publish(newClientInitKubeEvent())

	kubeClient, err := kubernetes.NewOutOfClusterKubeClient(kubeConfig.ConfigPath)
	if err != nil {
		return nil, errors.Wrap(err, "error kube box")
	}

	return &KubeBox{
		client:       kubeClient,
		clientConfig: kubeConfig,
		streams:      internalOpts.Streams,
		eventBus:     internalOpts.EventBus,
	}, nil
}

func (box *KubeBox) close() error {
	box.eventBus.Publish(newClientCloseKubeEvent())
	box.eventBus.Close()
	return box.client.Close()
}

func (box *KubeBox) createBox(template *model.BoxV1) (*model.BoxInfo, error) {
	namespace := box.clientConfig.Namespace

	// boxName
	containerName := template.GenerateName()
	deployment, service, err := buildSpec(containerName, namespace, template, box.clientConfig.Resource)
	if err != nil {
		return nil, err
	}

	// TODO LOADER start ??? local.loader.Refresh(fmt.Sprintf("creating %s/%s", local.box.ResourceOptions.Namespace, containerName))
	// TODO LOADER stop ??? local.loader.Halt(err, "error kube: apply template")
	if box.client.NamespaceApply(namespace) != nil {
		return nil, err
	}
	// TODO DEBUG box.OnSetupCallback(fmt.Sprintf("namespace %s successfully applied", namespace.Name))

	if template.HasPorts() {
		if box.client.ServiceCreate(namespace, service) != nil {
			return nil, err
		}
		// TODO DEBUG box.OnSetupCallback(fmt.Sprintf("service %s successfully created", service.Name))
	} else {
		// TODO DEBUG box.OnSetupCallback(fmt.Sprint("service not created"))
	}

	deploymentOpts := &kubernetes.DeploymentCreateOpts{
		Namespace: namespace,
		Spec:      deployment,
		OnStatusEventCallback: func(event string) {
			// TODO DEBUG box.OnSetupCallback(fmt.Sprintf("watch kube event: type=%v, condition=%v", event.Type, condition.Message))
		},
	}
	if box.client.DeploymentCreate(deploymentOpts) != nil {
		return nil, err
	}
	// TODO DEBUG box.OnSetupCallback(fmt.Sprintf("deployment %s successfully created", deployment.Name))

	// boxId
	podId, err := box.client.PodName(deployment)
	if err != nil {
		return nil, err
	}
	// TODO debug local.log.Debug().Msgf("found matching pod %s", pod.Name)

	return &model.BoxInfo{Id: podId, Name: containerName}, nil
}

func buildSpec(containerName string, namespace string, template *model.BoxV1, resourceOptions *kubernetes.KubeResource) (*appsv1.Deployment, *corev1.Service, error) {

	customLabels := buildLabels(containerName, util.ToKebabCase(template.Image.Repository), template.ImageVersion())
	objectMeta := metav1.ObjectMeta{
		Name:      containerName,
		Namespace: namespace,
		Labels:    customLabels,
	}
	pod, err := buildPod(objectMeta, template, resourceOptions.Memory, resourceOptions.Cpu)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error kube pod spec")
	}

	deployment := buildDeployment(objectMeta, pod)

	service, err := buildService(objectMeta, template)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error kube service spec")
	}

	return deployment, service, nil
}

type Labels map[string]string

func buildLabels(name, instance, version string) Labels {
	return map[string]string{
		"app.kubernetes.io/name":       name,
		"app.kubernetes.io/instance":   instance,
		"app.kubernetes.io/version":    version,
		"app.kubernetes.io/managed-by": common.ProjectName,
	}
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

	containerPorts, err := buildContainerPorts(template.NetworkPorts())
	if err != nil {
		return nil, err
	}

	return &corev1.Pod{
		ObjectMeta: objectMeta,
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            util.ToKebabCase(template.Image.Repository),
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
			Replicas: int32Ptr(1), // only 1 replica
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

	servicePorts, err := buildServicePorts(template.NetworkPorts())
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

func (box *KubeBox) openBox(template *model.BoxV1) error {

	info, err := box.createBox(template)
	if err != nil {
		return err
	}
	if err := box.portForward(template, info.Name); err != nil {
		return err
	}

	// TODO model.BoxShellNone
	// TODO exec

	return nil
}

// TODO deploymentName or podId
func (box *KubeBox) portForward(template *model.BoxV1, name string) error {

	if !template.HasPorts() {
		// exit, no service/port available to bind
		return nil
	}
	ports, err := toPortBindings(template.NetworkPorts(), func(port model.BoxPort) {
		// TODO OnTunnelCallback
	})
	if err != nil {
		return err
	}

	opts := &kubernetes.PortForwardOpts{
		Namespace: box.clientConfig.Namespace,
		PodName:   name,
		Ports:     ports,
		OnTunnelErrorCallback: func(err error) {
			// TODO event
		},
	}
	if err := box.client.PortForward(opts); err != nil {
		return err
	}

	return nil
}

func toPortBindings(ports []model.BoxPort, onPortBindCallback func(port model.BoxPort)) ([]string, error) {

	var portBindings []string
	for _, port := range ports {

		localPort, err := util.FindOpenPort(port.Local)
		if err != nil {
			return nil, errors.Wrap(err, "error kube local port: portForward")
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

func (box *KubeBox) listBoxes() ([]model.BoxInfo, error) {

	deployments, err := box.client.DeploymentList(box.clientConfig.Namespace)
	if err != nil {
		return nil, err
	}
	var result []model.BoxInfo
	for _, d := range deployments {
		result = append(result, model.BoxInfo{Id: d.PodId, Name: d.DeploymentName})
		// TODO box.eventBus.Publish(newContainerListDockerEvent(index, c.ContainerName, c.ContainerId))
	}

	return result, nil
}

func (box *KubeBox) deleteBox(name string) error {
	// TODO box.eventBus.Publish(newContainerRemoveDockerEvent(id))
	namespace := box.clientConfig.Namespace

	if err := box.client.DeploymentDelete(namespace, name); err != nil {
		return err
	}
	if err := box.client.ServiceDelete(namespace, name); err != nil {
		return err
	}
	return nil
}

func (box *KubeBox) deleteBoxes() ([]model.BoxInfo, error) {
	boxes, err := box.listBoxes()
	if err != nil {
		return nil, err
	}
	var deleted []model.BoxInfo
	for _, boxInfo := range boxes {
		if err := box.deleteBox(boxInfo.Name); err == nil {
			deleted = append(deleted, boxInfo)
		} else {
			// TODO same for docker
			// TODO box.eventBus.Publish: silently ignore error
		}
	}
	if err := box.client.NamespaceDelete(box.clientConfig.Namespace); err != nil {
		// TODO box.eventBus.Publish: silently ignore error
	}
	return deleted, nil
}
