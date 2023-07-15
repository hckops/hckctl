package kubernetes

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
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

	"github.com/hckops/hckctl/pkg/client/common"
)

func NewKubeClient(inCluster bool, configPath string) (*KubeClient, error) {
	if inCluster {
		return NewInClusterKubeClient()
	} else {
		return NewOutOfClusterKubeClient(configPath)
	}
}

func NewOutOfClusterKubeClient(configPath string) (*KubeClient, error) {

	var kubeConfig string
	if strings.TrimSpace(configPath) == "" {
		kubeConfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
	} else {
		// absolute path
		kubeConfig = configPath
	}

	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error out-of-cluster restConfig")
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error out-of-cluster clientSet")
	}

	return &KubeClient{
		ctx:            context.Background(),
		kubeRestConfig: restConfig,
		kubeClientSet:  clientSet,
	}, nil
}

func NewInClusterKubeClient() (*KubeClient, error) {

	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, errors.Wrap(err, "error in-cluster restConfig")
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error in-cluster clientSet")
	}

	return &KubeClient{
		ctx:            context.Background(),
		kubeRestConfig: restConfig,
		kubeClientSet:  clientSet,
	}, nil
}

func (client *KubeClient) Close() error {
	// TODO
	return errors.New("not implemented")
}

func (client *KubeClient) NamespaceApply(name string) error {
	coreClient := client.kubeClientSet.CoreV1()

	// https://github.com/kubernetes/client-go/issues/1036
	_, err := coreClient.Namespaces().Apply(client.ctx, applyv1.Namespace(name), metav1.ApplyOptions{FieldManager: "application/apply-patch"})
	if err != nil {
		return errors.Wrapf(err, "error namespace apply: name=%s", name)
	}
	return nil
}

func (client *KubeClient) NamespaceDelete(name string) error {
	coreClient := client.kubeClientSet.CoreV1()

	if err := coreClient.Namespaces().Delete(client.ctx, name, metav1.DeleteOptions{}); err != nil {
		return errors.Wrapf(err, "error namespace delete: name=%s", name)
	}
	return nil
}

func (client *KubeClient) DeploymentCreate(opts *DeploymentCreateOpts) error {
	appClient := client.kubeClientSet.AppsV1()

	deployment, err := appClient.Deployments(opts.Namespace).Create(client.ctx, opts.Spec, metav1.CreateOptions{})
	if err != nil {
		return errors.Wrapf(err, "error deployment create: namespace=%s name=%s", opts.Namespace, opts.Spec.Name)
	}

	// blocks until the deployment is available, then stop watching
	watcher, err := appClient.Deployments(opts.Namespace).Watch(client.ctx, metav1.SingleObject(deployment.ObjectMeta))
	if err != nil {
		return errors.Wrapf(err, "error deployment watch: namespace=%s name=%s", opts.Namespace, deployment.Name)
	}
	for event := range watcher.ResultChan() {
		switch event.Type {
		case watch.Modified:
			deploymentEvent := event.Object.(*appsv1.Deployment)

			for _, condition := range deploymentEvent.Status.Conditions {
				opts.OnStatusEventCallback(fmt.Sprintf("watch deployment event: namespace=%s name=%s type=%v condition=%v",
					opts.Namespace, deployment.Name, event.Type, condition.Message))

				if condition.Type == appsv1.DeploymentAvailable &&
					condition.Status == corev1.ConditionTrue {
					watcher.Stop()
				}
			}
		default:
			return errors.Wrapf(err, "error deployment event: type=%v", event.Type)
		}
	}
	return nil
}

func (client *KubeClient) DeploymentList(namespace string, namePrefix string, label string) ([]DeploymentInfo, error) {
	appClient := client.kubeClientSet.AppsV1()

	deployments, err := appClient.Deployments(namespace).List(client.ctx, metav1.ListOptions{
		// format <LABEL_KEY>=<SANITIZED_LABEL_VALUE>
		LabelSelector: label,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "error deployment list: namespace=%s", namespace)
	}
	var result []DeploymentInfo
	for _, deployment := range deployments.Items {

		// check manually because ListOptions/FieldSelector supports only full name
		if !strings.HasPrefix(deployment.Name, namePrefix) {
			// skip invalid prefix
			continue
		}

		podInfo, err := client.GetPodInfo(&deployment)
		if err != nil {
			// TODO verify if error sidecar container or continue
			return nil, err
		}
		servicePorts, err := client.GetServicePorts(namespace, deployment.Name)
		if err != nil {
			// TODO return or continue (silently ignore)
			return nil, err
		}

		deploymentInfo := DeploymentInfo{
			Namespace:      namespace,
			DeploymentName: deployment.Name,
			PodInfo:        podInfo,
			Healthy:        isDeploymentHealthy(deployment.Status),
			Labels:         deployment.Labels,
			ServicePorts:   servicePorts,
		}
		result = append(result, deploymentInfo)
	}
	return result, nil
}

func isDeploymentHealthy(status appsv1.DeploymentStatus) bool {
	// all conditions must be healthy
	var healthy bool
	for _, condition := range status.Conditions {
		if condition.Status == corev1.ConditionTrue {
			healthy = true
		} else {
			healthy = false
			// mark as unhealthy if at least one condition is invalid e.g. stuck to progressing due to lack of resources
			break
		}
	}
	return healthy
}

