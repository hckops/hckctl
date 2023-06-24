package box

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/command/common/flag"
)

func TestBoxProviders(t *testing.T) {
	assert.Equal(t, 3, len(boxProviders()))
	assert.Equal(t, "docker", boxProviders()[0].String())
	assert.Equal(t, "kube", boxProviders()[1].String())
	assert.Equal(t, "cloud", boxProviders()[2].String())
}

func TestToBoxProvider(t *testing.T) {
	docker, err := toBoxProvider(flag.DockerProviderFlag)
	assert.NoError(t, err)
	assert.Equal(t, model.Docker, docker)

	kube, err := toBoxProvider(flag.KubeProviderFlag)
	assert.NoError(t, err)
	assert.Equal(t, model.Kubernetes, kube)

	cloud, err := toBoxProvider(flag.CloudProviderFlag)
	assert.NoError(t, err)
	assert.Equal(t, model.Cloud, cloud)

	_, err = toBoxProvider(flag.UnknownProviderFlag)
	assert.EqualError(t, err, "invalid provider")
}

func TestBoxProviderIds(t *testing.T) {
	assert.Equal(t, 3, len(boxProviderIds()))

	assert.Equal(t, []string{"docker"}, boxProviderIds()[flag.DockerProviderFlag])
	assert.Equal(t, []string{"kube", "k8s", "kubernetes"}, boxProviderIds()[flag.KubeProviderFlag])
	assert.Equal(t, []string{"cloud"}, boxProviderIds()[flag.CloudProviderFlag])
}

func TestValidateBoxProviderConfig(t *testing.T) {
	var boxProviderFlag flag.ProviderFlag
	boxProviderFlag = flag.UnknownProviderFlag
	boxProvider, err := validateBoxProvider("docker", &boxProviderFlag)

	assert.NoError(t, err)
	assert.Equal(t, "docker", boxProvider.String())
}

func TestValidateBoxProviderFlag(t *testing.T) {
	var boxProviderFlag flag.ProviderFlag
	boxProviderFlag = flag.KubeProviderFlag
	boxProvider, err := validateBoxProvider("docker", &boxProviderFlag)

	assert.NoError(t, err)
	assert.Equal(t, "kube", boxProvider.String())
}
