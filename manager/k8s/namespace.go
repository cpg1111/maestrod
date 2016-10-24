package k8s

type nsMetadata struct {
	Name                       string `json:"name"`
	GenerateName               string `json:"generateName,omitempty"`
	Namespace                  string `json:"namespace,omitempty"`
	SelfLink                   string `json:"selfLink,omitempty"`
	UID                        string `json:"uid,omitempty"`
	ResourceVersion            string `json:"resourceVersion,omitempty"`
	Generation                 int    `json:"generation,omitempty"`
	CreationTimestamp          string `json:"creationTimestamp,omitempty"`
	DeletionTimestamp          string `json:"deletionTimestamp,omitempty"`
	DeletionGracePeriodSeconds int    `json:"deletionGracePeriodSeconds,omitempty"`
}

// Namespace is a struct for creating k8s Namespaces
type Namespace struct {
	Kind       string     `json:"kind"`
	ApiVersion string     `json:"apiVersion"`
	Metadata   nsMetadata `json:"metadata"`
}
