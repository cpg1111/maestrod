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
		} `json:"fieldRef,omitempty"`
		ConfigMapKeyRef struct {
			Name string `json:"name"`
			Key  string `json:"key"`
		} `json:"configMapKeyRef,omitempty"`
		SecretKeyRef struct {
			Name string `json:"name"`
			Key  string `json:"key"`
		} `json:"secretKeyRef,omitempty"`
	} `json:"valueFrom,omitempty"`
}

type volumeMount struct {
	Name      string `json:"name"`
	ReadOnly  bool   `json:"readOnly,omitempty"`
	MountPath string `json:"mountPath"`
}

type secCtx struct {
	Capabilities struct {
		Add  []string `json:"add"`
		Drop []string `json:"drop"`
	} `json:"capabilities,omitempty"`
	SELinuxOptions struct {
		User  string `json:"user"`
		Role  string `json:"role"`
		Type  string `json:"type"`
		Level string `json:"level"`
	} `json:"seLinuxOptions,omitempty"`
	RunAsUser              int  `json:"runAsUser"`
	RunAsNonRoot           bool `json:"runAsNonRoot"`
	ReadOnlyRootFileSystem bool `json:"readOnlyRootFileSystem"`
}

type Container struct {
	Name            string        `json:"name"`
	Image           string        `json:"image"`
	Command         []string      `json:"command,omitempty"`
	Args            []string      `json:"args,omitempty"`
	WorkingDir      string        `json:"workingDir,omitempty"`
	Ports           []port        `json:"ports,omitempty"`
	Env             []envVar      `json:"env,omitempty"`
	VolumeMounts    []volumeMount `json:"volumeMounts,omitempty"`
	SecurityContext *secCtx       `json:"securityContext,omitempty"`
}

func NewContainer(maestroVersion string, cmd []string, vol volumeMount, sec *secCtx) *Container {
	return &Container{
		Name:            fmt.Sprintf("maestro-%s", maestroVersion),
		Image:           fmt.Sprintf("cpg1111/maestro:%s", maestroVersion),
		Command:         cmd,
		VolumeMounts:    []volumeMount{vol},
		SecurityContext: sec,
	}
}
