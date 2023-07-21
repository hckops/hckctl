package model

import (
	"io"
	"os"

	"github.com/hckops/hckctl/pkg/event"
)

type BoxInfo struct {
	Id      string
	Name    string
	Healthy bool // TODO BoxStatus healthy, offline, unknown, error, etc
}

type BoxDetails struct {
	Info         BoxInfo
	Provider     BoxProvider
	Size         ResourceSize
	TemplateInfo *BoxTemplateInfo
	Env          []BoxEnv
	Ports        []BoxPort
	Docker       struct{ Network string }   // TODO missing
	Kube         struct{ Namespace string } // TODO missing
}

type BoxTemplateInfo struct {
	CachedTemplate *CachedTemplateInfo
	GitTemplate    *GitTemplateInfo
}

func (info *BoxTemplateInfo) IsCached() bool {
	return info.CachedTemplate != nil
}

type CachedTemplateInfo struct {
	Path string
}

type GitTemplateInfo struct {
	Url      string
	Revision string
	Commit   string
	Name     string
}

type BoxClientOptions struct {
	Provider   BoxProvider
	CommonOpts *CommonBoxOptions
	DockerOpts *DockerBoxOptions
	KubeOpts   *KubeBoxOptions
	CloudOpts  *CloudBoxOptions
}

type CommonBoxOptions struct {
	EventBus *event.EventBus
}

func NewCommonBoxOpts() *CommonBoxOptions {
	return &CommonBoxOptions{
		EventBus: event.NewEventBus(),
	}
}

type DockerBoxOptions struct {
	NetworkName          string
	IgnoreImagePullError bool
}

type KubeBoxOptions struct {
	InCluster  bool
	ConfigPath string
	Namespace  string
}

type CloudBoxOptions struct {
	Version  string
	Address  string
	Username string
	Token    string
}

type BoxStreams struct {
	In    io.ReadCloser
	Out   io.Writer
	Err   io.Writer
	IsTty bool // tty is false only for ssh tunnel
}

func NewDefaultStreams(tty bool) *BoxStreams {
	return &BoxStreams{
		In:    os.Stdin,
		Out:   os.Stdout,
		Err:   os.Stderr,
		IsTty: tty,
	}
}

type TemplateOptions struct {
	Template *BoxV1
	Size     ResourceSize
	Labels   BoxLabels
}

type TunnelOptions struct {
	Streams    *BoxStreams
	TunnelOnly bool
	NoTunnel   bool
}
