package client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/hckops/hckctl/pkg/old/common"
	"github.com/hckops/hckctl/pkg/old/model"
	"github.com/hckops/hckctl/pkg/old/schema"
	"github.com/hckops/hckctl/pkg/old/util"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
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
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/transport/spdy"
	"k8s.io/client-go/util/homedir"
	"k8s.io/kubectl/pkg/cmd/exec"
	"k8s.io/kubectl/pkg/scheme"
)

type KubeBox struct {
	ctx             context.Context
	KubeRestConfig  *rest.Config
	KubeClientSet   *kubernetes.Clientset
	Template        *schema.BoxV1
	ResourceOptions *ResourceOptions

	// TODO look for better design than callbacks
	OnSetupCallback       func(string)
	OnTunnelCallback      func(port schema.PortV1)
	OnTunnelErrorCallback func(error, string)
	OnExecCallback        func()
	OnCloseCallback       func(string)
	OnCloseErrorCallback  func(error, string)
}

type ResourceOptions struct {
	Namespace string
	Memory    string
	Cpu       string
}

func NewOutOfClusterKubeBox(template *schema.BoxV1, resourceOptions *ResourceOptions, configPath string) (*KubeBox, error) {

	restConfig, clientSet, err := NewOutOfClusterClients(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "error out-of-cluster clients")
	}

	return &KubeBox{
		ctx:             context.Background(),
		KubeRestConfig:  restConfig,
		KubeClientSet:   clientSet,
		Template:        template,
		ResourceOptions: resourceOptions,
	}, nil
}

func NewInClusterKubeBox(template *schema.BoxV1, resourceOptions *ResourceOptions) (*KubeBox, error) {

	restConfig, clientSet, err := NewInClusterClients()
	if err != nil {
		return nil, errors.Wrap(err, "error in-cluster clients")
	}

	return &KubeBox{
		ctx:             context.Background(),
		KubeRestConfig:  restConfig,
		KubeClientSet:   clientSet,
		Template:        template,
		ResourceOptions: resourceOptions,
	}, nil
}

func NewOutOfClusterClients(configPath string) (*rest.Config, *kubernetes.Clientset, error) {

	var kubeconfig string
	if strings.TrimSpace(configPath) == "" {
		kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
	} else {
		// absolute path
		kubeconfig = configPath
	}

	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error restConfig")
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error clientSet")
	}

	return restConfig, clientSet, nil
}

func NewInClusterClients() (*rest.Config, *kubernetes.Clientset, error) {

	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, nil, errors.Wrap(err, "error restConfig")
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error clientSet")
	}

	return restConfig, clientSet, nil
}

func (box *KubeBox) BuildSpec(containerName string) (*appsv1.Deployment, *corev1.Service, error) {
	return buildSpec(containerName, box.Template, box.ResourceOptions)
}

func (box *KubeBox) ApplyTemplate(deploymentSpec *appsv1.Deployment, serviceSpec *corev1.Service) error {
	coreClient := box.KubeClientSet.CoreV1()
	appClient := box.KubeClientSet.AppsV1()

	// apply namespace, see https://github.com/kubernetes/client-go/issues/1036
	namespace, err := coreClient.Namespaces().Apply(box.ctx, applyv1.Namespace(box.ResourceOptions.Namespace), metav1.ApplyOptions{FieldManager: "application/apply-patch"})
	if err != nil {
		return errors.Wrapf(err, "error kube apply namespace: %s", box.ResourceOptions.Namespace)
	}
	box.OnSetupCallback(fmt.Sprintf("namespace %s successfully applied", namespace.Name))

	// create deployment
	deployment, err := appClient.Deployments(namespace.Name).Create(box.ctx, deploymentSpec, metav1.CreateOptions{})
	if err != nil {
		return errors.Wrapf(err, "error kube create deployment: %s", deployment.Name)
	}
	box.OnSetupCallback(fmt.Sprintf("deployment %s successfully created", deployment.Name))

	// create service: can't create a service without ports
	var service *corev1.Service
	if box.Template.HasPorts() {
		service, err = coreClient.Services(namespace.Name).Create(box.ctx, serviceSpec, metav1.CreateOptions{})
		if err != nil {
			return errors.Wrapf(err, "error kube create service: %s", service.Name)
		}
		box.OnSetupCallback(fmt.Sprintf("service %s successfully created", service.Name))
	} else {
		box.OnSetupCallback(fmt.Sprint("service not created"))
	}

	// blocks until the deployment is available, then stop watching
	watcher, err := appClient.Deployments(namespace.Name).Watch(box.ctx, metav1.SingleObject(deployment.ObjectMeta))
	if err != nil {
		return errors.Wrapf(err, "error kube watch deployment: %s", deployment.Name)
	}
	for event := range watcher.ResultChan() {
		switch event.Type {
		case watch.Modified:
			deploymentEvent := event.Object.(*appsv1.Deployment)

			for _, condition := range deploymentEvent.Status.Conditions {
				box.OnSetupCallback(fmt.Sprintf("watch kube event: type=%v, condition=%v", event.Type, condition.Message))

				if condition.Type == appsv1.DeploymentAvailable &&
					condition.Status == corev1.ConditionTrue {
					watcher.Stop()
				}
			}

		default:
			return errors.Wrapf(err, "error kube event type: %s", event.Type)
		}
	}
	return nil
}

