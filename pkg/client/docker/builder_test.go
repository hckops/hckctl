package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

func TestBuildContainerConfig(t *testing.T) {
	envs := []ContainerEnv{
		{Key: "TTYD_USERNAME", Value: "username"},
		{Key: "TTYD_PASSWORD", Value: "password"},
	}
	ports := []ContainerPort{
		{Local: "123", Remote: "123"},
		{Local: "456", Remote: "789"},
	}
	expected := &container.Config{
		Hostname:     "myContainerName",
		Image:        "myImageName",
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		StdinOnce:    true,
		Tty:          true,
		Entrypoint:   []string{"xyz", "abc"},
		Cmd:          []string{"foo", "bar"},
		Env: []string{
			"TTYD_USERNAME=username",
			"TTYD_PASSWORD=password",
		},
		ExposedPorts: nat.PortSet{
			"123/tcp": struct{}{},
			"789/tcp": struct{}{},
		},
		Labels: map[string]string{
			"a.b.c": "hello",
			"x.y.z": "world",
		},
	}
	opts := &ContainerConfigOpts{
		ImageName:  "myImageName",
		Hostname:   "myContainerName",
		Env:        envs,
		Ports:      ports,
		Tty:        true,
		Entrypoint: []string{"xyz", "abc"},
		Cmd:        []string{"foo", "bar"},
		Labels: map[string]string{
			"a.b.c": "hello",
			"x.y.z": "world",
		},
	}

	result, err := BuildContainerConfig(opts)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestBuildHostConfig(t *testing.T) {
	// 1024 is the first user port available https://superuser.com/questions/1631280/what-exactly-are-user-ports
	ports := []ContainerPort{
		{Local: "1024", Remote: "1024"},
	}
	expected := &container.HostConfig{
		NetworkMode: "myNetworkMode",
		PortBindings: nat.PortMap{
			"1024/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "1024"}},
		},
		Mounts: []mount.Mount{{
			Type:   mount.TypeBind,
			Source: "/tmp/hck/share",
			Target: "/hck/share",
		}},
	}
	opts := &ContainerHostConfigOpts{
		NetworkMode: "myNetworkMode",
		PortConfig: &ContainerPortConfigOpts{
			Ports:              ports,
			OnPortBindCallback: func(port ContainerPort) {},
		},
		Volumes: []ContainerVolume{
			{HostDir: "/tmp/hck/share", ContainerDir: "/hck/share"},
		},
	}

	result, err := BuildHostConfig(opts)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestBuildVpnHostConfig(t *testing.T) {
	ports := []ContainerPort{
		{Local: "1024", Remote: "1024"},
	}
	expected := &container.HostConfig{
		CapAdd:  []string{"NET_ADMIN"},
		Sysctls: map[string]string{"net.ipv6.conf.all.disable_ipv6": "0"},
		Resources: container.Resources{
			Devices: []container.DeviceMapping{
				{
					PathOnHost:        "/dev/net/tun",
					PathInContainer:   "/dev/net/tun",
					CgroupPermissions: "rwm",
				},
			},
		},
		PortBindings: nat.PortMap{
			"1024/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "1024"}},
		},
	}
	result, err := BuildVpnHostConfig(&ContainerPortConfigOpts{
		Ports:              ports,
		OnPortBindCallback: func(port ContainerPort) {},
	})
	assert.Equal(t, expected, result)
	assert.NoError(t, err)
}

func TestDefaultNetworkMode(t *testing.T) {
	assert.Equal(t, "default", DefaultNetworkMode())
}

func TestContainerNetworkMode(t *testing.T) {
	assert.Equal(t, "container:myIdOrName", ContainerNetworkMode("myIdOrName"))
}

func TestBuildNetworkingConfig(t *testing.T) {
	expected := &network.NetworkingConfig{EndpointsConfig: map[string]*network.EndpointSettings{"myNetwork": {NetworkID: "123"}}}

	result := BuildNetworkingConfig("myNetwork", "123")
	assert.Equal(t, expected, result)
}

func TestDefaultPlatform(t *testing.T) {
	expected := &ocispec.Platform{Architecture: "amd64", OS: "linux"}
	assert.Equal(t, expected, DefaultPlatform())
}
