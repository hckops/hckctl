package model

import (
	"fmt"
)

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

type Image struct {
	Repository string
	Version    string
}

func (image *Image) Name() string {
	return fmt.Sprintf("%s:%s", image.Repository, image.ResolveVersion())
}

func (image *Image) ResolveVersion() string {
	var version string
	if image.Version == "" {
		version = "latest"
	} else {
		version = image.Version
	}
	return version
}
