package request

const (
	destinationServer string = "http://202.120.40.177:16380"
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
)