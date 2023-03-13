package client

import (
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/hckops/hckctl/internal/config"
	"github.com/hckops/hckctl/pkg/schema"
)

func BuildSpec(namespaceName string, containerName string, template *schema.BoxV1, config *config.KubeConfig) (*appsv1.Deployment, *corev1.Service) {

	labels := buildLabels(containerName, template.SafeName(), template.ImageVersion())
	objectMeta := metav1.ObjectMeta{
		Name:      containerName,
		Namespace: namespaceName,
		Labels:    labels,
	}
	pod := buildPod(objectMeta, template, config)
	deployment := buildDeployment(objectMeta, pod)
	service := buildService(objectMeta, template)

	return deployment, service
}

type Labels map[string]string

func buildLabels(name, instance, version string) Labels {
	return map[string]string{
		"app.kubernetes.io/name":       name,
		"app.kubernetes.io/instance":   instance,
		"app.kubernetes.io/version":    version,
		"app.kubernetes.io/managed-by": config.ProjectName,
	}
}

func buildContainerPorts(ports []schema.PortV1) []corev1.ContainerPort {

	containerPorts := make([]corev1.ContainerPort, 0)
	for _, port := range ports {

		portNumber, err := strconv.Atoi(port.Remote)
		if err != nil {
			log.Fatal().Err(err).Msg("error kube container port")
		}

		containerPort := corev1.ContainerPort{
			Name:          fmt.Sprintf("%s-svc", port.Alias),
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: int32(portNumber),
		}
		containerPorts = append(containerPorts, containerPort)
	}
	return containerPorts
}

func buildPod(objectMeta metav1.ObjectMeta, template *schema.BoxV1, config *config.KubeConfig) *corev1.Pod {

	containerPorts := buildContainerPorts(template.NetworkPorts())

	return &corev1.Pod{
		ObjectMeta: objectMeta,
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            template.SafeName(),
					Image:           template.ImageName(),
					ImagePullPolicy: corev1.PullIfNotPresent,
					TTY:             true,
					Stdin:           true,
					Ports:           containerPorts,
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							"memory": resource.MustParse(config.Resources.Memory),
						},
						Requests: corev1.ResourceList{
							"cpu":    resource.MustParse(config.Resources.Cpu),
							"memory": resource.MustParse(config.Resources.Memory),
						},
					},
				},
			},
		},
	}
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

func buildServicePorts(ports []schema.PortV1) []corev1.ServicePort {

	servicePorts := make([]corev1.ServicePort, 0)
	for _, port := range ports {

		portNumber, err := strconv.Atoi(port.Remote)
		if err != nil {
			log.Fatal().Err(err).Msg("error kube service port")
		}

		containerPort := corev1.ServicePort{
			Name:       port.Alias,
			Protocol:   corev1.ProtocolTCP,
			Port:       int32(portNumber),
			TargetPort: intstr.FromString(fmt.Sprintf("%s-svc", port.Alias)),
		}
		servicePorts = append(servicePorts, containerPort)
	}
	return servicePorts
}

func buildService(objectMeta metav1.ObjectMeta, template *schema.BoxV1) *corev1.Service {

	servicePorts := buildServicePorts(template.NetworkPorts())

	return &corev1.Service{
		ObjectMeta: objectMeta,
		Spec: corev1.ServiceSpec{
			Selector: objectMeta.Labels,
			Type:     corev1.ServiceTypeClusterIP,
			Ports:    servicePorts,
		},
	}
}
