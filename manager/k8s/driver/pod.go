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

type podMetadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type podSpec struct {
	Volumes       []Volume    `json:"volumes"`
	Containers    []Container `json:"containers"`
	RestartPolicy string      `json:"restartPolicy,omitempty"`
}

// Pod is a struct for creating k8s pods
type Pod struct {
	Kind       string      `json:"kind"`
	ApiVersion string      `json:"apiVersion"`
	Metadata   podMetadata `json:"metadata"`
	Spec       podSpec     `json:"spec"`
}
