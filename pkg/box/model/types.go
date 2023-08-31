package model

import (
	"time"

	commonModel "github.com/hckops/hckctl/pkg/common/model"
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
	CachedTemplate *commonModel.CachedTemplateInfo
	GitTemplate    *commonModel.GitTemplateInfo
}

func (info *BoxTemplateInfo) IsCached() bool {
	return info.CachedTemplate != nil
}

type BoxProviderInfo struct {
	Provider       BoxProvider
	DockerProvider *commonModel.DockerProviderInfo
	KubeProvider   *commonModel.KubeProviderInfo
}
