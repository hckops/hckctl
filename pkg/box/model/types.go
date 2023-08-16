package model

import (
	"time"
)

type BoxInfo struct {
	Id      string
	Name    string
	Healthy bool // TODO BoxStatus healthy, offline, unknown, error, etc
}

type BoxDetails struct {
	Info         BoxInfo
	TemplateInfo *BoxTemplateInfo
	ProviderInfo *BoxProviderInfo
	Size         ResourceSize
	Env          []BoxEnv  // TODO map[string]BoxEnv
	Ports        []BoxPort // TODO map[string]BoxPort
	Created      time.Time
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

type BoxProviderInfo struct {
	Provider       BoxProvider
	DockerProvider *DockerProviderInfo
	KubeProvider   *KubeProviderInfo
}

type DockerProviderInfo struct {
	Network string
	Ip      string
}

type KubeProviderInfo struct {
	Namespace string
}
