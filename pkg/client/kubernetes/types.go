package kubernetes

import (
	"context"
	"io"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	SingleReplica = 1
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
	Created     string
	Annotations map[string]string
}

type PodInfo struct {
	Namespace string
	Id        string
	Name      string
	Env       map[string]string
}

type ServiceInfo struct {
	Namespace string
	Name      string
	Ports     []ServicePort
}

type ServicePort struct {
	Name string
	Port string
}

type DeploymentCreateOpts struct {
	Namespace             string
	Spec                  *appsv1.Deployment
	OnStatusEventCallback func(event string)
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
