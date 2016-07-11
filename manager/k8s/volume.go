package k8s

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/cpg1111/maestrod/config"
)

type downwardAPIObj struct {
	Path     string `json:"path"`
	FieldRef struct {
		Name string `json:"name"`
	} `json:"fieldRef"`
}

type keyPath struct {
	Key  string `json:"key"`
	Path string `json:"path"`
}

type hostPath struct {
	Path string `json:"path"`
}

type emptyDir struct {
	Medium string `json:"medium"`
}

type gcePersistentDisk struct {
	PDName    string `json:"pdName"`
	FSType    string `json:"fsType"`
	Partition int    `json:"partition"`
	ReadOnly  bool   `json:"readOnly"`
}

type awsElasticBlockStore struct {
	VolumeID  string `json:"volumeID"`
	FSType    string `json:"fsType"`
	Partition int    `json:"partition"`
	ReadOnly  bool   `json:"readOnly"`
}

type nfs struct {
	Server   string `json:"server"`
	Path     string `json:"path"`
	ReadOnly bool   `json:"readOnly"`
}

type gluster struct {
	Endpoints string `json:"endpoints"`
	Path      string `json:"path"`
	ReadOnly  bool   `json:"readOnly"`
}

type claim struct {
	ClaimName string `json:"claimName"`
	ReadOnly  bool   `json:"readOnly"`
}

type cinder struct {
	VolumeID string `json:"volumeID"`
	FSType   string `json:"fsType"`
	ReadOnly bool   `json:"readOnly"`
}

type flocker struct {
	DatasetName string `json:"datasetName"`
}

type Volume struct {
	Name                  string               `json:"-"`
	Type                  string               `json:"-"`
	HostPath              hostPath             `json:"hostPath,omitempty"`
	EmptyDir              emptyDir             `json:"emptDir,omitempty"`
	GCEPersistentDisk     gcePersistentDisk    `json:"gcePersistentDisk,omitempty"`
	AWSElasticBlockStore  awsElasticBlockStore `json:"awsElasticBlockStore,omitempty"`
	NFS                   nfs                  `json:"nfs,omitempty"`
	GlusterFS             gluster              `json:"glusterfs,omitempty"`
	PersistentVolumeClaim claim                `json:"persistentVolumeClaim,omitempty"`
	Cinder                cinder               `json:"cinder,omitempty"`
	Flocker               flocker              `json:"flocker,omitempty"`
}

func (v *Volume) DelegateType(confTarget, confTargetPrefix string, volumeConf *config.Mount) {
	if volumeConf != nil && volumeConf.Kind != confTargetPrefix {
		log.Fatal("mount configuration mismatch")
	}
	v.Type = confTargetPrefix
	switch confTargetPrefix {
	case "hostPath":
		v.HostPath = hostPath{Path: confTarget}
		break
	case "emptyDir":
		v.EmptyDir = emptyDir{}
		break
	case "gcePersistentDisk":
		v.GCEPersistentDisk = gcePersistentDisk{
			PDName:   volumeConf.ID,
			FSType:   volumeConf.FSType,
			ReadOnly: volumeConf.ReadOnly,
		}
		break
	case "awsElasticBlockStore":
		v.AWSElasticBlockStore = awsElasticBlockStore{
			VolumeID: volumeConf.ID,
			FSType:   volumeConf.FSType,
			ReadOnly: volumeConf.ReadOnly,
		}
		break
	case "nfs":
		v.NFS = nfs{
			Server:   volumeConf.Server,
			Path:     volumeConf.Path,
			ReadOnly: volumeConf.ReadOnly,
		}
		break
	case "glusterfs":
		v.GlusterFS = gluster{
			Endpoints: volumeConf.Endpoints,
			Path:      volumeConf.Path,
			ReadOnly:  volumeConf.ReadOnly,
		}
		break
	case "persistentVolumeClaim":
		v.PersistentVolumeClaim = claim{
			ClaimName: volumeConf.Path,
			ReadOnly:  volumeConf.ReadOnly,
		}
		break
	case "flocker":
		v.Flocker = flocker{
			DatasetName: volumeConf.Name,
		}
		break
	default:
		log.Fatal("volume mount type not supported, please create an issue @ github.com/cpg1111/maestrod")
	}
}

func checkForVolume(host string, vol *Volume) (bool, error) {
	res, getErr := http.Get(fmt.Sprintf("%s/api/v1/persistentVolume/%s", host, vol.Name))
	if getErr != nil {
		return false, getErr
	}
	if res.StatusCode == 403 {
		return false, errors.New("Bad permissions for volume mount")
	} else if res.StatusCode >= 404 {
		return false, nil
	}
	return true, nil
}

type payloadMetadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type createPayload struct {
	Kind       string          `json:"kind"`
	APIVersion string          `json:"apiVersion"`
	Metadata   payloadMetadata `json:"metadata"`
	Spec       interface{}     `json:"spec"`
}

func createVolume(host string, vol *Volume) error {
	client := http.Client{}
	payload := &createPayload{
		Kind:       "persistentVolume",
		APIVersion: "v1",
		Metadata: payloadMetadata{
			Name:      vol.Name,
			Namespace: "maestro",
		},
		Spec: vol,
	}
	body, marshErr := json.Marshal(payload)
	if marshErr != nil {
		return marshErr
	}
	bodyReader := bytes.NewReader(body)
	res, postErr := client.Post(fmt.Sprintf("%s/api/v1/persistentVolumes/", host), "application/json", bodyReader)
	if postErr != nil {
		return postErr
	}
	log.Println(res)
	return nil
}

func getVolume(name, confTarget string) Volume {
	vol := Volume{
		Name: fmt.Sprintf("%s_config", name),
	}
	if !strings.Contains(confTarget, "://") {
		vol.Type = "hostPath"
		vol.HostPath = hostPath{Path: confTarget}
		return vol
	}
	sepIndex := strings.Index(confTarget, "://")
	prefix := confTarget[0:sepIndex]
	vol.DelegateType(confTarget, prefix, nil)
	return vol
}

func NewVolume(name, confTarget, host string) (*Volume, error) {
	vol := getVolume(name, confTarget)
	volPtr := &vol
	if vol.Type != "hostPath" {
		found, checkErr := checkForVolume(host, volPtr)
		if checkErr != nil {
			return nil, checkErr
		}
		if !found {
			createErr := createVolume(host, volPtr)
			if createErr != nil {
				return nil, createErr
			}
		}
	}
	return volPtr, nil
}
