package kubernetes

//import (
//	"bytes"
//	"context"
//	"fmt"
//	"net/http"
//	"path/filepath"
//	"strings"
//
//	"k8s.io/apimachinery/pkg/labels"
//	"k8s.io/client-go/tools/portforward"
//	"k8s.io/client-go/tools/remotecommand"
//	"k8s.io/client-go/transport/spdy"
//	"k8s.io/kubectl/pkg/cmd/exec"
//	"github.com/pkg/errors"
//	appsv1 "k8s.io/api/apps/v1"
//	corev1 "k8s.io/api/core/v1"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	"k8s.io/apimachinery/pkg/watch"
//	applyv1 "k8s.io/client-go/applyconfigurations/core/v1"
//	"k8s.io/client-go/kubernetes"
//	"k8s.io/client-go/rest"
//	"k8s.io/client-go/tools/clientcmd"
//	"k8s.io/client-go/util/homedir"
//
//	"github.com/hckops/hckctl/internal/schema" // TODO remove
//	"github.com/hckops/hckctl/pkg/util"
//)

import (
	"context"
	"fmt"
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
	"k8s.io/client-go/util/homedir"
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
				// TODO fields only?
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

type DeploymentInfo struct {
	Namespace      string
	DeploymentName string
	PodId          string
}

func (client *KubeClient) DeploymentList(namespace string) ([]DeploymentInfo, error) {
	appClient := client.kubeClientSet.AppsV1()

	deployments, err := appClient.Deployments(namespace).List(client.ctx, metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "error deployment list: namespace=%s", namespace)
	}
	var result []DeploymentInfo
	for _, deployment := range deployments.Items {

		if podId, err := client.PodName(&deployment); err != nil {
			// TODO verify error sidecar container
			return nil, err
		} else {
			info := DeploymentInfo{
				Namespace:      namespace,
				DeploymentName: deployment.Name,
				PodId:          podId,
			}
			result = append(result, info)
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

func (client *KubeClient) PodName(deployment *appsv1.Deployment) (string, error) {
	coreClient := client.kubeClientSet.CoreV1()

	labelSet := labels.Set(deployment.Spec.Selector.MatchLabels)
	listOptions := metav1.ListOptions{LabelSelector: labelSet.AsSelector().String()}

	pods, err := coreClient.Pods(deployment.Namespace).List(client.ctx, listOptions)
	if err != nil {
		return "", errors.Wrapf(err, "error pod list: namespace=%s labels=%v", deployment.Namespace, labelSet)
	}
	// TODO verify with sidecar container ??? select by podName == deploymentName
	if len(pods.Items) != 1 {
		return "", errors.Wrapf(err, "found %d pods, expected only 1 pod for deployment: namespace=%s labels=%v", len(pods.Items), deployment.Namespace, labelSet)
	}

	pod := pods.Items[0]

	// podId = pod.Name + unique generated name
	return pod.ObjectMeta.Name, nil
}

//func (client *KubeClient) PortForward(podName, namespace string) {
//	coreClient := client.kubeClientSet.CoreV1()
//
//	if !box.Template.HasPorts() {
//		// exit, no service/port available to bind
//		return
//	}
//
//	var portBindings []string
//	for _, port := range box.Template.NetworkPorts() {
//		localPort, _ := util.FindOpenPort(port.Local)
//
//		box.OnTunnelCallback(schema.PortV1{
//			Alias:  port.Alias,
//			Local:  localPort,
//			Remote: port.Remote,
//		})
//
//		portBindings = append(portBindings, fmt.Sprintf("%s:%s", localPort, port.Remote))
//	}
//
//	restRequest := coreClient.RESTClient().
//		Post().
//		Resource("pods").
//		Namespace(namespace).
//		Name(podName).
//		SubResource("portforward")
//
//	transport, upgrader, err := spdy.RoundTripperFor(client.kubeRestConfig)
//	if err != nil {
//		box.OnTunnelErrorCallback(err, "error kube round tripper")
//	}
//	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, restRequest.URL())
//
//	stopChannel := client.ctx.Done()
//	readyChannel := make(chan struct{}, 1)
//	out := new(bytes.Buffer)
//	errOut := new(bytes.Buffer)
//
//	forwarder, err := portforward.New(dialer, portBindings, stopChannel, readyChannel, out, errOut)
//	if err != nil {
//		box.OnTunnelErrorCallback(err, "error kube new portforward")
//	}
//
//	// wait until interrupted
//	go func() {
//		if err := forwarder.ForwardPorts(); err != nil {
//			box.OnTunnelErrorCallback(err, "error kube forwarding")
//		}
//	}()
//	for range readyChannel {
//	}
//
//	if len(errOut.String()) != 0 {
//		box.OnTunnelErrorCallback(err, fmt.Sprintf("error kube new portforward: %s", errOut.String()))
//	}
//}
//
//func (client *KubeClient) Exec(pod *corev1.Pod, streams *model.BoxStreams) error {
//	coreClient := client.kubeClientSet.CoreV1()
//
//	streamOptions := buildStreamOptions(streams)
//	tty := streamOptions.SetupTTY()
//
//	var sizeQueue remotecommand.TerminalSizeQueue
//	if tty.Raw {
//		sizeQueue = tty.MonitorSize(tty.GetSize())
//		streamOptions.ErrOut = nil
//	}
//
//	// TODO add to template or default
//	shell := "/bin/bash"
//
//	// exec remote shell
//	execUrl := coreClient.RESTClient().
//		Post().
//		Namespace(pod.Namespace).
//		Resource("pods").
//		Name(pod.Name).
//		SubResource("exec").
//		VersionedParams(&corev1.PodExecOptions{
//			Container: box.Template.SafeName(), // pod.Spec.Containers[0].Name
//			Command:   []string{shell},
//			Stdin:     true,
//			Stdout:    true,
//			Stderr:    true,
//			TTY:       tty.Raw,
//		}, scheme.ParameterCodec).
//		URL()
//
//	executor := exec.DefaultRemoteExecutor{}
//
//	box.OnExecCallback()
//
//	fn := func() error {
//		return executor.Execute(http.MethodPost, execUrl, client.kubeRestConfig, streamOptions.In, streamOptions.Out, streamOptions.ErrOut, tty.Raw, sizeQueue)
//	}
//	if err := tty.Safe(fn); err != nil {
//		return errors.Wrap(err, "terminal session closed")
//	}
//	return nil
//}
