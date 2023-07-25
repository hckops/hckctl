package flag

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProviderFlag(t *testing.T) {
	assert.Equal(t, 3, len(allProviderIds))
	assert.Equal(t, []string{"docker"}, allProviderIds[DockerProviderFlag])
	assert.Equal(t, []string{"kube", "k8s", "kubernetes"}, allProviderIds[KubeProviderFlag])
	assert.Equal(t, []string{"cloud"}, allProviderIds[CloudProviderFlag])
}

func TestProviderString(t *testing.T) {
	assert.Equal(t, "unknown", UnknownProviderFlag.String())
	assert.Equal(t, "docker", DockerProviderFlag.String())
	assert.Equal(t, "kube", KubeProviderFlag.String())
	assert.Equal(t, "cloud", CloudProviderFlag.String())
}

func TestProviderIds(t *testing.T) {
	providerFlags := []ProviderFlag{CloudProviderFlag, KubeProviderFlag}
	providerIds := ProviderIds(providerFlags)

	assert.Equal(t, 2, len(providerIds))
	assert.Equal(t, []string{"cloud"}, providerIds[CloudProviderFlag])
	assert.Equal(t, []string{"kube", "k8s", "kubernetes"}, providerIds[KubeProviderFlag])
}

func TestProviderValues(t *testing.T) {
	providerValues := ProviderValues(allProviderIds)
	expected := []string{"cloud", "docker", "k8s", "kube", "kubernetes"}

	assert.Equal(t, 5, len(providerValues))
	assert.Equal(t, expected, providerValues)

	// the return value of SearchStrings is the index to insert x if x is not present
	assert.Equal(t, 5, sort.SearchStrings(providerValues, "unknown"))
}

func TestExistProvider(t *testing.T) {
	providerFlag, err := ExistProvider(allProviderIds, "cloud")
	assert.NoError(t, err)
	assert.Equal(t, CloudProviderFlag, providerFlag)

	_, err = ExistProvider(allProviderIds, "unknown")
	assert.EqualError(t, err, "invalid provider")
}
