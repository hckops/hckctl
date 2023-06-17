package kubernetes

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"strings"
)

//import (
//	"fmt"
//	appsv1 "k8s.io/api/apps/v1"
//	corev1 "k8s.io/api/core/v1"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	"k8s.io/apimachinery/pkg/watch"
//	applyv1 "k8s.io/client-go/applyconfigurations/core/v1"
//	"path/filepath"
//	"strings"
//
//	"github.com/pkg/errors"
//	"k8s.io/client-go/kubernetes"
//	"k8s.io/client-go/rest"
//	"k8s.io/client-go/tools/clientcmd"
//	"k8s.io/client-go/util/homedir"
//)

func NewKubeClient() (*KubeClient, error) {
	return &KubeClient{}, nil
}

func NewOutOfClusterClient(configPath string) (*rest.Config, *kubernetes.Clientset, error) {

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

func NewInClusterClient() (*rest.Config, *kubernetes.Clientset, error) {

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

func (client *KubeClient) Close() error {
	return errors.New("not implemented")
}

//func (client *KubeClient) ApplyTemplate(deploymentSpec *appsv1.Deployment, serviceSpec *corev1.Service) error {
//	coreClient := box.KubeClientSet.CoreV1()
//	appClient := box.KubeClientSet.AppsV1()
//
//	// apply namespace, see https://github.com/kubernetes/client-go/issues/1036
//	namespace, err := coreClient.Namespaces().Apply(box.ctx, applyv1.Namespace(box.ResourceOptions.Namespace), metav1.ApplyOptions{FieldManager: "application/apply-patch"})
//	if err != nil {
//		return errors.Wrapf(err, "error kube apply namespace: %s", box.ResourceOptions.Namespace)
//	}
//	box.OnSetupCallback(fmt.Sprintf("namespace %s successfully applied", namespace.Name))
//
//	// create deployment
//	deployment, err := appClient.Deployments(namespace.Name).Create(box.ctx, deploymentSpec, metav1.CreateOptions{})
//	if err != nil {
//		return errors.Wrapf(err, "error kube create deployment: %s", deployment.Name)
//	}
//	box.OnSetupCallback(fmt.Sprintf("deployment %s successfully created", deployment.Name))
//
//	// create service: can't create a service without ports
//	var service *corev1.Service
//	if box.Template.HasPorts() {
//		service, err = coreClient.Services(namespace.Name).Create(box.ctx, serviceSpec, metav1.CreateOptions{})
//		if err != nil {
//			return errors.Wrapf(err, "error kube create service: %s", service.Name)
//		}
//		box.OnSetupCallback(fmt.Sprintf("service %s successfully created", service.Name))
//	} else {
//		box.OnSetupCallback(fmt.Sprint("service not created"))
//	}
//
//	// blocks until the deployment is available, then stop watching
//	watcher, err := appClient.Deployments(namespace.Name).Watch(box.ctx, metav1.SingleObject(deployment.ObjectMeta))
//	if err != nil {
//		return errors.Wrapf(err, "error kube watch deployment: %s", deployment.Name)
//	}
//	for event := range watcher.ResultChan() {
//		switch event.Type {
//		case watch.Modified:
//			deploymentEvent := event.Object.(*appsv1.Deployment)
//
//			for _, condition := range deploymentEvent.Status.Conditions {
//				box.OnSetupCallback(fmt.Sprintf("watch kube event: type=%v, condition=%v", event.Type, condition.Message))
//
//				if condition.Type == appsv1.DeploymentAvailable &&
//					condition.Status == corev1.ConditionTrue {
//					watcher.Stop()
//				}
//			}
//
//		default:
//			return errors.Wrapf(err, "error kube event type: %s", event.Type)
//		}
//	}
//	return nil
//}
