package model

import (
	"io"
	"os"
)

type DockerOptions struct {
	NetworkName          string
	IgnoreImagePullError bool
}

type KubeOptions struct {
	InCluster  bool
	ConfigPath string
	Namespace  string
}

type CloudOptions struct {
	Version  string
	Address  string
	Username string
	Token    string
}

type StreamOptions struct {
	In    io.ReadCloser
	Out   io.Writer
	Err   io.Writer
	IsTty bool // tty is false for ssh tunnel or logs
}

func NewStdStreamOpts(tty bool) *StreamOptions {
	return &StreamOptions{
		In:    os.Stdin,
		Out:   os.Stdout,
		Err:   os.Stderr,
		IsTty: tty,
	}
}

type SidecarVpnInjectOpts struct {
	MainContainerId string
	NetworkVpn      *NetworkVpnInfo
}

type SidecarShareInjectOpts struct {
	MainContainerName string
	ShareDir          *ShareDirInfo
}

type SidecarShareUploadOpts struct {
	Namespace string
	PodName   string
	ShareDir  *ShareDirInfo
}
