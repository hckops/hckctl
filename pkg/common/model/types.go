package model

type DockerProviderInfo struct {
	Network string
	Ip      string
}

type KubeProviderInfo struct {
	Namespace string
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
