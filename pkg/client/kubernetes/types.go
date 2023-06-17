package kubernetes

import (
	"context"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubeClient struct {
	ctx            context.Context
	KubeRestConfig *rest.Config
	KubeClientSet  *kubernetes.Clientset
}

type ResourceOptions struct {
	Namespace string
	Memory    string
	Cpu       string
}
