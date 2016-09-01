package k8s

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cpg1111/maestrod/config"
)

type downwardAPIObj struct {
	Path     string `json:"path,omitempty"`
	FieldRef struct {
		Name string `json:"name,omitempty"`
	} `json:"fieldRef,omitempty"`
}

type keyPath struct {
	Key  string `json:"key,omitempty"`
	Path string `json:"path,omitempty"`
}

type hostPath struct {
	Path string `json:"path,omitempty"`
}

type emptyDir struct {
	Medium string `json:"medium,omitempty"`
}

type gcePersistentDisk struct {
	PDName    string `json:"pdName,omitempty"`
	FSType    string `json:"fsType,omitempty"`
	Partition int    `json:"partition,omitempty"`
	ReadOnly  bool   `json:"readOnly,omitempty"`
}

type awsElasticBlockStore struct {
	VolumeID  string `json:"volumeID,omitempty"`
	FSType    string `json:"fsType,omitempty"`
	Partition int    `json:"partition,omitempty"`
	ReadOnly  bool   `json:"readOnly,omitempty"`
}

type nfs struct {
	Server   string `json:"server,omitempty"`
	Path     string `json:"path,omitempty"`
	ReadOnly bool   `json:"readOnly,omitempty"`
}

type gluster struct {
	Endpoints string `json:"endpoints,omitempty"`
	Path      string `json:"path,omitempty"`
	ReadOnly  bool   `json:"readOnly,omitempty"`
}

type claim struct {
	ClaimName string `json:"claimName,omitempty"`
	ReadOnly  bool   `json:"readOnly,omitempty"`
}

type cinder struct {
	VolumeID string `json:"volumeID,omitempty"`
	FSType   string `json:"fsType,omitempty"`
	ReadOnly bool   `json:"readOnly,omitempty"`
}

type flocker struct {
	DatasetName string `json:"datasetName,omitempty"`
}

type Volume struct {
	Name                  string                `json:"-"`
	Type                  string                `json:"-"`
	HostPath              *hostPath             `json:"hostPath,omitempty"`
	GCEPersistentDisk     *gcePersistentDisk    `json:"gcePersistentDisk,omitempty"`
	AWSElasticBlockStore  *awsElasticBlockStore `json:"awsElasticBlockStore,omitempty"`
	NFS                   *nfs                  `json:"nfs,omitempty"`
	GlusterFS             *gluster              `json:"glusterfs,omitempty"`
	PersistentVolumeClaim *claim                `json:"persistentVolumeClaim,omitempty"`
	Cinder                *cinder               `json:"cinder,omitempty"`
	Flocker               *flocker              `json:"flocker,omitempty"`
}

func (v *Volume) DelegateType(confTarget, confTargetPrefix string, volumeConf *config.Mount) {
	if volumeConf != nil && volumeConf.Kind != confTargetPrefix {
		log.Fatal("mount configuration mismatch")
	}
	v.Type = confTargetPrefix
	switch confTargetPrefix {
	case "hostPath":
		v.HostPath = &hostPath{Path: confTarget}
		break
	case "gcePersistentDisk":
		v.GCEPersistentDisk = &gcePersistentDisk{
			PDName:   volumeConf.ID,
			FSType:   volumeConf.FSType,
			ReadOnly: volumeConf.ReadOnly,
		}
		break
	case "awsElasticBlockStore":
		v.AWSElasticBlockStore = &awsElasticBlockStore{
			VolumeID: volumeConf.ID,
			FSType:   volumeConf.FSType,
			ReadOnly: volumeConf.ReadOnly,
		}
		break
	case "nfs":
		v.NFS = &nfs{
			Server:   volumeConf.Server,
			Path:     volumeConf.Path,
			ReadOnly: volumeConf.ReadOnly,
		}
		break
	case "glusterfs":
		v.GlusterFS = &gluster{
			Endpoints: volumeConf.Endpoints,
			Path:      volumeConf.Path,
			ReadOnly:  volumeConf.ReadOnly,
		}
		break
	case "persistentVolumeClaim":
		v.PersistentVolumeClaim = &claim{
			ClaimName: volumeConf.Path,
			ReadOnly:  volumeConf.ReadOnly,
		}
		break
	case "flocker":
		v.Flocker = &flocker{
			DatasetName: volumeConf.Name,
		}
		break
	default:
		log.Fatal("volume mount type not supported, please create an issue @ github.com/cpg1111/maestrod")
	}
}

func checkForVolume(driver *Driver, vol *Volume) (bool, error) {
	res, getErr := driver.Client.Get(fmt.Sprintf("%s/api/v1/persistentVolume/%s", driver.Host, vol.Name))
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

func createVolume(driver *Driver, vol *Volume) error {
	client := driver.Client
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
	res, postErr := client.Post(fmt.Sprintf("%s/api/v1/persistentVolumes/", driver.Host), "application/json", bodyReader)
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
		vol.HostPath = &hostPath{Path: confTarget}
		return vol
	}
	sepIndex := strings.Index(confTarget, "://")
	prefix := confTarget[0:sepIndex]
	if len(vol.Type) == 0 {
		vol.DelegateType(confTarget, prefix, nil)
	}
	return vol
}

func NewVolume(name, confTarget string, driver *Driver) (*Volume, error) {
	vol := getVolume(name, confTarget)
	volPtr := &vol
	if vol.Type != "hostPath" {
		found, checkErr := checkForVolume(driver, volPtr)
		if checkErr != nil {
			return nil, checkErr
		}
		if !found {
			createErr := createVolume(driver, volPtr)
			if createErr != nil {
				return nil, createErr
			}
		}
	}
	return volPtr, nil
}
