package kubernetes

import (
	"context"
	"io"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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

type DeploymentCreateOpts struct {
	Namespace             string
	Spec                  *appsv1.Deployment
	OnStatusEventCallback func(event string)
}

type DeploymentInfo struct {
	Namespace      string
	DeploymentName string
	PodInfo        *PodInfo
	Healthy        bool
	Labels         map[string]string
	ServicePorts   []ServicePort
}

type PodInfo struct {
	Id   string
	Name string
}

type ServicePort struct {
	Name string
	Port string
}

type PodPortForwardOpts struct {
	Namespace             string
	PodId                 string
	Ports                 []string // format "LOCAL:REMOTE"
	IsWait                bool
	OnTunnelErrorCallback func(error)
}

type PodExecOpts struct {
	Namespace      string
	PodName        string
	PodId          string
	Shell          string
	InStream       io.ReadCloser
	OutStream      io.Writer
	ErrStream      io.Writer
	IsTty          bool
	OnExecCallback func()
}
