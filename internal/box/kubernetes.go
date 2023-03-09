package box

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	applyv1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/transport/spdy"
	"k8s.io/client-go/util/homedir"
	"k8s.io/kubectl/pkg/cmd/exec"

	"github.com/hckops/hckctl/internal/common"
	model "github.com/hckops/hckctl/internal/model"
	"github.com/hckops/hckctl/internal/terminal"
)

// TODO add log with context?
type KubeBox struct {
	ctx            context.Context
	loader         *terminal.Loader
	config         *model.KubeConfig
	template       *model.BoxV1
	kubeRestConfig *rest.Config
	kubeClientSet  *kubernetes.Clientset
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
		ctx:            context.Background(),
		loader:         terminal.NewLoader(),
		config:         config,
		template:       template,
		kubeRestConfig: restConfig,
		kubeClientSet:  clientSet,
	}
}

func (b *KubeBox) OpenBox(streams *model.BoxStreams) {
	log.Debug().Msgf("init kube box: \n%v\n", b.template.Pretty())
	b.loader.Start(fmt.Sprintf("loading %s", b.template.Name))
	b.loader.Sleep(1)

	pod, deleteResources := b.applyTemplate()
	defer deleteResources()
	log.Info().Msgf("open new box: image=%s, namespace=%s, podName=%s", b.template.ImageName(), pod.Namespace, pod.Name)

	b.portForwardPod(pod)
	b.execPod(pod, streams)
}

func (b *KubeBox) applyTemplate() (*corev1.Pod, func()) {
	coreClient := b.kubeClientSet.CoreV1()
	appClient := b.kubeClientSet.AppsV1()

	// apply namespace, see https://github.com/kubernetes/client-go/issues/1036
	namespace, err := coreClient.Namespaces().Apply(b.ctx, applyv1.Namespace(b.config.Namespace), metav1.ApplyOptions{FieldManager: "application/apply-patch"})
	if err != nil {
		log.Fatal().Err(err).Msgf("error kube apply namespace: %s", namespace.Name)
	}
	log.Debug().Msgf("namespace %s successfully applied", namespace.Name)

	containerName := b.template.GenerateName()
	deploymentSpec, serviceSpec := buildSpec(namespace.Name, containerName, b.template, b.config)

	b.loader.Refresh(fmt.Sprintf("creating %s/%s", namespace.Name, containerName))

	// create deployment
	deployment, err := appClient.Deployments(namespace.Name).Create(b.ctx, deploymentSpec, metav1.CreateOptions{})
	if err != nil {
		log.Fatal().Err(err).Msgf("error kube create deployment: %s", deployment.Name)
	}
	log.Debug().Msgf("deployment %s successfully created", deployment.Name)

	// create service: can't create a service without ports
	var service *corev1.Service
	if b.template.HasPorts() {
		service, err = coreClient.Services(namespace.Name).Create(b.ctx, serviceSpec, metav1.CreateOptions{})
		if err != nil {
			log.Fatal().Err(err).Msgf("error kube create service: %s", service.Name)
		}
		log.Debug().Msgf("service %s successfully created", service.Name)
	} else {
		log.Debug().Msg("service not created")
	}

	// blocks until the deployment is available, then stop watching
	watcher, err := appClient.Deployments(namespace.Name).Watch(b.ctx, metav1.SingleObject(deployment.ObjectMeta))
	if err != nil {
		log.Fatal().Err(err).Msgf("error kube watch deployment: %s", deployment.Name)
	}
	for event := range watcher.ResultChan() {
		switch event.Type {
		case watch.Modified:
			deploymentEvent := event.Object.(*appsv1.Deployment)

			for _, condition := range deploymentEvent.Status.Conditions {
				log.Debug().Msgf("watch kube event: type=%v, condition=%v", event.Type, condition.Message)

				if condition.Type == appsv1.DeploymentAvailable &&
					condition.Status == corev1.ConditionTrue {
					watcher.Stop()
				}
			}

		default:
			log.Fatal().Msgf("error kube event type: %s", event.Type)
		}
	}

	// find unique pod for deployment
	pod := b.getPod(deployment)

	cleanupCallback := func() {
		if err := appClient.Deployments(namespace.Name).Delete(b.ctx, deployment.Name, metav1.DeleteOptions{}); err != nil {
			log.Fatal().Err(err).Msgf("error kube delete deployment: %s", deployment.Name)
		}
		log.Debug().Msgf("deployment %s successfully deleted", deployment.Name)

		if service != nil {
			if err := coreClient.Services(namespace.Name).Delete(b.ctx, service.Name, metav1.DeleteOptions{}); err != nil {
				log.Fatal().Err(err).Msgf("error kube delete service: %s", service.Name)
			}
			log.Debug().Msgf("service %s successfully deleted", service.Name)
		}
	}

	return pod, cleanupCallback
}

