package k8s

type saMetadata struct {
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

// ServiceAccount is a struct for creating svc accnts
type ServiceAccount struct {
	Kind       string     `json:"kind"`
	ApiVersion string     `json:"apiVersion"`
	Metadata   saMetadata `json:"metadata"`
}
