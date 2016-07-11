package k8s

import (
	"fmt"
)

type port struct {
	Name          string `json:"name"`
	HostPort      int    `json:"hostPort"`
	ContainerPort int    `json:"containerPort"`
	Protocol      string `json:"protocol"`
	HostIP        string `json:"hostIP,omitempty"`
}

type envVar struct {
	Name      string `json:"name"`
	Value     string `json:"value,omitempty"`
	ValueFrom struct {
		FieldRef struct {
			APIVersion string `json:"apiVersion"`
			FieldPath  string `json:"fieldPath"`
		} `json:"fieldRef"`
		ConfigMapKeyRef struct {
			Name string `json:"name"`
			Key  string `json:"key"`
		} `json:"configMapKeyRef"`
		SecretKeyRef struct {
			Name string `json:"name"`
			Key  string `json:"key"`
		} `json:"secretKeyRef"`
	} `json:"valueFrom,omitempty"`
}

type volumeMount struct {
	Name      string `json:"name"`
	ReadOnly  bool   `json:"readOnly"`
	MountPath string `json:"mountPath"`
}

type secCtx struct {
	Capabilities struct {
		Add  []string `json:"add"`
		Drop []string `json:"drop"`
	} `json:"capabilities"`
	SELinuxOptions struct {
		User  string `json:"user"`
		Role  string `json:"role"`
		Type  string `json:"type"`
		Level string `json:"level"`
	} `json:"seLinuxOptions"`
	RunAsUser              int  `json:"runAsUser"`
	RunAsNonRoot           bool `json:"runAsNonRoot"`
	ReadOnlyRootFileSystem bool `json:"readOnlyRootFileSystem"`
}

type Container struct {
	Name            string        `json:"name"`
	Image           string        `json:"image"`
	Command         []string      `json:"command"`
	Args            []string      `json:"args"`
	WorkingDir      string        `json:"workingDir"`
	Ports           []port        `json:"ports"`
	Env             []envVar      `json:"env"`
	VolumeMounts    []volumeMount `json:"volumeMounts"`
	SecurityContext secCtx        `json:"securityContext"`
}

func NewContainer(maestroVersion string, cmd []string, vol volumeMount, sec secCtx) *Container {
	return &Container{
		Name:            fmt.Sprintf("maestro:%s", maestroVersion),
		Image:           fmt.Sprintf("cpg1111/maestro:%s", maestroVersion),
		Command:         cmd,
		VolumeMounts:    []volumeMount{vol},
		SecurityContext: sec,
	}
}
