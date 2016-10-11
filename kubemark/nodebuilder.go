package kubemark

import (
	classType "github.com/wy2745/kubernetes-deployment-tool/type137"
	"net/http"
	"github.com/wy2745/kubernetes-deployment-tool/json"
	"fmt"
	"bytes"
	"io/ioutil"
)

const (
	userName string = "admin"
	password string = "FlbY3CD6mcFUfZvb"

	DestinationServer_Test string = "http://202.120.40.178:1080"
	destinationServer_Test string = "http://202.120.40.177:16380"
	kubemarkServer_Test string = "http://202.120.40.177:17080"
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

func InvokeRequest(method string, url string, body []byte) *http.Response {
	client := http.Client{}
	var req *http.Request
	var err error
	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if method != "GET" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := client.Do(req)
	fmt.Println(resp.Header)
	//fmt.Println(resp.Status)
	//fmt.Println(resp.StatusCode)
	if err != nil {
		fmt.Print(err)
	}
	return resp
}

func BuildNode(num int32) {
	url := destinationServer_Test + GenerateReplicationControllerNamespaceUrl("kubemark")
	b := GenerateReplicationcontrollerBody(num)
	InvokeRequest("POST", url, b)
}
func DeleteNode() {
	url := destinationServer_Test + GenerateReplicationControllerNamespaceUrl("kubemark") + "/" + "hollow-node"
	InvokeRequest("DELETE", url, nil)

	url = destinationServer_Test + GeneratePodNamespaceUrl("kubemark")
	resp := InvokeRequest("GET", url, nil)
	if (resp != nil) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
		}
		var v classType.PodList
		jsonParse.JsonUnmarsha(body, &v)
		for _, pod := range v.Items {
			if pod.Labels["name"] == "hollow-node" {
				url = destinationServer_Test + GeneratePodNameUrl("kubemark", pod.Name)
				InvokeRequest("DELETE", url, nil)
			}
		}
	}
}
func GetNode() {
	url := kubemarkServer_Test + CreateNode_POST
	resp := InvokeRequest("GET", url, nil)
	if (resp != nil) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
		}
		var v classType.NodeList
		jsonParse.JsonUnmarsha(body, &v)
		fmt.Println(v)
	}
}

func GenerateReplicationcontrollerBody(replic int32) []byte {
	var name = "hollow-node"
	var image = "docker.io/wy2745/kubemark:latest"
	var containerName = "hollow-kubelet"
	var labelName = "hollow-node"
	//生成typeMedata
	var typeMedata classType.TypeMeta
	typeMedata.APIVersion = "v1"
	typeMedata.Kind = "ReplicationController"
	//生成objectMedata
	var objectMedata classType.ObjectMeta
	objectMedata.Labels = make(map[string]string)
	objectMedata.Labels["name"] = labelName
	objectMedata.Name = name
	objectMedata.Namespace = "kubemark"

	//生成PodObjectMedata
	var objectMedata2 classType.ObjectMeta
	objectMedata2.Labels = make(map[string]string)
	objectMedata2.Labels["name"] = labelName
	objectMedata2.Namespace = "kubemark"

	//volumns
	var volumn classType.Volume
	volumn.Name = "kubeconfig-volume"
	volumn.Secret.SecretName = "kubeconfig"
	vslice := []classType.Volume{volumn}

	//生成PodTemplateSpec.container
	var Ports []classType.ContainerPort
	var port classType.ContainerPort
	port.ContainerPort = 4194
	Ports = append(Ports, port)
	port.ContainerPort = 10250
	Ports = append(Ports, port)
	port.ContainerPort = 10255
	Ports = append(Ports, port)

	var args []string
	args = append(args, "--v=3")
	args = append(args, "--morph=proxy")
	args = append(args, "$(CONTENT_TYPE)")

	var volumnMount classType.VolumeMount
	volumnMount.Name = "kubeconfig-volume"
	volumnMount.MountPath = "/kubeconfig"
	vmslice := []classType.VolumeMount{volumnMount}

	var resource classType.ResourceRequirements
	resource.Requests = make(map[classType.ResourceName]string)
	resource.Requests["cpu"] = "20m"
	resource.Requests["memory"] = "100M"

	var env classType.EnvVar
	env.Name = "CONTENT_TYPE"
	var vf classType.EnvVarSource
	var vfs classType.ConfigMapKeySelector
	vfs.Name = "node-configmap"
	vfs.Key = "content.type"
	vf.ConfigMapKeyRef = &vfs
	env.ValueFrom = &vf
	eslice := []classType.EnvVar{env}

	var container classType.Container
	container.Name = containerName
	container.Image = image
	container.Ports = Ports
	container.Env = eslice
	container.Args = args
	container.VolumeMounts = vmslice
	container.Resources = resource

	var command []string
	command = append(command, "./kubemark.sh")
	container.Command = command

	slice := []classType.Container{container}
	var podTemplateSpec classType.PodTemplateSpec
	podTemplateSpec.ObjectMeta = objectMedata2
	podTemplateSpec.Spec.Containers = slice
	podTemplateSpec.Spec.Volumes = vslice

	var replicationControllerSpec classType.ReplicationControllerSpec
	replicationControllerSpec.Template = &podTemplateSpec

	replicationControllerSpec.Replicas = &replic
	replicationControllerSpec.Selector = make(map[string]string)
	replicationControllerSpec.Selector["name"] = labelName

	var replicationController classType.ReplicationController
	replicationController.Spec = replicationControllerSpec
	replicationController.ObjectMeta = objectMedata
	replicationController.TypeMeta = typeMedata
	replicationController.Status.Replicas = replic

	b := jsonParse.JsonMarsha(replicationController)
	fmt.Print(string(b))
	return b

}

