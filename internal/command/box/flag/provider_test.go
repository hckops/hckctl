package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hckops/hckctl/internal/command/common/flag"
	"github.com/hckops/hckctl/pkg/box/model"
)

func TestBoxProviders(t *testing.T) {
	assert.Equal(t, 3, len(BoxProviders()))
	assert.Equal(t, "docker", BoxProviders()[0].String())
	assert.Equal(t, "kube", BoxProviders()[1].String())
	assert.Equal(t, "cloud", BoxProviders()[2].String())
}

func TestToBoxProvider(t *testing.T) {
	docker, err := ToBoxProvider(flag.DockerProviderFlag)
	assert.NoError(t, err)
	assert.Equal(t, model.Docker, docker)

	kube, err := ToBoxProvider(flag.KubeProviderFlag)
	assert.NoError(t, err)
	assert.Equal(t, model.Kubernetes, kube)

	cloud, err := ToBoxProvider(flag.CloudProviderFlag)
	assert.NoError(t, err)
	assert.Equal(t, model.Cloud, cloud)

	_, err = ToBoxProvider(flag.UnknownProviderFlag)
	assert.EqualError(t, err, "invalid provider")
}

func TestBoxProviderIds(t *testing.T) {
	assert.Equal(t, 3, len(boxProviderIds()))

	assert.Equal(t, []string{"docker"}, boxProviderIds()[flag.DockerProviderFlag])
	assert.Equal(t, []string{"kube"}, boxProviderIds()[flag.KubeProviderFlag])
	assert.Equal(t, []string{"cloud"}, boxProviderIds()[flag.CloudProviderFlag])
}

func TestValidateBoxProviderConfig(t *testing.T) {
	var boxProviderFlag flag.ProviderFlag
	boxProviderFlag = flag.UnknownProviderFlag
	boxProvider, err := ValidateBoxProvider("docker", &boxProviderFlag)

	assert.NoError(t, err)
	assert.Equal(t, "docker", boxProvider.String())
}

func TestValidateBoxProviderFlag(t *testing.T) {
	var boxProviderFlag flag.ProviderFlag
	boxProviderFlag = flag.KubeProviderFlag
	boxProvider, err := ValidateBoxProvider("docker", &boxProviderFlag)

	assert.NoError(t, err)
	assert.Equal(t, "kube", boxProvider.String())
}
