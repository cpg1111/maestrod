package k8s

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
