package kubernetes

import (
	"context"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubeClient struct {
	ctx            context.Context
	kubeRestConfig *rest.Config
	kubeClientSet  *kubernetes.Clientset
}

type KubeClientConfig struct {
	ConfigPath string
	Namespace  string
	Resource   *KubeResource
}

type KubeResource struct {
	Memory string
	Cpu    string
}
