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

func New(host, maestroVersion string, conf *config.Server) *Driver {
	authTransport, authErr := NewAuthTransport(conf)
	if authErr != nil {
		panic(authErr)
	}
	return &Driver{
		Host:           host,
		MaestroVersion: maestroVersion,
		Client: &http.Client{
			Transport: authTransport,
		},
	}
}

func (d *Driver) CreateNamespace(namespace string) error {
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
	bodyReader := bytes.NewReader(body)
	res, postErr := d.Client.Post(fmt.Sprintf("%s/api/v1/namespaces", d.Host), "application/json", bodyReader)
	if postErr != nil {
		return postErr
	}
	defer res.Body.Close()
	resBody, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return readErr
	}
	if res.StatusCode != 201 {
		return fmt.Errorf("did not create maestro namespace, received %v \n %s", res.StatusCode, (string)(resBody))
	}
	return nil
}

func (d *Driver) CreateSvcAccnt(name string) error {
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
	bodyReader := bytes.NewReader(body)
	res, postErr := d.Client.Post(fmt.Sprintf("%s/api/v1/namespaces/%s/serviceaccounts", d.Host, "maestro"), "application/json", bodyReader)
	if postErr != nil {
		return postErr
	}
	defer res.Body.Close()
	resBody, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return readErr
	}
	if res.StatusCode != 201 {
		return fmt.Errorf("did not create service account, received %s \n %s", res.StatusCode, (string)(resBody))
	}
	return nil
}

func (d *Driver) createPod(newPod *Pod) error {
	body, marshErr := json.Marshal(newPod)
	if marshErr != nil {
		return marshErr
	}
	bodyReader := bytes.NewReader(body)
	res, postErr := d.Client.Post(fmt.Sprintf("%s/api/v1/namespaces/maestro/pods", d.Host), "application/json", bodyReader)
	if postErr != nil {
		return postErr
	}
	defer res.Body.Close()
	resBody, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return readErr
	}
	if res.StatusCode != 201 {
		return fmt.Errorf("did not create maestro worker, received a status of %v \n %s", res.StatusCode, (string)(resBody))
	}
	return nil
}

func (d *Driver) Run(name, confTarget, hostVolume string, args []string) error {
	confVol, volErr := NewVolume(fmt.Sprintf("%s_conf", name), hostVolume, d)
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
		Kind:       "Pod",
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
	return d.createPod(newPod)
}
