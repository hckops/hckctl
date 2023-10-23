package kubernetes

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	applyv1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	app "k8s.io/client-go/kubernetes/typed/apps/v1"
	batch "k8s.io/client-go/kubernetes/typed/batch/v1"
	core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/transport/spdy"
	"k8s.io/client-go/util/homedir"
	"k8s.io/kubectl/pkg/cmd/exec"
	"k8s.io/kubectl/pkg/scheme"

	"github.com/hckops/hckctl/pkg/util"
)

func NewKubeClient(inCluster bool, configPath string) (*KubeClient, error) {
	ctx := context.Background()
	if inCluster {
		return newInClusterKubeClient(ctx)
	} else {
		return newOutOfClusterKubeClient(ctx, configPath)
	}
}

func newOutOfClusterKubeClient(ctx context.Context, configPath string) (*KubeClient, error) {
	restConfig, clientSet, err := NewOutOfClusterKubeConfig(configPath)
	if err != nil {
		return nil, err
	}
	return &KubeClient{
		ctx:            ctx,
		kubeRestConfig: restConfig,
		kubeClientSet:  clientSet,
	}, nil
}

func NewOutOfClusterKubeConfig(configPath string) (*rest.Config, *kubernetes.Clientset, error) {

	kubeConfig := NormalizeKubeConfig(configPath)

	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error out-of-cluster restConfig")
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error out-of-cluster clientSet")
	}

	return restConfig, clientSet, nil
}

func NormalizeKubeConfig(configPath string) string {
	if strings.TrimSpace(configPath) == "" {
		return filepath.Join(homedir.HomeDir(), ".kube", "config")
	} else {
		// absolute path
		return configPath
	}
}

func newInClusterKubeClient(ctx context.Context) (*KubeClient, error) {
	restConfig, clientSet, err := NewInClusterKubeConfig()
	if err != nil {
		return nil, err
	}
	return &KubeClient{
		ctx:            ctx,
		kubeRestConfig: restConfig,
		kubeClientSet:  clientSet,
	}, nil
}

func NewInClusterKubeConfig() (*rest.Config, *kubernetes.Clientset, error) {

	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, nil, errors.Wrap(err, "error in-cluster restConfig")
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error in-cluster clientSet")
	}

	return restConfig, clientSet, nil
}

func (client *KubeClient) Close() error {
	return errors.New("not implemented")
}

func (client *KubeClient) RestApi() *rest.Config {
	return client.kubeRestConfig
}

func (client *KubeClient) CoreApi() core.CoreV1Interface {
	return client.kubeClientSet.CoreV1()
}

func (client *KubeClient) AppApi() app.AppsV1Interface {
	return client.kubeClientSet.AppsV1()
}

func (client *KubeClient) BatchApi() batch.BatchV1Interface {
	return client.kubeClientSet.BatchV1()
}

func (client *KubeClient) NamespaceApply(name string) error {

	// https://github.com/kubernetes/client-go/issues/1036
	_, err := client.CoreApi().Namespaces().Apply(client.ctx, applyv1.Namespace(name), metav1.ApplyOptions{FieldManager: "application/apply-patch"})
	if err != nil {
		return errors.Wrapf(err, "error namespace apply: name=%s", name)
	}
	return nil
}

func (client *KubeClient) NamespaceDelete(name string) error {

	if err := client.CoreApi().Namespaces().Delete(client.ctx, name, metav1.DeleteOptions{}); err != nil {
		return errors.Wrapf(err, "error namespace delete: name=%s", name)
	}
	return nil
}

