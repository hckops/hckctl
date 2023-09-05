package kubernetes

import (
	"context"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	SingleReplica      = 1
	SidecarPrefix      = "sidecar-"
	LabelKubeName      = "app.kubernetes.io/name"
	LabelKubeInstance  = "app.kubernetes.io/instance"
	LabelKubeVersion   = "app.kubernetes.io/version"
	LabelKubeManagedBy = "app.kubernetes.io/managed-by"
)

type KubeClient struct {
	ctx            context.Context
	kubeRestConfig *rest.Config
	kubeClientSet  *kubernetes.Clientset
}

type KubeResource struct {
	Memory string
	Cpu    string
}

type DeploymentInfo struct {
	Namespace string
	Name      string
	Healthy   bool
	PodInfo   *PodInfo
}

type DeploymentDetails struct {
	Info        *DeploymentInfo
	Created     time.Time
	Annotations map[string]string
}

type PodInfo struct {
	Namespace     string
	PodName       string
	ContainerName string
	ImageName     string
	Env           []KubeEnv
	Resource      *KubeResource
}

type KubeEnv struct {
	Key   string
	Value string
}

type ServiceInfo struct {
	Namespace string
	Name      string
	Ports     []KubePort
}

type KubePort struct {
	Name string
	Port string
}
