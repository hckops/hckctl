package box

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"golang.org/x/crypto/ssh/terminal"

	privcommon "github.com/hckops/hckctl/internal/common"
	pubcommon "github.com/hckops/hckctl/pkg/common"
)

type DockerBox struct {
	ctx          context.Context
	dockerClient *client.Client
	loader       *privcommon.Loader
	boxTemplate  *pubcommon.BoxV1
}

func NewDockerBox(box *pubcommon.BoxV1) *DockerBox {

	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("error docker client: %v", err)
	}
	defer docker.Close()

	return &DockerBox{
		ctx:          context.Background(),
		dockerClient: docker,
		loader:       privcommon.NewLoader(),
		boxTemplate:  box,
	}
}

func (d *DockerBox) InitBox() {
	d.loader.Start(fmt.Sprintf("loading %s", d.boxTemplate.Name))

	reader, err := d.dockerClient.ImagePull(d.ctx, d.boxTemplate.ImageName(), types.ImagePullOptions{})
	if err != nil {
		log.Fatalf("error image pull: %v", err)
	}
	defer reader.Close()

	d.loader.Refresh(fmt.Sprintf("pulling %s", d.boxTemplate.ImageName()))
	// suppress output
	io.Copy(ioutil.Discard, reader)

	containerName := d.boxTemplate.GenerateName()

	d.loader.Refresh(fmt.Sprintf("starting %s", containerName))

	// TODO is port is busy start in port+1 ? or attach to existing ?
	newContainer, err := d.dockerClient.ContainerCreate(
		d.ctx,
		buildContainerConfig(d.boxTemplate), // containerConfig
		buildHostConfig(d.boxTemplate),      // hostConfig
		nil,                                 // networkingConfig
		nil,                                 // platform
		containerName)
	if err != nil {
		log.Fatalf("error container create: %v", err)
	}

	containerId := newContainer.ID

	// if err := d.dockerClient.ContainerStart(d.ctx, containerId, types.ContainerStartOptions{}); err != nil {
	// 	log.Fatalf("error container start: %v", err)
	// }

	//d.todoAttach(containerId)
	// TODO tty false for tunnel only
	d.todoExecAttach(containerId, true)
}

func buildContainerConfig(boxTemplate *pubcommon.BoxV1) *container.Config {

	// TODO iterate over ports

	return &container.Config{
		Image:        boxTemplate.ImageName(),
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		StdinOnce:    true,
		Tty:          true,
		// TODO
		ExposedPorts: nat.PortSet{
			nat.Port("5900/tcp"): {},
			nat.Port("6080/tcp"): {},
		},
	}
}

func buildHostConfig(boxTemplate *pubcommon.BoxV1) *container.HostConfig {

	// TODO iterate over ports

	return &container.HostConfig{
		PortBindings: nat.PortMap{
			"5900/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "5900"}},
			"6080/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "6080"}},
		},
	}
}

func (d *DockerBox) todoAttach(containerId string) {

	attach, err := d.dockerClient.ContainerAttach(
		d.ctx,
		containerId,
		types.ContainerAttachOptions{
			Stream: true,
			Stdin:  true,
			Stdout: true,
			Stderr: true,
			Logs:   true,
		},
	)
	if err != nil {
		log.Fatalf("error container attach: %v", err)
	}

	closeCallback := func() {
		err = d.dockerClient.ContainerRemove(d.ctx, containerId, types.ContainerRemoveOptions{})
		if err != nil {
			log.Fatalf("error container remove: %v", err)
		}
	}

	var once sync.Once
	go func() {
		_, err = io.Copy(os.Stdout, attach.Reader)
		if err != nil {
			log.Fatalf("error copy docker->stdout: %v", err)
		}
		once.Do(closeCallback)
	}()
	// go func() {
	// 	_, err = io.Copy(os.Stderr, attach.Reader)
	// 	if err != nil {
	// 		log.Fatalf("error copy docker->stderr: %v", err)
	// 	}
	// 	once.Do(closeCallback)
	// }()
	go func() {
		_, _ = io.Copy(attach.Conn, os.Stdin)
		if err != nil {
			log.Fatalf("error copy stdin->docker: %v", err)
		}
		once.Do(closeCallback)
	}()

	if err := d.dockerClient.ContainerStart(d.ctx, containerId, types.ContainerStartOptions{}); err != nil {
		log.Fatalf("error container start: %v", err)
	}

	fd := int(os.Stdin.Fd())
	var oldState *terminal.State
	if terminal.IsTerminal(fd) {
		oldState, err = terminal.MakeRaw(fd)
		if err != nil {
			// print error
		}
		defer terminal.Restore(fd, oldState)
	}

	// clear
	//fmt.Print("\033[H\033[2J")
	d.loader.Stop()
	fmt.Print("\033[F")

	statusCh, errCh := d.dockerClient.ContainerWait(d.ctx, containerId, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Fatalf("error container wait: %v", err)
		}
		//log.Printf("close container wait errCh")
	case <-statusCh:
		//log.Printf("close container wait statusCh: %v", status.StatusCode)
	}

	terminal.Restore(fd, oldState)
}

