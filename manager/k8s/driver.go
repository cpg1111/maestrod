package k8s

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cpg1111/maestrod/config"
	"github.com/cpg1111/maestrod/manager"
)

// Driver is a struct for the k8s Driver
type Driver struct {
	manager.Driver
	Host           string
	MaestroVersion string
	Client         *http.Client
}

// New returns a pointer to a k8s Driver
func New(maestroVersion string, conf *config.Server) *Driver {
	authTransport, authErr := NewAuthTransport(conf)
	if authErr != nil {
		panic(authErr)
	}
	return &Driver{
		Host:           manager.GetTarget(conf),
		MaestroVersion: maestroVersion,
		Client: &http.Client{
			Transport: authTransport,
		},
	}
}

func (d *Driver) create(url, errObj string, body []byte) error {
	bodyReader := bytes.NewReader(body)
	res, postErr := d.Client.Post(fmt.Sprintf("%s%s", d.Host, url), "application/json", bodyReader)
	if postErr != nil {
		return postErr
	}
	defer res.Body.Close()
	resBody, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return readErr
	}
	if res.StatusCode != 201 {
		return fmt.Errorf("did not create %s, received  %v \n %s", errObj, res.StatusCode, (string)(resBody))
	}
	return nil
}

func (d *Driver) check(url string) (bool, error) {
	res, getErr := d.Client.Get(fmt.Sprintf("%s%s", d.Host, url))
	if res.StatusCode == 404 {
		return false, nil
	} else if res.StatusCode == 200 {
		return true, nil
	}
	return false, getErr
}

// CreateNamespace creates a maestro k8s namespace if one does not exist
func (d *Driver) CreateNamespace(namespace string) error {
	exists, checkErr := d.check(fmt.Sprintf("/api/v1/namespaces/%s", namespace))
	if checkErr != nil {
		return checkErr
	} else if exists {
		return nil
	}
	newNamespace := &Namespace{
		Kind:       "Namespace",
		ApiVersion: "v1",
		Metadata: nsMetadata{
			Name:      namespace,
			Namespace: namespace,
		},
	}
	body, marshErr := json.Marshal(newNamespace)
	if marshErr != nil {
		return marshErr
	}
	return d.create("/api/v1/namespaces", "namespace", body)
}

// CreateSvcAccnt creates a kubernetes svc accnt
func (d *Driver) CreateSvcAccnt(name string) error {
	exists, checkErr := d.check(fmt.Sprintf("/api/v1/namespaces/maestro/serviceaccounts/%s", name))
	if checkErr != nil {
		return checkErr
	} else if exists {
		return nil
	}
	newSvcAccnt := &ServiceAccount{
		Kind:       "ServiceAccount",
		ApiVersion: "v1",
		Metadata: saMetadata{
			Name:      name,
			Namespace: "maestro",
		},
	}
	body, marshErr := json.Marshal(newSvcAccnt)
	if marshErr != nil {
		return marshErr
	}
	return d.create("/api/namepsaces/maestro/serviceaccounts", "service account", body)
}

func (d *Driver) createPod(newPod *Pod) error {
	body, marshErr := json.Marshal(newPod)
	if marshErr != nil {
		return marshErr
	}
	return d.create("/api/v1/namespaces/maestro/pods", "maestro worker", body)
}

// Run will run a maestro pod in kubernetes
func (d Driver) Run(name, confTarget, hostVolume string, args []string) error {
	dPtr := &d
	name = strings.Replace(strings.Replace(name, "/", "-", -1), "_", "-", -1)
	confVol, volErr := NewVolume(fmt.Sprintf("%s-conf", name), hostVolume, dPtr)
	if volErr != nil {
		return volErr
	}
	dockerVol, dockerErr := NewVolume("docker-sock", "/var/run/docker.sock", dPtr)
	if dockerErr != nil {
		return dockerErr
	}
	sec := &secCtx{}
	confContainerVol := newMount(confVol.Name)
	dockerContainerVol := newMount(dockerVol.Name)
	mounts := []volumeMount{*confContainerVol, *dockerContainerVol}
	maestroContainer := NewContainer(d.MaestroVersion, args, mounts, sec)
	newPod := &Pod{
		Kind:       "Pod",
		ApiVersion: "v1",
		Metadata: podMetadata{
			Name:      name,
			Namespace: "maestro",
		},
		Spec: podSpec{
			Volumes:       []Volume{*confVol, *dockerVol},
			Containers:    []Container{*maestroContainer},
			RestartPolicy: "Never",
		},
	}
	return dPtr.createPod(newPod)
}
