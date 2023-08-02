package kubernetes

import (
	"io"

	appsv1 "k8s.io/api/apps/v1"
)

type ResourcesOpts struct {
	Namespace   string
	Name        string
	Annotations map[string]string
	Labels      map[string]string
	Ports       []KubePort
	PodInfo     *PodInfo
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
	OnTunnelStartCallback func()
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
