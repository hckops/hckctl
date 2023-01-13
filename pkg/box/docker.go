package box

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dchest/uniuri"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	privcommon "github.com/hckops/hckctl/internal/common"
	pubcommon "github.com/hckops/hckctl/pkg/common"
)

type DockerBox struct {
	ctx    context.Context
	loader *privcommon.Loader
	box    *pubcommon.BoxV1
}

func NewDockerBox(box *pubcommon.BoxV1) *DockerBox {
	return &DockerBox{
		ctx:    context.Background(),
		loader: privcommon.NewLoader(),
		box:    box,
	}
}

func (d *DockerBox) InitBox() {
	d.loader.Start(fmt.Sprintf("loading %s", d.box.Name))

	//containerName := d.box.GenerateName()

	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("error docker client: %v", err)
	}
	defer docker.Close()

	reader, err := docker.ImagePull(d.ctx, d.box.ImageName(), types.ImagePullOptions{})
	if err != nil {
		log.Fatalf("error image pull: %v", err)
	}
	defer reader.Close()

	// suppress output
	io.Copy(ioutil.Discard, reader)

	d.loader.Stop()
}

func InitDockerBoxOld(box *pubcommon.BoxV1) {

	ctx := context.Background()

	var imageVersion string

	imageName := fmt.Sprintf("%s:%s", box.Image.Repository, imageVersion)
	containerName := fmt.Sprintf("%s-%s", box.Name, uniuri.NewLen(5))

	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("error docker client: %v", err)
	}
	defer docker.Close()

	reader, err := docker.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		log.Fatalf("error image pull: %v", err)
	}
	defer reader.Close()

	// suppress output
	io.Copy(ioutil.Discard, reader)

	containerConfig := &container.Config{
		Image:       imageName,
		AttachStdin: true,
		//AttachStdout: true,
		//AttachStderr: true,
		OpenStdin: true,
		StdinOnce: true,
		//Tty:       true,
		// TODO
		ExposedPorts: nat.PortSet{
			nat.Port("5900/tcp"): {},
			nat.Port("6080/tcp"): {},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"5900/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "5900"}},
			"6080/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "6080"}},
		},
	}

	newContainer, err := docker.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig, // hostConfig
		nil,        // networkingConfig
		nil,        // platform
		containerName)
	if err != nil {
		log.Fatalf("error container create: %v", err)
	}

	containerId := newContainer.ID

	if err := docker.ContainerStart(ctx, containerId, types.ContainerStartOptions{}); err != nil {
		log.Fatalf("error container start: %v", err)
	}

	execCreateResponse, err := docker.ContainerExecCreate(ctx, containerId, types.ExecConfig{
		AttachStdout: true,
		AttachStdin:  true,
		AttachStderr: true,
		Detach:       false,
		Tty:          true,
		Cmd:          []string{"/bin/ash"},
	})
	if err != nil {
		log.Fatalf("error docker exec create: %v", err)
	}

	execAttachResponse, err := docker.ContainerExecAttach(ctx, execCreateResponse.ID, types.ExecStartCheck{
		Tty: true,
	})
	if err != nil {
		log.Fatalf("error docker exec attach: %v", err)
	}
	defer execAttachResponse.Close()

	closeChannel := func() {
		log.Printf("removing docker container: id=%s", containerId)

		if err := docker.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{Force: true}); err != nil {
			log.Fatalf("error docker remove: %v", err)
		}
	}

	var once sync.Once
	go func() {
		// use with TTY=false only, with TTY=true returns: "Unrecognized input header: 13"
		//_, err := stdcopy.StdCopy(os.Stdout, os.Stderr, execAttachResponse.Reader)

		// TTY
		_, err := io.Copy(os.Stdout, execAttachResponse.Reader)
		if err != nil {
			log.Fatalf("error copy docker->local: %v", err)
		}

		log.Printf("close docker->local")
		once.Do(closeChannel)
	}()

	go func() {
		_, err = io.Copy(execAttachResponse.Conn, os.Stdin)
		if err != nil {
			log.Fatalf("error copy local->docker: %v", err)
		}

		log.Printf("close local->docker")
		once.Do(closeChannel)
	}()

	// TODO CTRL+C should NOT exit
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	// signal.Notify(c, os.Interrupt)
	go func() {
		// for sig := range c {
		// 	// sig is a ^C, handle it
		// 	log.Printf("CTRL+C handler %v", sig)
		// }
		<-signalCh

		log.Printf("CTRL+C handler")
		once.Do(closeChannel)
		//os.Exit(0)
	}()

	statusCh, errCh := docker.ContainerWait(ctx, containerId, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Fatalf("error container wait: %v", err)
		}
		log.Printf("close container wait errCh")
	case status := <-statusCh:
		log.Printf("close container wait statusCh: %v", status.StatusCode)
	}
}
