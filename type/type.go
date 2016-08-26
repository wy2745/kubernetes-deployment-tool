package classType

type Metadata struct{
	Name string `json:"name,omitempty"`
	GenerateName string `json:"generateName,omitempty"`
	SelfLink string `json:"selfLink,omitempty"`
	ResourceVersion string `json:"resourceVersion,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Uid string `json:"uid,omitempty"`
	Generation int `json:"generation,omitempty"`
	CreationTimestamp string `json:"CreationTimestamp,omitempty"`
	DeletionTimestamp string `json:"deletionTimestamp,omitempty"`
	DeletionGracePeriodSeconds int `json:"deletionGracePeriodSeconds,omitempty"`
}


type Item struct{
	Kind string `json:"kind,omitempty"`
	ApiVersion string `json:"apiVersion,omitempty"`
	Metadata Metadata `json:"metadata,omitempty"`
	Spec Spec `json:"spec,omitempty"`
	Status Status `json:"status,omitempty"`
}

type Spec struct{
	PodCIDR string `json:"podCIDR,omitempty"`
	ExternalID string `json:"externalID,omitempty"`
	ProviderID string `json:"providerID,omitempty"`
	Unschedulable bool `json:"unschedulable,omitempty"`
}

type Condition struct{
	Type string `json:"type,omitempty"`
	Status string `json:"status,omitempty"`
	LastHeartbeatTime string `json:"LastHeartbeatTime,omitempty"`
	LastTransitionTime string `json:"lastTransitionTime,omitempty"`
	Reason string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}
type Status struct{
	Phase string `json:"phase,omitempty"`
	Conditions []Condition `json:"conditions,omitempty"`
	Addresses []Address `json:"addresses,omitempty"`
}
type Address struct{
	Type string `json:"type,omitempty"`
	Address string `json:"address,omitempty"`
}

type KubectlEndpoint struct{
	Port int `json:"Port,omitempty"`
}
type DaemonEndpoints struct{
	KubeletEndpoints KubectlEndpoint `json:"kubeletEndpoint,omitempty"`
}
type NodeInfo struct{
	MachineID string `json:"machineID,omitempty"`
	SystemUUID string `json:"systemUUID,omitempty"`
	BootID string `json:"bootID,omitempty"`
	KernelVersion string `json:"kernelVersion,omitempty"`
	OsImage string `json:"osImage,omitempty"`
	ContainerRuntimeVersion string `json:"containerRuntimeVersion,omitempty"`
	KubeletVersion string `json:"kubeletVersion,omitempty"`
	KubeProxyVersion string `json:"kubeProxyVersion,omitempty"`
}



type Node struct{

	Kind  string `json:"kind,omitempty"`
	ApiVersion  string `json:"apiVersion,omitempty"`
	Metadata Metadata `json:"metadata,omitempty"`
	Items []Item `json:"items,omitempty"`
	DaemonEndpoints DaemonEndpoints `json:"daemonEndpoints,omitempty"`
	NodeInfo NodeInfo `json:"nodeInfo,omitempty"`
}