func (client *KubeClient) DeploymentCreate(opts *DeploymentCreateOpts) error {

	deployment, err := client.AppApi().Deployments(opts.Namespace).Create(client.ctx, opts.Spec, metav1.CreateOptions{})
	if err != nil {
		return errors.Wrapf(err, "error deployment create: namespace=%s name=%s", opts.Namespace, opts.Spec.Name)
	}

	// blocks until the deployment is available, then stop watching
	watcher, err := client.AppApi().Deployments(opts.Namespace).Watch(client.ctx, metav1.SingleObject(deployment.ObjectMeta))
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

func (client *KubeClient) DeploymentList(namespace string, namePrefix string, labelSelector string) ([]DeploymentInfo, error) {

	deployments, err := client.AppApi().Deployments(namespace).List(client.ctx, metav1.ListOptions{
		// comma separated values with format <LABEL_KEY>=<SANITIZED_LABEL_VALUE>
		LabelSelector: labelSelector,
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

		podInfo, err := client.PodDescribeFromDeployment(&deployment)
		if err != nil {
			// skip invalid pod
			continue
		}

		result = append(result, newDeploymentInfo(&deployment, podInfo))
	}
	return result, nil
}

func newDeploymentInfo(deployment *appsv1.Deployment, podInfo *PodInfo) DeploymentInfo {
	return DeploymentInfo{
		Namespace: deployment.Namespace,
		Name:      deployment.Name,
		Healthy:   isDeploymentHealthy(deployment.Status),
		PodInfo:   podInfo,
	}
}

func isDeploymentHealthy(status appsv1.DeploymentStatus) bool {
	// all conditions must be true to be healthy
	var healthy bool
	for _, condition := range status.Conditions {
		if condition.Status == corev1.ConditionTrue {
			healthy = true
		} else {
			healthy = false
			// mark as unhealthy if at least one condition is false e.g. stuck due to lack of resources
			break
		}
	}
	return healthy
}

func (client *KubeClient) DeploymentDescribe(namespace string, name string) (*DeploymentDetails, error) {

	deployment, err := client.AppApi().Deployments(namespace).Get(client.ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "error deployment describe: namespace=%s name=%s", namespace, name)
	}

	podInfo, err := client.PodDescribeFromDeployment(deployment)
	if err != nil {
		return nil, err
	}

	return newDeploymentDetails(deployment, podInfo), nil
}

func newDeploymentDetails(deployment *appsv1.Deployment, podInfo *PodInfo) *DeploymentDetails {
	deploymentInfo := newDeploymentInfo(deployment, podInfo)
	return &DeploymentDetails{
		Info:        &deploymentInfo,
		Created:     deployment.CreationTimestamp.Time.UTC(),
		Annotations: deployment.Annotations,
	}
}

func (client *KubeClient) DeploymentDelete(namespace string, name string) error {

	err := client.AppApi().Deployments(namespace).Delete(client.ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return errors.Wrapf(err, "error deployment delete: namespace=%s name=%s", namespace, name)
	}
	return nil
}

func (client *KubeClient) ServiceCreate(namespace string, spec *corev1.Service) error {

	_, err := client.CoreApi().Services(namespace).Create(client.ctx, spec, metav1.CreateOptions{})
	if err != nil {
		return errors.Wrapf(err, "error service create: namespace=%s name=%s", namespace, spec.Name)
	}
	return nil
}

func (client *KubeClient) ServiceDescribe(namespace string, name string) (*ServiceInfo, error) {

	service, err := client.CoreApi().Services(namespace).Get(client.ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "error service describe: namespace=%s name=%s", namespace, name)
	}

	return newServiceInfo(service), nil
}

func newServiceInfo(service *corev1.Service) *ServiceInfo {
	var ports []KubePort
	for _, port := range service.Spec.Ports {
		ports = append(ports, KubePort{Name: port.Name, Port: strconv.Itoa(int(port.Port))})
	}

	return &ServiceInfo{
		Namespace: service.Namespace,
		Name:      service.Name,
		Ports:     ports,
	}
}

func (client *KubeClient) ServiceDelete(namespace string, name string) error {

	if err := client.CoreApi().Services(namespace).Delete(client.ctx, name, metav1.DeleteOptions{}); err != nil {
		return errors.Wrapf(err, "error service delete: namespace=%s name=%s", namespace, name)
	}
	return nil
}

func (client *KubeClient) PodDescribeFromDeployment(deployment *appsv1.Deployment) (*PodInfo, error) {
	labelSet := labels.Set(deployment.Spec.Selector.MatchLabels)
	listOptions := metav1.ListOptions{
		LabelSelector: labelSet.AsSelector().String(),
	}

	return client.podDescribe(deployment.Namespace, listOptions)
}

func (client *KubeClient) podDescribe(namespace string, listOptions metav1.ListOptions) (*PodInfo, error) {

	pods, err := client.CoreApi().Pods(namespace).List(client.ctx, listOptions)
	if err != nil {
		return nil, errors.Wrapf(err, "error pod describe: namespace=%s labels=%v", namespace, listOptions.LabelSelector)
	}

	return newPodInfo(namespace, pods)
}

