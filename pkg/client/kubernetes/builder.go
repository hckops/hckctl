package kubernetes

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/hckops/hckctl/pkg/util"
)

func BuildResources(opts *ResourcesOpts) (*appsv1.Deployment, *corev1.Service, error) {

	objectMeta := metav1.ObjectMeta{
		Name:        opts.Name,
		Namespace:   opts.Namespace,
		Annotations: opts.Annotations,
		Labels:      opts.Labels,
	}
	pod, err := buildPod(objectMeta, opts.PodInfo, opts.Ports)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error kube pod spec")
	}

	deployment := buildDeployment(objectMeta, pod)

	service, err := buildService(objectMeta, opts.Ports)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error kube service spec")
	}

	return deployment, service, nil
}

func BuildLabels(name, instance, version string, extra map[string]string) map[string]string {
	// default
	labels := map[string]string{
		LabelKubeName:      name,
		LabelKubeInstance:  util.ToLowerKebabCase(instance),
		LabelKubeVersion:   version,
		LabelKubeManagedBy: "hckops", // TODO common?
	}
	maps.Copy(labels, extra)
	return labels
}

func buildPod(objectMeta metav1.ObjectMeta, podInfo *PodInfo, ports []KubePort) (*corev1.Pod, error) {

	containerPorts, err := buildContainerPorts(ports)
	if err != nil {
		return nil, err
	}

	return &corev1.Pod{
		ObjectMeta: objectMeta,
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            util.ToLowerKebabCase(podInfo.ContainerName),
					Image:           podInfo.ImageName,
					ImagePullPolicy: corev1.PullIfNotPresent,
					TTY:             true,
					Stdin:           true,
					Ports:           containerPorts,
					Env:             buildEnvVars(podInfo.Env),
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceMemory: resource.MustParse(podInfo.Resource.Memory),
						},
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse(podInfo.Resource.Cpu),
							corev1.ResourceMemory: resource.MustParse(podInfo.Resource.Memory),
						},
					},
				},
			},
		},
	}, nil
}

func buildEnvVars(envs []KubeEnv) []corev1.EnvVar {
	var envVars []corev1.EnvVar
	for _, env := range envs {
		envVars = append(envVars, corev1.EnvVar{Name: env.Key, Value: env.Value})
	}
	return envVars
}

func buildContainerPorts(ports []KubePort) ([]corev1.ContainerPort, error) {

	containerPorts := make([]corev1.ContainerPort, 0)
	for _, p := range ports {

		portNumber, err := strconv.Atoi(p.Port)
		if err != nil {
			return nil, errors.Wrap(err, "error kube container port")
		}

		containerPort := corev1.ContainerPort{
			Name:          fmt.Sprintf("%s-svc", p.Name),
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: int32(portNumber),
		}
		containerPorts = append(containerPorts, containerPort)
	}
	return containerPorts, nil
}

func int32Ptr(i int32) *int32 { return &i }

func buildDeployment(objectMeta metav1.ObjectMeta, pod *corev1.Pod) *appsv1.Deployment {

	return &appsv1.Deployment{
		ObjectMeta: objectMeta,
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(SingleReplica), // only 1 replica
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

func buildService(objectMeta metav1.ObjectMeta, ports []KubePort) (*corev1.Service, error) {

	servicePorts, err := buildServicePorts(ports)
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

func buildServicePorts(ports []KubePort) ([]corev1.ServicePort, error) {

	servicePorts := make([]corev1.ServicePort, 0)
	for _, p := range ports {

		portNumber, err := strconv.Atoi(p.Port)
		if err != nil {
			return nil, errors.Wrap(err, "error kube service port")
		}

		containerPort := corev1.ServicePort{
			Name:       p.Name,
			Protocol:   corev1.ProtocolTCP,
			Port:       int32(portNumber),
			TargetPort: intstr.FromString(fmt.Sprintf("%s-svc", p.Name)),
		}
		servicePorts = append(servicePorts, containerPort)
	}
	return servicePorts, nil
}

func BuildJob(opts *JobOpts) *batchv1.Job {

	objectMeta := metav1.ObjectMeta{
		Name:        opts.Name,
		Namespace:   opts.Namespace,
		Annotations: opts.Annotations,
		Labels:      opts.Labels,
	}

	return &batchv1.Job{
		ObjectMeta: objectMeta,
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            util.ToLowerKebabCase(opts.PodInfo.ContainerName),
							Image:           opts.PodInfo.ImageName,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command:         []string{},
							Args:            opts.PodInfo.Arguments,
							Env:             buildEnvVars(opts.PodInfo.Env),
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
			BackoffLimit: int32Ptr(0), // attempt only once
		},
	}
}
