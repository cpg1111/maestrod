package docker

import (
	"fmt"
	"log"

	"github.com/cpg1111/maestrod/manager"

	dockerEngine "github.com/docker/engine-api/client"
	dockerTypes "github.com/docker/engine-api/types"
	dockerContainer "github.com/docker/engine-api/types/container"
	dockerNetwork "github.com/docker/engine-api/types/network"
	"golang.org/x/net/context"
)

// Driver is the driver to the docker daemon
type Driver struct {
	manager.Driver
	client      *dockerEngine.Client
	containerID string
	image       string
	cmd         []string
	confTarget  string
	hostVolume  string
}

// New returns a pointer to an instance of Driver and an error
func New(host, apiVersion, maestroVersion string) (*Driver, error) {
	dClient, dockerErr := dockerEngine.NewEnvClient()
	if dockerErr != nil {
		return nil, dockerErr
	}
	return &Driver{
		client: dClient,
		image:  fmt.Sprintf("cpg1111/maestro:%s", maestroVersion),
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

type empty struct{}

func (d *Driver) getContainerConfig() *dockerContainer.Config {
	labels := make(map[string]string)
	labels["NAME"] = d.containerID
	timeout := 5
	volumes := make(map[string]struct{})
	volumes[d.confTarget] = empty{}
	return &dockerContainer.Config{
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		OpenStdin:    false,
		StdinOnce:    false,
		Cmd:          d.cmd[1:],
		Healthcheck: &dockerContainer.HealthConfig{
			Test: []string{""},
		},
		ArgsEscaped:     true,
		Image:           d.image,
		Volumes:         volumes,
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
		Binds:           []string{fmt.Sprintf("%s:%s:ro", d.confTarget, d.hostVolume)},
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

func (d *Driver) needToRemove(ctx context.Context) (bool, string, error) {
	options := dockerTypes.ContainerListOptions{
		All: true,
	}
	list, listErr := d.client.ContainerList(ctx, options)
	if listErr != nil {
		return false, "", listErr
	}
	for i := range list {
		if list[i].Labels["NAME"] == d.containerID {
			return true, list[i].ID, nil
		}
	}
	return false, "", nil
}

func (d *Driver) remove(ctx context.Context, id string) error {
	options := dockerTypes.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         true,
	}
	return d.client.ContainerRemove(ctx, id, options)
}

func (d *Driver) create(ctx context.Context) error {
	createResp, err := d.client.ContainerCreate(ctx, d.getContainerConfig(), d.getHostConfig(), d.getNetworkConfig(), d.containerID)
	if err != nil {
		return err
	}
	log.Println(createResp)
	return nil
}

func (d *Driver) start(ctx context.Context) error {
	return d.client.ContainerStart(ctx, d.containerID, dockerTypes.ContainerStartOptions{})
}

// Run runs a maestro container on the docker daemon
func (d Driver) Run(name, confTarget, hostVolume string, args []string) error {
	d.containerID = fmt.Sprintf("maestro-%s", name)
	d.confTarget = confTarget
	d.hostVolume = hostVolume
	d.cmd = args
	needToPull, checkErr := d.needToPull(context.Background())
	if checkErr != nil {
		return checkErr
	}
	if needToPull {
		pullErr := d.pull(context.Background())
		if pullErr != nil {
			return pullErr
		}
	}
	needToRemoveOld, removalID, checkRemoveErr := d.needToRemove(context.Background())
	if checkRemoveErr != nil {
		return checkRemoveErr
	}
	if needToRemoveOld {
		removeErr := d.remove(context.Background(), removalID)
		if removeErr != nil {
			return removeErr
		}
	}
	createErr := d.create(context.Background())
	if createErr != nil {
		return createErr
	}
	return d.start(context.Background())
}
