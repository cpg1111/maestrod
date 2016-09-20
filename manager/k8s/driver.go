package k8s

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cpg1111/maestrod/config"
	"github.com/cpg1111/maestrod/manager"
)

type Driver struct {
	manager.Driver
	Host           string
	MaestroVersion string
	Client         *http.Client
}

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

func (d Driver) Run(name, confTarget, hostVolume string, args []string) error {
	dPtr := &d
	confVol, volErr := NewVolume(fmt.Sprintf("%s-conf", name), hostVolume, dPtr)
	if volErr != nil {
		return volErr
	}
	sec := &secCtx{}
	confContainerVol := newMount(confVol.Name)
	maestroContainer := NewContainer(d.MaestroVersion, args, *confContainerVol, sec)
	newPod := &Pod{
		Kind:       "Pod",
		ApiVersion: "v1",
		Metadata: podMetadata{
			Name:      name,
			Namespace: "maestro",
		},
		Spec: podSpec{
			Volumes:       []Volume{*confVol},
			Containers:    []Container{*maestroContainer},
			RestartPolicy: "Never",
		},
	}
	return dPtr.createPod(newPod)
}