func (client *KubeClient) DeploymentDelete(namespace string, name string) error {
	appClient := client.kubeClientSet.AppsV1()

	err := appClient.Deployments(namespace).Delete(client.ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return errors.Wrapf(err, "error deployment delete: namespace=%s name=%s", namespace, name)
	}
	return nil
}

func (client *KubeClient) ServiceCreate(namespace string, spec *corev1.Service) error {
	coreClient := client.kubeClientSet.CoreV1()

	_, err := coreClient.Services(namespace).Create(client.ctx, spec, metav1.CreateOptions{})
	if err != nil {
		return errors.Wrapf(err, "error service create: namespace=%s name=%s", namespace, spec.Name)
	}
	return nil
}

func (client *KubeClient) ServiceDelete(namespace string, name string) error {
	coreClient := client.kubeClientSet.CoreV1()

	if err := coreClient.Services(namespace).Delete(client.ctx, name, metav1.DeleteOptions{}); err != nil {
		return errors.Wrapf(err, "error service delete: namespace=%s name=%s", namespace, name)
	}
	return nil
}

func (client *KubeClient) GetServicePorts(namespace string, name string) ([]ServicePort, error) {
	coreClient := client.kubeClientSet.CoreV1()

	service, err := coreClient.Services(namespace).Get(client.ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "error get service ports: namespace=%s name=%s", namespace, name)
	}

	var ports []ServicePort
	for _, port := range service.Spec.Ports {
		ports = append(ports, ServicePort{Name: port.Name, Port: strconv.Itoa(int(port.Port))})
	}

	return ports, nil
}

func (client *KubeClient) GetPodInfo(deployment *appsv1.Deployment) (*PodInfo, error) {
	coreClient := client.kubeClientSet.CoreV1()

	labelSet := labels.Set(deployment.Spec.Selector.MatchLabels)
	listOptions := metav1.ListOptions{LabelSelector: labelSet.AsSelector().String()}

	pods, err := coreClient.Pods(deployment.Namespace).List(client.ctx, listOptions)
	if err != nil {
		return nil, errors.Wrapf(err, "error pod list: namespace=%s labels=%v", deployment.Namespace, labelSet)
	}
	// TODO verify with sidecar container ??? select by podName == deploymentName
	if len(pods.Items) != 1 {
		return nil, errors.Wrapf(err, "found %d pods, expected only 1 pod for deployment: namespace=%s labels=%v", len(pods.Items), deployment.Namespace, labelSet)
	}

	pod := pods.Items[0]
	info := &PodInfo{
		Id:   pod.ObjectMeta.Name, // pod.Name + unique generated suffix
		Name: pod.Name,
	}
	return info, nil
}

func (client *KubeClient) PodPortForward(opts *PodPortForwardOpts) error {
	coreClient := client.kubeClientSet.CoreV1()

	restRequest := coreClient.RESTClient().
		Post().
		Resource("pods").
		Namespace(opts.Namespace).
		Name(opts.PodId).
		SubResource("portforward")

	transport, upgrader, err := spdy.RoundTripperFor(client.kubeRestConfig)
	if err != nil {
		return errors.Wrap(err, "error kube round tripper")
	}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, restRequest.URL())

	stopChannel := client.ctx.Done()
	readyChannel := make(chan struct{}, 1)
	out := new(bytes.Buffer)
	errOut := new(bytes.Buffer)

	forwarder, err := portforward.New(dialer, opts.Ports, stopChannel, readyChannel, out, errOut)
	if err != nil {
		return errors.Wrap(err, "error kube new portforward")
	}

	go func() {
		if err := forwarder.ForwardPorts(); err != nil {
			opts.OnTunnelErrorCallback(err)
		}
	}()

	if opts.IsWait {
		// wait until interrupted
		select {
		case <-stopChannel:
			return errors.Wrap(err, "error kube stopped")
		}
	} else {
		// continue as soon as ready
		for range readyChannel {
		}
	}

	if len(errOut.String()) != 0 {
		return errors.Wrapf(err, "error kube stream: %s", errOut.String())
	}
	return nil
}

func (client *KubeClient) PodExec(opts *PodExecOpts) error {
	coreClient := client.kubeClientSet.CoreV1()

	streamOptions := exec.StreamOptions{
		Stdin: true,
		TTY:   opts.IsTty,
		IOStreams: genericclioptions.IOStreams{
			In:     opts.InStream,
			Out:    opts.OutStream,
			ErrOut: opts.ErrStream,
		},
	}
	tty := streamOptions.SetupTTY()

	var sizeQueue remotecommand.TerminalSizeQueue
	if tty.Raw {
		sizeQueue = tty.MonitorSize(tty.GetSize())
		streamOptions.ErrOut = nil
	}

	// exec remote shell
	execUrl := coreClient.RESTClient().
		Post().
		Namespace(opts.Namespace).
		Resource("pods").
		Name(opts.PodId).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: opts.PodName,
			Command:   []string{common.DefaultShell(opts.Shell)},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       tty.Raw,
		}, scheme.ParameterCodec).
		URL()

	executor := exec.DefaultRemoteExecutor{}

	opts.OnExecCallback()

	fn := func() error {
		return executor.Execute(http.MethodPost, execUrl, client.kubeRestConfig, streamOptions.In, streamOptions.Out, streamOptions.ErrOut, tty.Raw, sizeQueue)
	}
	if err := tty.Safe(fn); err != nil {
		return errors.Wrap(err, "terminal session closed")
	}
	return nil
}
