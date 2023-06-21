package kubernetes

import (
	"bytes"
	"context"
	"fmt"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/cmd/exec"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/watch"
	applyv1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"k8s.io/client-go/util/homedir"
	"k8s.io/kubectl/pkg/scheme"
)

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
		return nil, errors.Wrap(err, "error restConfig")
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error clientSet")
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
		return nil, errors.Wrap(err, "error restConfig")
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error clientSet")
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

func (client *KubeClient) DeploymentList(namespace string) ([]DeploymentInfo, error) {
	appClient := client.kubeClientSet.AppsV1()

	// TODO filter list: "box-" prefix and status running ?
	deployments, err := appClient.Deployments(namespace).List(client.ctx, metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "error deployment list: namespace=%s", namespace)
	}
	var result []DeploymentInfo
	for _, deployment := range deployments.Items {

		if podInfo, err := client.GetPodInfo(&deployment); err != nil {
			// TODO verify if error sidecar container
			return nil, err
		} else {
			deploymentInfo := DeploymentInfo{
				Namespace:      namespace,
				DeploymentName: deployment.Name,
				PodInfo:        podInfo,
			}
			result = append(result, deploymentInfo)
		}
	}
	return result, nil
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
		Id:   pod.ObjectMeta.Name, // pod.Name + unique generated name
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
	// TODO alternative to callback (uncomment)
	//failedChannel := make(chan error, 1)
	out := new(bytes.Buffer)
	errOut := new(bytes.Buffer)

	forwarder, err := portforward.New(dialer, opts.Ports, stopChannel, readyChannel, out, errOut)
	if err != nil {
		return errors.Wrap(err, "error kube new portforward")
	}

	// wait until interrupted
	go func() {
		if err := forwarder.ForwardPorts(); err != nil {
			// TODO alternative to callback: verify if callback is more stable i.e. ignore errors (replace line below)
			// failedChannel <- err
			opts.OnTunnelErrorCallback(err)
		}
	}()
	for range readyChannel {
	}
	// TODO alternative to callback (replace line above)
	//select {
	//case err = <-failedChannel:
	//	return errors.Wrap(err, "error kube forwarding")
	//case <-readyChannel:
	//}

	if len(errOut.String()) != 0 {
		return errors.Wrapf(err, "error kube portforward: %s", errOut.String())
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
			Command:   []string{opts.Shell},
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
