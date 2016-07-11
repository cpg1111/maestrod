package k8s

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cpg1111/maestrod/manager"
)

type Driver struct {
	manager.Driver
	Host           string
	MaestroVersion string
}

func New(host, maestroVersion string) *Driver {
	return &Driver{
		Host:           host,
		MaestroVersion: maestroVersion,
	}
}

type podMetadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type podSpec struct {
	Volumes       []Volume    `json:"volumes"`
	Containers    []Container `json:"containers"`
	RestartPolicy string      `json:"restartPolicy"`
}

type Pod struct {
	Kind       string      `json:"kind"`
	ApiVersion string      `json:"apiVersion"`
	Metadata   podMetadata `json:"metadata"`
	Spec       podSpec     `json:"spec"`
}

func (d *Driver) Run(name, confTarget, hostVolume string, args []string) error {
	confVol, volErr := NewVolume(fmt.Sprintf("%s_conf", name), hostVolume, d.Host)
	if volErr != nil {
		return volErr
	}
	confContainerVol := volumeMount{
		Name:      confVol.Name,
		ReadOnly:  false,
		MountPath: confTarget,
	}
	sec := secCtx{}
	maestroContainer := NewContainer(d.MaestroVersion, args, confContainerVol, sec)
	newPod := &Pod{
		Kind:       "pod",
		ApiVersion: "v1",
		Metadata: podMetadata{
			Name:      name,
			Namespace: "maestro",
		},
		Spec: podSpec{
			Volumes:    []Volume{*confVol},
			Containers: []Container{*maestroContainer},
		},
	}
	body, marshErr := json.Marshal(newPod)
	if marshErr != nil {
		return marshErr
	}
	bodyReader := bytes.NewReader(body)
	res, postErr := http.Post("%s/api/v1/namespaces/maestro/pods", "application/json", bodyReader)
	if postErr != nil {
		return postErr
	}
	if res.StatusCode != 201 {
		return errors.New("did not create maestro worker")
	}
	return nil
}
