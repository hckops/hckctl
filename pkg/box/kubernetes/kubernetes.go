package kubernetes

import (
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/kubernetes"
	"github.com/pkg/errors"
)

//import (
//	"fmt"
//	"strconv"
//
//	"github.com/pkg/errors"
//	appsv1 "k8s.io/api/apps/v1"
//	corev1 "k8s.io/api/core/v1"
//	"k8s.io/apimachinery/pkg/api/resource"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	"k8s.io/apimachinery/pkg/util/intstr"
//
//	"github.com/hckops/hckctl/pkg/box/model"
//	"github.com/hckops/hckctl/pkg/client/kubernetes"
//	"github.com/hckops/hckctl/pkg/command/common"
//)

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

//func (box *KubeBox) BuildSpec(containerName string) (*appsv1.Deployment, *corev1.Service, error) {
//	return buildSpec(containerName, box.Template, box.ResourceOptions)
//}
//
//func buildSpec(containerName string, template *model.BoxV1, resourceOptions *kubernetes.ResourceOptions) (*appsv1.Deployment, *corev1.Service, error) {
//
//	customLabels := buildLabels(containerName, template.SafeName(), template.ImageVersion())
//	objectMeta := metav1.ObjectMeta{
//		Name:      containerName,
//		Namespace: resourceOptions.Namespace,
//		Labels:    customLabels,
//	}
//	pod, err := buildPod(objectMeta, template, resourceOptions.Memory, resourceOptions.Cpu)
//	if err != nil {
//		return nil, nil, errors.Wrap(err, "error kube pod spec")
//	}
//
//	deployment := buildDeployment(objectMeta, pod)
//
//	service, err := buildService(objectMeta, template)
//	if err != nil {
//		return nil, nil, errors.Wrap(err, "error kube service spec")
//	}
//
//	return deployment, service, nil
//}
//
//type Labels map[string]string
//
//func buildLabels(name, instance, version string) Labels {
//	return map[string]string{
//		"app.kubernetes.io/name":       name,
//		"app.kubernetes.io/instance":   instance,
//		"app.kubernetes.io/version":    version,
//		"app.kubernetes.io/managed-by": common.ProjectName,
//	}
//}
//
//func buildContainerPorts(ports []model.BoxPort) ([]corev1.ContainerPort, error) {
//
//	containerPorts := make([]corev1.ContainerPort, 0)
//	for _, port := range ports {
//
//		portNumber, err := strconv.Atoi(port.Remote)
//		if err != nil {
//			return nil, errors.Wrap(err, "error kube container port")
//		}
//
//		containerPort := corev1.ContainerPort{
//			Name:          fmt.Sprintf("%s-svc", port.Alias),
//			Protocol:      corev1.ProtocolTCP,
//			ContainerPort: int32(portNumber),
//		}
//		containerPorts = append(containerPorts, containerPort)
//	}
//	return containerPorts, nil
//}
//
//func buildPod(objectMeta metav1.ObjectMeta, template *model.BoxV1, memory string, cpu string) (*corev1.Pod, error) {
//
//	containerPorts, err := buildContainerPorts(template.NetworkPorts())
//	if err != nil {
//		return nil, err
//	}
//
//	return &corev1.Pod{
//		ObjectMeta: objectMeta,
//		Spec: corev1.PodSpec{
//			Containers: []corev1.Container{
//				{
//					Name:            template.SafeName(),
//					Image:           template.ImageName(),
//					ImagePullPolicy: corev1.PullIfNotPresent,
//					TTY:             true,
//					Stdin:           true,
//					Ports:           containerPorts,
//					Resources: corev1.ResourceRequirements{
//						Limits: corev1.ResourceList{
//							"memory": resource.MustParse(memory),
//						},
//						Requests: corev1.ResourceList{
//							"cpu":    resource.MustParse(cpu),
//							"memory": resource.MustParse(memory),
//						},
//					},
//				},
//			},
//		},
//	}, nil
//}
//
//func int32Ptr(i int32) *int32 { return &i }
//
//func buildDeployment(objectMeta metav1.ObjectMeta, pod *corev1.Pod) *appsv1.Deployment {
//
//	return &appsv1.Deployment{
//		ObjectMeta: objectMeta,
//		Spec: appsv1.DeploymentSpec{
//			Replicas: int32Ptr(1), // only 1 replica
//			Selector: &metav1.LabelSelector{
//				MatchLabels: objectMeta.Labels,
//			},
//			Template: corev1.PodTemplateSpec{
//				ObjectMeta: pod.ObjectMeta,
//				Spec:       pod.Spec,
//			},
//		},
//	}
//}
//
//func buildServicePorts(ports []model.BoxPort) ([]corev1.ServicePort, error) {
//
//	servicePorts := make([]corev1.ServicePort, 0)
//	for _, port := range ports {
//
//		portNumber, err := strconv.Atoi(port.Remote)
//		if err != nil {
//			return nil, errors.Wrap(err, "error kube service port")
//		}
//
//		containerPort := corev1.ServicePort{
//			Name:       port.Alias,
//			Protocol:   corev1.ProtocolTCP,
//			Port:       int32(portNumber),
//			TargetPort: intstr.FromString(fmt.Sprintf("%s-svc", port.Alias)),
//		}
//		servicePorts = append(servicePorts, containerPort)
//	}
//	return servicePorts, nil
//}
//
//func buildService(objectMeta metav1.ObjectMeta, template *model.BoxV1) (*corev1.Service, error) {
//
//	servicePorts, err := buildServicePorts(template.NetworkPorts())
//	if err != nil {
//		return nil, err
//	}
//
//	return &corev1.Service{
//		ObjectMeta: objectMeta,
//		Spec: corev1.ServiceSpec{
//			Selector: objectMeta.Labels,
//			Type:     corev1.ServiceTypeClusterIP,
//			Ports:    servicePorts,
//		},
//	}, nil
//}
