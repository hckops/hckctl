package box

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	applyv1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	model "github.com/hckops/hckctl/internal/model"
	"github.com/hckops/hckctl/internal/terminal"
)

// TODO add log?
type KubeBox struct {
	ctx       context.Context
	loader    *terminal.Loader
	config    *model.KubeConfig
	template  *model.BoxV1
	clientSet *kubernetes.Clientset
}

func NewKubeBox(template *model.BoxV1, config *model.KubeConfig) *KubeBox {

	kubeconfig := filepath.Join(homedir.HomeDir(), strings.ReplaceAll(config.ConfigPath, "~/", ""))
	log.Debug().Msgf("read config: configPath=%s, kubeconfig=%s", config.ConfigPath, kubeconfig)

	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal().Err(err).Msg("error restConfig")
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("error clientSet")
	}

	return &KubeBox{
		ctx:       context.Background(),
		loader:    terminal.NewLoader(),
		config:    config,
		template:  template,
		clientSet: clientSet,
	}
}

func (b *KubeBox) Init() {
	coreClient := b.clientSet.CoreV1()
	appClient := b.clientSet.AppsV1()

	log.Debug().Msgf("init kube box: \n%v\n", b.template.Pretty())
	b.loader.Start(fmt.Sprintf("loading %s", b.template.Name))
	b.loader.Sleep(1)

	// apply namespace, see https://github.com/kubernetes/client-go/issues/1036
	namespace, err := coreClient.Namespaces().Apply(b.ctx, applyv1.Namespace(b.config.Namespace), metav1.ApplyOptions{FieldManager: "application/apply-patch"})
	if err != nil {
		log.Fatal().Err(err).Msg("error apply namespace")
	}
	log.Debug().Msgf("namespace %s successfully applied", namespace.Name)

	containerName := b.template.GenerateFullName()
	deploymentSpec, serviceSpec := buildSpec(namespace.Name, containerName, b.template, b.config)

	// create deployment
	deployment, err := appClient.Deployments(namespace.Name).Create(b.ctx, deploymentSpec, metav1.CreateOptions{})
	if err != nil {
		log.Fatal().Err(err).Msg("error create deployment")
	}
	defer appClient.Deployments(namespace.Name).Delete(b.ctx, deployment.Name, metav1.DeleteOptions{})
	log.Debug().Msgf("deployment %s successfully created", deployment.Name)

	// create service
	if b.template.HasPorts() {
		service, err := coreClient.Services(namespace.Name).Create(b.ctx, serviceSpec, metav1.CreateOptions{})
		if err != nil {
			log.Fatal().Err(err).Msg("error create service")
		}
		defer coreClient.Services(namespace.Name).Delete(b.ctx, service.Name, metav1.DeleteOptions{})
		log.Debug().Msgf("service %s successfully created", service.Name)
	} else {
		log.Debug().Msg("service not created")
	}

	b.loader.Refresh(fmt.Sprintf("creating %s/%s", namespace.Name, containerName))
	// TODO defer close + exec + forward
	b.loader.Sleep(5)

	b.loader.Stop()
}

func buildSpec(namespaceName string, containerName string, template *model.BoxV1, config *model.KubeConfig) (*appsv1.Deployment, *corev1.Service) {

	labels := buildLabels(containerName, template.ImageVersion())
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

func buildLabels(name, version string) Labels {
	return map[string]string{
		"app.kubernetes.io/name":    name,
		"app.kubernetes.io/version": version,
	}
}

func buildContainerPorts(ports []model.PortV1) []corev1.ContainerPort {

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

func buildPod(objectMeta metav1.ObjectMeta, template *model.BoxV1, config *model.KubeConfig) *corev1.Pod {

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
			Replicas: int32Ptr(1),
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

func buildServicePorts(ports []model.PortV1) []corev1.ServicePort {

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

func buildService(objectMeta metav1.ObjectMeta, template *model.BoxV1) *corev1.Service {

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
