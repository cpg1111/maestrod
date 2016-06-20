package docker

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cpg1111/maestrod/manager"

	dockerEngine "github.com/docker/engine-api/client"
	dockerTypes "github.com/docker/engine-api/types"
	dockerContainer "github.com/docker/engine-api/types/container"
	dockerNetwork "github.com/docker/engine-api/types/network"
	"golang.org/x/net/context"
)

type Driver struct {
	manager.Driver
	client      *dockerEngine.Client
	containerID string
	image       string
	cmd         []string
}

func New(host, apiVersion, maestroVersion, name string) (*Driver, error) {
	hClient := &http.Client{}
	dClient, dockerErr := dockerEngine.NewClient(host, apiVersion, hClient, make(map[string]string))
	if dockerErr != nil {
		return nil, dockerErr
	}
	return &Driver{
		client:      dClient,
		containerID: name,
		image:       fmt.Sprintf("cpg1111/maestro:%s", maestroVersion),
	}, nil
}

func (d *Driver) needToPull(ctx context.Context) (bool, error) {
	listOptions := &dockerTypes.ImageListOptions{
		MatchName: "cpg1111/maestro",
		All:       true,
	}
	images, listErr := d.client.ImageList(ctx, *listOptions)
	if listErr != nil {
		return false, listErr
	}
	return len(images) > 0, nil
}

func (d *Driver) pull(ctx context.Context) error {
	pullOptions := dockerTypes.ImagePullOptions{}
	res, resErr := d.client.ImagePull(ctx, d.image, pullOptions)
	defer res.Close()
	if resErr != nil {
		return resErr
	}
	resp := make([]byte, 4096)
	_, readErr := res.Read(resp)
	if readErr != nil {
		return readErr
	}
	log.Println((string)(resp))
	return nil
}

func (d *Driver) getContainerConfig() *dockerContainer.Config {
	labels := make(map[string]string)
	labels["NAME"] = d.containerID
	timeout := 5
	return &dockerContainer.Config{
		User:         "maestro",
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		OpenStdin:    false,
		StdinOnce:    false,
		Cmd:          d.cmd,
		Healthcheck: &dockerContainer.HealthConfig{
			Test: []string{""},
		},
		ArgsEscaped:     true,
		Image:           d.image,
		Volumes:         make(map[string]struct{}),
		WorkingDir:      "/",
		NetworkDisabled: false,
		Labels:          labels,
		StopSignal:      "SIGTERM",
		StopTimeout:     &timeout,
	}
}

func (d *Driver) getHostConfig() *dockerContainer.HostConfig {
	logConf := make(map[string]string)
	return &dockerContainer.HostConfig{
		Binds:           []string{},
		ContainerIDFile: "/tmp/containers",
		LogConfig: dockerContainer.LogConfig{
			Type:   "json-file",
			Config: logConf,
		},
		NetworkMode: "host",
		RestartPolicy: dockerContainer.RestartPolicy{
			Name:              "never",
			MaximumRetryCount: 0,
		},
		AutoRemove:      true,
		VolumeDriver:    "local",
		Privileged:      false,
		PublishAllPorts: false,
		ReadonlyRootfs:  false,
	}
}

func (d *Driver) getNetworkConfig() *dockerNetwork.NetworkingConfig {
	return &dockerNetwork.NetworkingConfig{}
}

func (d *Driver) create(ctx context.Context) error {
	createResp, err := d.client.ContainerCreate(ctx, d.getContainerConfig(), d.getHostConfig(), d.getNetworkConfig(), d.containerID)
	if err != nil {
		return err
	}
	log.Println(createResp)
	return nil
}

func (d Driver) Run(args []string) error {
	d.cmd = args
	ctx := context.Background()
	needToPull, checkErr := d.needToPull(ctx)
	if checkErr != nil {
		return checkErr
	}
	if needToPull {
		pullErr := d.pull(ctx)
		if pullErr != nil {
			return pullErr
		}
	}
	createErr := d.create(ctx)
	if createErr != nil {
		return createErr
	}
	return nil
}