func newPodInfo(namespace string, pods *corev1.PodList) (*PodInfo, error) {
	if len(pods.Items) != SingleReplica {
		return nil, fmt.Errorf("found %d pods, expected only 1 pod: namespace=%s", len(pods.Items), namespace)
	}

	podItem := pods.Items[0]

	var containers []corev1.Container
	for _, c := range podItem.Spec.Containers {
		// exclude injected sidecar containers
		if !strings.HasPrefix(c.Name, SidecarPrefix) {
			containers = append(containers, c)
		}
	}
	if len(containers) != 1 {
		return nil, fmt.Errorf("found %d containers, expected only 1 container for pod: namespace=%s", len(podItem.Spec.Containers), namespace)
	}
	containerItem := containers[0]

	var envs []KubeEnv
	for _, e := range containerItem.Env {
		envs = append(envs, KubeEnv{Key: e.Name, Value: e.Value})
	}

	return &PodInfo{
		Namespace:     podItem.Namespace,
		PodName:       podItem.Name, // pod.Name + unique generated suffix
		ContainerName: containerItem.Name,
		ImageName:     containerItem.Image, // <REPOSITORY>/<NAME>:<VERSION>
		Env:           envs,
		Arguments:     []string{}, // ignore containerItem.Command and containerItem.Args
		Resource: &KubeResource{
			Memory: containerItem.Resources.Requests.Memory().String(),
			Cpu:    containerItem.Resources.Requests.Cpu().String(),
		},
	}, nil
}

func (client *KubeClient) PodPortForward(opts *PodPortForwardOpts) error {

	restRequest := client.CoreApi().RESTClient().
		Post().
		Resource("pods").
		Namespace(opts.Namespace).
		Name(opts.PodName).
		SubResource("portforward")

	transport, upgrader, err := spdy.RoundTripperFor(client.RestApi())
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

	opts.OnTunnelStartCallback()

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

func (client *KubeClient) newRestRequestExec(opts *PodExecOpts, isTty bool) *rest.Request {
	return client.CoreApi().RESTClient().
		Post().
		Namespace(opts.Namespace).
		Resource("pods").
		Name(opts.PodName).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: opts.ContainerName,
			Command:   opts.Commands,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       isTty,
		}, scheme.ParameterCodec)
}

func (client *KubeClient) PodExecShell(opts *PodExecOpts) error {

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
	execUrl := client.newRestRequestExec(opts, true).URL()
	executor := exec.DefaultRemoteExecutor{}

	opts.OnExecCallback()

	fn := func() error {
		isTty := tty.Raw && opts.IsTty
		return executor.Execute(http.MethodPost, execUrl, client.RestApi(), streamOptions.In, streamOptions.Out, streamOptions.ErrOut, isTty, sizeQueue)
	}
	if err := tty.Safe(fn); err != nil {
		return errors.Wrap(err, "terminal session closed")
	}
	return nil
}

func (client *KubeClient) PodExecCommand(opts *PodExecOpts) error {
	isTty := false
	execUrl := client.newRestRequestExec(opts, isTty).URL()
	return (&exec.DefaultRemoteExecutor{}).Execute(http.MethodPost, execUrl, client.RestApi(), opts.InStream, opts.OutStream, opts.ErrStream, isTty, nil)
}

func (client *KubeClient) podLogsStream(opts *PodLogsOpts) (io.ReadCloser, error) {

	logOptions := &corev1.PodLogOptions{
		Container: opts.ContainerName,
		Follow:    true,
	}
	outStream, err := client.CoreApi().
		Pods(opts.Namespace).
		GetLogs(opts.PodName, logOptions).
		Stream(client.ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error pod logs stream")
	}
	return outStream, nil
}

func (client *KubeClient) PodLogs(opts *PodLogsOpts) error {

	outStream, err := client.podLogsStream(opts)
	if err != nil {
		return err
	}
	defer outStream.Close()

	// blocks until the stream is finished
	if _, err = io.Copy(opts.OutStream, outStream); err != nil {
		return errors.Wrapf(err, "error pod logs std copy")
	}
	return nil
}

func (client *KubeClient) PodLogsTee(opts *PodLogsOpts, logFileName string) error {

	outStream, err := client.podLogsStream(opts)
	if err != nil {
		return err
	}
	defer outStream.Close()

	logFile, err := util.OpenFile(logFileName)
	if err != nil {
		return errors.Wrap(err, "error pod logs file")
	}
	multiWriter := io.MultiWriter(opts.OutStream, logFile)
	defer logFile.Close()

	if _, err = io.Copy(multiWriter, outStream); err != nil {
		return errors.Wrapf(err, "error pod logs tee copy")
	}
	return nil
}

