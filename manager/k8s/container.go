/*
Copyright 2016 Christian Grabowski All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

type volumeMount struct {
	Name      string `json:"name"`
	ReadOnly  bool   `json:"readOnly,omitempty"`
	MountPath string `json:"mountPath"`
}

func newMount(name, path string) *volumeMount {
	return &volumeMount{
		Name:      name,
		ReadOnly:  false,
		MountPath: path,
	}
}

func newMounts(vols []Volume) []volumeMount {
	mounts := make([]volumeMount, len(vols))
	for v := range vols {
		mount := newMount(vols[v].Name, vols[v].HostPath.Path)
		mounts[v] = *mount
	}
	return mounts
}

// Container is a struct reoresenting a k8s container API object
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
	TTY             bool          `json:"tty,omitempty"`
}

// NewContainer returns a pointer to a Container struct
func NewContainer(maestroVersion string, cmd []string, vols []volumeMount, sec *secCtx) *Container {
	return &Container{
		Name:            "maestro",
		Image:           fmt.Sprintf("cpg1111/maestro:%s", maestroVersion),
		Command:         cmd,
		VolumeMounts:    vols,
		SecurityContext: sec,
		TTY:             false,
	}
}