func (b *KubeBox) getPod(deployment *appsv1.Deployment) *corev1.Pod {
	coreClient := b.kubeClientSet.CoreV1()

	labelSet := labels.Set(deployment.Spec.Selector.MatchLabels)
	listOptions := metav1.ListOptions{LabelSelector: labelSet.AsSelector().String()}

	pods, err := coreClient.Pods(deployment.Namespace).List(b.ctx, listOptions)
	if err != nil {
		log.Fatal().Err(err).Msg("error list pods")
	}

	if len(pods.Items) != 1 {
		log.Fatal().Msgf("only 1 pod expected for deployment with labels=%s", deployment.Spec.Selector.MatchLabels)
	}

	pod := pods.Items[0]
	log.Debug().Msgf("found matching pod %s", pod.Name)

	return &pod
}

func (b *KubeBox) portForwardPod(pod *corev1.Pod) {
	coreClient := b.kubeClientSet.CoreV1()

	if !b.template.HasPorts() {
		// exit, no service/port available to bind
		return
	}

	var portBindings []string
	for _, port := range b.template.NetworkPorts() {
		localPort := common.GetLocalPort(port.Local)
		log.Info().Msgf("[%s] forwarding %s (local) -> %s (remote)", port.Alias, localPort, port.Remote)

		portBindings = append(portBindings, fmt.Sprintf("%s:%s", localPort, port.Remote))
	}

	restRequest := coreClient.RESTClient().
		Post().
		Resource("pods").
		Namespace(pod.Namespace).
		Name(pod.Name).
		SubResource("portforward")

	transport, upgrader, err := spdy.RoundTripperFor(b.kubeRestConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("error kube round tripper")
	}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, restRequest.URL())

	stopChannel := b.ctx.Done()
	readyChannel := make(chan struct{}, 1)
	out := new(bytes.Buffer)
	errOut := new(bytes.Buffer)

	forwarder, err := portforward.New(dialer, portBindings, stopChannel, readyChannel, out, errOut)
	if err != nil {
		log.Fatal().Err(err).Msg("error kube new portforward")
	}

	// wait until interrupted
	log.Debug().Msgf("forwarding all ports for pod %s", pod.Name)
	go func() {
		if err := forwarder.ForwardPorts(); err != nil {
			log.Fatal().Err(err).Msg("error kube forwarding")
		}
	}()
	for range readyChannel {
	}

	if len(errOut.String()) != 0 {
		log.Fatal().Msgf("error kube new portforward: %s", errOut.String())
	}
}

func (b *KubeBox) execPod(pod *corev1.Pod, streams *model.BoxStreams) {
	coreClient := b.kubeClientSet.CoreV1()

	streamOptions := exec.StreamOptions{
		Stdin: true,
		TTY:   streams.IsTty,
		IOStreams: genericclioptions.IOStreams{
			In:     streams.Stdin,
			Out:    streams.Stdout,
			ErrOut: streams.Stderr,
		},
	}

	tty := streamOptions.SetupTTY()

	var sizeQueue remotecommand.TerminalSizeQueue
	if tty.Raw {
		sizeQueue = tty.MonitorSize(tty.GetSize())
		streamOptions.ErrOut = nil
	}

	// exec remote shell
	restRequest := coreClient.RESTClient().
		Post().
		Namespace(pod.Namespace).
		Resource("pods").
		Name(pod.Name).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: pod.Spec.Containers[0].Name,
			Command:   []string{"/bin/bash"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       tty.Raw,
		}, scheme.ParameterCodec)

	executor := exec.DefaultRemoteExecutor{}

	log.Debug().Msgf("exec into pod %s", pod.Name)
	b.loader.Stop()

	fn := func() error {
		return executor.Execute(http.MethodPost, restRequest.URL(), b.kubeRestConfig, streamOptions.In, streamOptions.Out, streamOptions.ErrOut, tty.Raw, sizeQueue)
	}
	if err := tty.Safe(fn); err != nil {
		log.Warn().Err(err).Msg("terminal session closed")
	}
}

func buildSpec(namespaceName string, containerName string, template *model.BoxV1, config *model.KubeConfig) (*appsv1.Deployment, *corev1.Service) {

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
		"app.kubernetes.io/managed-by": common.ProjectName,
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