// LAB: configure remove/keepalive timeout, override name, shell, installed packages etc.
func (d *DockerBox) todoExecAttach(containerId string, tty bool) {

	if err := d.dockerClient.ContainerStart(d.ctx, containerId, types.ContainerStartOptions{}); err != nil {
		log.Fatalf("error container start: %v", err)
	}

	// TODO always bash
	execCreateResponse, err := d.dockerClient.ContainerExecCreate(d.ctx, containerId, types.ExecConfig{
		AttachStdout: true,
		AttachStdin:  true,
		AttachStderr: true,
		Detach:       false,
		Tty:          tty,
		Cmd:          []string{"/bin/bash"},
	})
	if err != nil {
		log.Fatalf("error container exec create: %v", err)
	}

	execAttachResponse, err := d.dockerClient.ContainerExecAttach(d.ctx, execCreateResponse.ID, types.ExecStartCheck{
		Tty: tty,
	})
	if err != nil {
		log.Fatalf("error container exec attach: %v", err)
	}
	defer execAttachResponse.Close()

	closeCallback := func() {
		if err := d.dockerClient.ContainerRemove(d.ctx, containerId, types.ContainerRemoveOptions{Force: true}); err != nil {
			log.Fatalf("error docker remove: %v", err)
		}
	}

	var once sync.Once
	go func() {

		if tty {
			_, err := io.Copy(os.Stdout, execAttachResponse.Reader)
			if err != nil {
				log.Fatalf("error copy docker->local: %v", err)
			}
		} else {
			_, err := stdcopy.StdCopy(os.Stdout, os.Stderr, execAttachResponse.Reader)
			if err != nil {
				log.Fatalf("error copy docker->local: %v", err)
			}
		}

		once.Do(closeCallback)
	}()

	go func() {
		_, err = io.Copy(execAttachResponse.Conn, os.Stdin)
		if err != nil {
			log.Fatalf("error copy local->docker: %v", err)
		}

		once.Do(closeCallback)
	}()

	fd := int(os.Stdin.Fd())
	var oldState *terminal.State
	if terminal.IsTerminal(fd) {
		oldState, err = terminal.MakeRaw(fd)
		if err != nil {
			// print error
		}
		defer terminal.Restore(fd, oldState)
	}

	// clear
	//fmt.Print("\033[H\033[2J")
	d.loader.Stop()

	// // TODO <<<

	// // TODO CTRL+C should NOT exit
	// signalCh := make(chan os.Signal)
	// signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	// // signal.Notify(c, os.Interrupt)
	// go func() {
	// 	// for sig := range c {
	// 	// 	// sig is a ^C, handle it
	// 	// 	log.Printf("CTRL+C handler %v", sig)
	// 	// }
	// 	<-signalCh

	// 	log.Printf("CTRL+C handler")
	// 	once.Do(closeCallback)
	// 	//os.Exit(0)
	// }()

	statusCh, errCh := d.dockerClient.ContainerWait(d.ctx, containerId, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Fatalf("error container wait: %v", err)
		}
		log.Printf("close container wait errCh")
	case <-statusCh:
	}
}
