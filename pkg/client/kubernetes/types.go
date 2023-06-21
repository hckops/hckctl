package kubernetes

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
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

type DeploymentCreateOpts struct {
	Namespace             string
	Spec                  *appsv1.Deployment
	OnStatusEventCallback func(event string)
}

type DeploymentInfo struct {
	Namespace      string
	DeploymentName string
	PodName        string // unique generated name
}

type PodPortForwardOpts struct {
	Namespace             string
	PodName               string
	Ports                 []string // format "LOCAL:REMOTE"
	OnTunnelErrorCallback func(error)
}