func (box *KubeBox) RemoveTemplate(deployment *appsv1.Deployment, service *corev1.Service) {
	coreClient := box.KubeClientSet.CoreV1()
	appClient := box.KubeClientSet.AppsV1()

	if err := appClient.Deployments(box.ResourceOptions.Namespace).Delete(box.ctx, deployment.Name, metav1.DeleteOptions{}); err != nil {
		box.OnCloseErrorCallback(err, fmt.Sprintf("error kube delete deployment: %s", deployment.Name))
	}
	box.OnCloseCallback(fmt.Sprintf("deployment %s successfully deleted", deployment.Name))

	if service != nil {
		if err := coreClient.Services(box.ResourceOptions.Namespace).Delete(box.ctx, service.Name, metav1.DeleteOptions{}); err != nil {
			box.OnCloseErrorCallback(err, fmt.Sprintf("error kube delete service: %s", service.Name))
		}
		box.OnCloseCallback(fmt.Sprintf("service %s successfully deleted", service.Name))
	}
}

func (box *KubeBox) GetPod(deployment *appsv1.Deployment) (*corev1.Pod, error) {
	coreClient := box.KubeClientSet.CoreV1()

	labelSet := labels.Set(deployment.Spec.Selector.MatchLabels)
	listOptions := metav1.ListOptions{LabelSelector: labelSet.AsSelector().String()}

	pods, err := coreClient.Pods(deployment.Namespace).List(box.ctx, listOptions)
	if err != nil {
		return nil, errors.Wrap(err, "error list pods")
	}

	if len(pods.Items) != 1 {
		return nil, errors.Wrapf(err, "only 1 pod expected for deployment with labels=%s", deployment.Spec.Selector.MatchLabels)
	}

	pod := pods.Items[0]

	return &pod, nil
}

func (box *KubeBox) PortForward(podName, namespace string) {
	coreClient := box.KubeClientSet.CoreV1()

	if !box.Template.HasPorts() {
		// exit, no service/port available to bind
		return
	}

	var portBindings []string
	for _, port := range box.Template.NetworkPorts() {
		localPort, _ := util.GetLocalPort(port.Local)

		box.OnTunnelCallback(schema.PortV1{
			Alias:  port.Alias,
			Local:  localPort,
			Remote: port.Remote,
		})

		portBindings = append(portBindings, fmt.Sprintf("%s:%s", localPort, port.Remote))
	}

	restRequest := coreClient.RESTClient().
		Post().
		Resource("pods").
		Namespace(namespace).
		Name(podName).
		SubResource("portforward")

	transport, upgrader, err := spdy.RoundTripperFor(box.KubeRestConfig)
	if err != nil {
		box.OnTunnelErrorCallback(err, "error kube round tripper")
	}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, restRequest.URL())

	stopChannel := box.ctx.Done()
	readyChannel := make(chan struct{}, 1)
	out := new(bytes.Buffer)
	errOut := new(bytes.Buffer)

	forwarder, err := portforward.New(dialer, portBindings, stopChannel, readyChannel, out, errOut)
	if err != nil {
		box.OnTunnelErrorCallback(err, "error kube new portforward")
	}

	// wait until interrupted
	go func() {
		if err := forwarder.ForwardPorts(); err != nil {
			box.OnTunnelErrorCallback(err, "error kube forwarding")
		}
	}()
	for range readyChannel {
	}

	if len(errOut.String()) != 0 {
		box.OnTunnelErrorCallback(err, fmt.Sprintf("error kube new portforward: %s", errOut.String()))
	}
}

func (box *KubeBox) Exec(pod *corev1.Pod, streams *model.BoxStreams) error {
	coreClient := box.KubeClientSet.CoreV1()

	streamOptions := buildStreamOptions(streams)
	tty := streamOptions.SetupTTY()

	var sizeQueue remotecommand.TerminalSizeQueue
	if tty.Raw {
		sizeQueue = tty.MonitorSize(tty.GetSize())
		streamOptions.ErrOut = nil
	}

	// TODO add to template or default
	shell := "/bin/bash"

	// exec remote shell
	execUrl := coreClient.RESTClient().
		Post().
		Namespace(pod.Namespace).
		Resource("pods").
		Name(pod.Name).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: box.Template.SafeName(), // pod.Spec.Containers[0].Name
			Command:   []string{shell},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       tty.Raw,
		}, scheme.ParameterCodec).
		URL()

	executor := exec.DefaultRemoteExecutor{}

	box.OnExecCallback()

	fn := func() error {
		return executor.Execute(http.MethodPost, execUrl, box.KubeRestConfig, streamOptions.In, streamOptions.Out, streamOptions.ErrOut, tty.Raw, sizeQueue)
	}
	if err := tty.Safe(fn); err != nil {
		return errors.Wrap(err, "terminal session closed")
	}
	return nil
}

func buildStreamOptions(streams *model.BoxStreams) exec.StreamOptions {
	return exec.StreamOptions{
		Stdin: true,
		TTY:   streams.IsTty,
		IOStreams: genericclioptions.IOStreams{
			In:     streams.Stdin,
			Out:    streams.Stdout,
			ErrOut: streams.Stderr,
		},
	}
}

func buildSpec(containerName string, template *schema.BoxV1, resourceOptions *ResourceOptions) (*appsv1.Deployment, *corev1.Service, error) {

	customLabels := buildLabels(containerName, template.SafeName(), template.ImageVersion())
	objectMeta := metav1.ObjectMeta{
		Name:      containerName,
		Namespace: resourceOptions.Namespace,
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

func buildContainerPorts(ports []schema.PortV1) ([]corev1.ContainerPort, error) {

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

func buildPod(objectMeta metav1.ObjectMeta, template *schema.BoxV1, memory string, cpu string) (*corev1.Pod, error) {

	containerPorts, err := buildContainerPorts(template.NetworkPorts())
	if err != nil {
		return nil, err
	}

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

func buildServicePorts(ports []schema.PortV1) ([]corev1.ServicePort, error) {

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

func buildService(objectMeta metav1.ObjectMeta, template *schema.BoxV1) (*corev1.Service, error) {

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