func (client *KubeClient) CopyToPod(opts *CopyPodOpts) error {

	if util.PathNotExist(opts.LocalPath) {
		return fmt.Errorf("error copy invalid localPath=%s", opts.LocalPath)
	}

	reader, writer := io.Pipe()
	defer reader.Close()

	// create archive
	go func() {
		defer writer.Close()

		// TODO unhandled error
		cpMakeTar(opts.LocalPath, opts.RemotePath, writer)
	}()

	// upload and extract archive
	execArchive := &PodExecOpts{
		Namespace:      opts.Namespace,
		PodName:        opts.PodName,
		ContainerName:  opts.ContainerName,
		Commands:       []string{"tar", "-xmf", "-", "-C", filepath.Dir(opts.RemotePath)},
		InStream:       reader, // input stream reader
		OutStream:      io.Discard,
		ErrStream:      io.Discard,
		IsTty:          false,
		OnExecCallback: func() {},
	}
	if err := client.PodExecCommand(execArchive); err != nil {
		return errors.Wrapf(err, "error copy archive")
	}
	return nil
}

func buildJobLabelSelector(job *batchv1.Job) (metav1.ListOptions, error) {
	labelMap, err := metav1.LabelSelectorAsMap(job.Spec.Selector)
	if err != nil {
		return metav1.ListOptions{}, errors.Wrapf(err, "error pod describe from job")
	}
	listOptions := metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labelMap).String(),
	}
	return listOptions, nil
}

func (client *KubeClient) JobCreate(opts *JobCreateOpts) error {

	job, err := client.BatchApi().Jobs(opts.Namespace).Create(client.ctx, opts.Spec, metav1.CreateOptions{})
	if err != nil {
		return errors.Wrapf(err, "error job create: namespace=%s name=%s", opts.Namespace, opts.Spec.Name)
	}

	// TODO move in the commands
	if opts.CaptureInterrupt {
		util.InterruptHandler(func() {
			opts.OnContainerInterruptCallback(job.Name)
		})
	}

	// blocks until the job is ready, then stop watching
	listOptions, err := buildJobLabelSelector(job)
	if err != nil {
		return err
	}
	watcher, err := client.CoreApi().Pods(opts.Namespace).Watch(client.ctx, listOptions)
	if err != nil {
		return errors.Wrapf(err, "error job watch: namespace=%s name=%s", opts.Namespace, job.Name)
	}
	for event := range watcher.ResultChan() {
		podEvent, ok := event.Object.(*corev1.Pod)
		if !ok {
			// ignore unrelated events
			continue
		}

		opts.OnStatusEventCallback(fmt.Sprintf("watch job event: namespace=%s name=%s type=%v phase=%v",
			opts.Namespace, job.Name, event.Type, podEvent.Status.Phase))

		if podEvent.Status.Phase != corev1.PodPending {
			watcher.Stop()
		}
	}
	return nil
}

func (client *KubeClient) JobDescribe(namespace string, name string) (*PodInfo, error) {

	job, err := client.BatchApi().Jobs(namespace).Get(client.ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "error job describe: namespace=%s name=%s", namespace, name)
	}

	listOptions, err := buildJobLabelSelector(job)
	if err != nil {
		return nil, err
	}

	return client.podDescribe(job.Namespace, listOptions)
}

func (client *KubeClient) JobDelete(namespace string, name string) error {

	// delete job and all pods
	backgroundDeletion := metav1.DeletePropagationBackground
	err := client.BatchApi().Jobs(namespace).Delete(client.ctx, name, metav1.DeleteOptions{
		PropagationPolicy: &backgroundDeletion,
	})
	if err != nil {
		return errors.Wrapf(err, "error job delete: namespace=%s name=%s", namespace, name)
	}

	return nil
}

func (client *KubeClient) SecretCreate(namespace string, spec *corev1.Secret) error {

	_, err := client.CoreApi().Secrets(namespace).Create(client.ctx, spec, metav1.CreateOptions{})
	if err != nil {
		return errors.Wrapf(err, "error secret create: namespace=%s name=%s", namespace, spec.Name)
	}
	return nil
}

func (client *KubeClient) SecretDelete(namespace string, name string) (bool, error) {

	_, err := client.CoreApi().Secrets(namespace).Get(client.ctx, name, metav1.GetOptions{})
	if err != nil {
		// it means the secret doesn't exist
		return false, nil
	}

	if err := client.CoreApi().Secrets(namespace).Delete(client.ctx, name, metav1.DeleteOptions{}); err != nil {
		return false, errors.Wrapf(err, "error secret delete: namespace=%s name=%s", namespace, name)
	}
	return true, nil
}
