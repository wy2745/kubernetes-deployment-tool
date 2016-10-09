package request

const (
	userName string = "admin"
	password string = "FlbY3CD6mcFUfZvb"

	destinationServer_Test string = "http://202.120.40.177:16380"
	KubemarkServer_Test string = "http://202.120.40.177:17080"
	destinationServer_Caicloud string = "https://sjtu.caicloudapp.com"
	GetNodeList_GET string = "/api/v1/nodes"   //list or watch objects of kind Node
	CreateNode_POST string = "/api/v1/nodes"    //create a Node
	DeleteNode_DELETE string = "/api/v1/nodes/{name}"  //delete a Node
	ReadNode_GET string = "/api/v1/nodes/{name}"     //read the specified Node
	GetNamespaces_GET string = "/api/v1/namespaces"    //list or watch objects of kind Namespace
	CreateNamespace_POST string = "/api/v1/namespaces" //create a Namespace
	DeleteNamespace_DELETE string = "/api/v1/namespaces/{name}" //delete a Namespace
	ReadNamespace_GET string = "/api/v1/namespaces/{name}" //read the specified Namespace

	ReadPodsListOfNamespace_GET string = "/api/v1/namespaces/{namespace}/pods" //list or watch objects of kind Pod
	CreatePods_POST string = "/api/v1/namespaces/{namespace}/pods"  //create a Pod
	DeletePod_DELETE string = "/api/v1/namespaces/{namespace}/pods/{name}" //delete a Pod
	ReadPod_GET string = "/api/v1/namespaces/{namespace}/pods/{name}" //read the specified Pod

	ReadAllService_GET string = "/api/v1/services" //list or watch objects of kind Service
	ReadServicesListOfNamespace_GET string = "/api/v1/namespaces/{namespace}/services"  //list or watch objects of kind Service
	CreateService_POST string = "/api/v1/namespaces/{namespace}/services" //create a Service
	DeleteService_DELETE string = "/api/v1/namespaces/{namespace}/services/{name}" //delete a Service
	ReadService_GET string = "/api/v1/namespaces/{namespace}/services/{name}" //read the specified Service

	ReadAllReplicationController_GET string = "/api/v1/replicationcontrollers" //list or watch objects of kind ReplicationController
	ReadReplicationControllerListOfNamespace_GET string = "/api/v1/namespaces/{namespace}/replicationcontrollers" //list or watch
	// objects of kind ReplicationController

	extensionApiNamespace string = "/apis/batch/v1/namespaces"
	ReadAllJob_GET string = "/apis/batch/v1/jobs"
)

func GeneratePodNamespaceUrl(namespace string) string {
	return GetNamespaces_GET + "/" + namespace + "/pods"
}
func GenerateJobNamespaceUrl(namespace string) string {
	return extensionApiNamespace + "/" + namespace + "/jobs"
}
func GeneratePodNameUrl(namespace string, name string) string {
	return GetNamespaces_GET + "/" + namespace + "/pods" + "/" + name
}
func GenerateJobNameUrl(namespace string, name string) string {
	return extensionApiNamespace + "/" + namespace + "/jobs" + "/" + name
}
func GenerateServiceListNamespaceUrl(namespace string) string {
	return GetNamespaces_GET + "/" + namespace + "/services"
}
func GenerateServiceListNameUrl(namespace string, name string) string {
	return GetNamespaces_GET + "/" + namespace + "/services" + "/" + name
}
func GenerateReplicationControllerNamespaceUrl(namespace string) string {
	return GetNamespaces_GET + "/" + namespace + "/replicationcontrollers"
}
func GenerateReplicationControllerNameUrl(namespace string, name string) string {
	return GetNamespaces_GET + "/" + namespace + "/replicationcontrollers" + "/" + name
}
func GenerateNodeNameUrl(name string) string {
	return CreateNode_POST + "/" + name
}
func GenerateNodeUrl() string {
	return CreateNode_POST
}

