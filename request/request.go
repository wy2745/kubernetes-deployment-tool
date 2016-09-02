package request

import (
	"net/http"
	jsonParse "../json"
	classType "../type"
	"io/ioutil"
	"io"
	"fmt"
	"bytes"
)

func setBasicAuthOfCaicloud(r *http.Request) {
	r.SetBasicAuth(userName, password)
}

func InvokeRequest_Caicloud(method string, url string, body io.Reader) *http.Response {
	client := http.Client{}
	req, err := http.NewRequest(method, url, body)
	setBasicAuthOfCaicloud(req)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err)
	}
	return resp
}

func InvokeGetReuqest(url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Print(err)
		return nil
	}
	return resp
}

func GetAllNode_Caicloud() {
	resp := InvokeRequest_Caicloud("GET", destinationServer_Caicloud + GetNodeList_GET, nil)
	if (resp != nil) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
			return
		}
		var v classType.NodeList
		jsonParse.JsonUnmarsha(body, &v)
		for _, item := range v.Items {
			classType.PrintNode(item)
		}
	}
}

func GetAllNode_Test() {
	resp := InvokeGetReuqest(destinationServer_Test + GetNodeList_GET)
	if (resp != nil) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
			return
		}
		var v classType.NodeList
		jsonParse.JsonUnmarsha(body, &v)
		for _, item := range v.Items {
			classType.PrintNode(item)
		}
	}
}

func GetAllNamespace_Test() {
	resp := InvokeGetReuqest(destinationServer_Test + GetNamespaces_GET)
	if (resp != nil) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
			return
		}
		var v classType.NamespaceList
		jsonParse.JsonUnmarsha(body, &v)
		for _, item := range v.Items {
			classType.PrintNamespace(item)
		}
	}
}

func GetAllReplicationcontrollers_Test() {
	resp := InvokeGetReuqest(destinationServer_Test + ReadAllReplicationController_GET)
	if (resp != nil) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
			return
		}
		var v classType.ReplicationControllerList
		jsonParse.JsonUnmarsha(body, &v)
		for _, item := range v.Items {
			classType.PrintReplicationController(item)
		}
	}
}

func GetAllService_Test(namespace string) {
	resp := InvokeGetReuqest(destinationServer_Test + ReadAllService_GET)
	if (resp != nil) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
			return
		}
		var v classType.ServiceList
		jsonParse.JsonUnmarsha(body, &v)
		for _, item := range v.Items {
			classType.PrintService(item)
		}
	}
}

func GetPodsOfNamespace_Test(namespace string) {
	str := GeneratePodNamespaceUrl(namespace)
	resp := InvokeGetReuqest(destinationServer_Test + str)
	if (resp != nil) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
			return
		}
		var v classType.PodList
		jsonParse.JsonUnmarsha(body, &v)
		for _, item := range v.Items {
			classType.PrintPod(item)
		}
	}
}

func PostUrl_test(url string, byte []byte) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(byte))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
func DeleteUrl_test(url string, byte []byte) {
	fmt.Print(url)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func CreatePod_test(namespace string, image string, name string, cpu_min string, cpu_max string, mem_min string, mem_max string) {
	var resource classType.ResourceRequirements
	resource.Limits = make(map[classType.ResourceName]string)
	resource.Limits["cpu"] = cpu_max
	resource.Limits["memory"] = mem_max
	resource.Requests = make(map[classType.ResourceName]string)
	resource.Requests["cpu"] = cpu_min
	resource.Requests["memory"] = mem_min
	byte := GeneratePodBody(namespace, image, name, resource)
	url := destinationServer_Test + GeneratePodNamespaceUrl(namespace)
	fmt.Print(url + "\n")
	PostUrl_test(url, byte)
}

func CreateService_test(name string, label_name string, namespace string, port int32) {
	byte := GenerateServiceBody(name, label_name, namespace, port)
	url := destinationServer_Test + GenerateServiceListNamespaceUrl(namespace)
	fmt.Print(url + "\n")
	PostUrl_test(url, byte)
}

func CreateReplicationController_test(namespace string, image string, name string, podName string, labelName string, replic int32, cpu_min string, cpu_max string, mem_min string, mem_max string) {
	var resource classType.ResourceRequirements
	resource.Limits = make(map[classType.ResourceName]string)
	resource.Limits["cpu"] = cpu_max
	resource.Limits["memory"] = mem_max
	resource.Requests = make(map[classType.ResourceName]string)
	resource.Requests["cpu"] = cpu_min
	resource.Requests["memory"] = mem_min
	byte := GenerateReplicationcontrollerBody(namespace, image, name, podName, labelName, replic, resource)
	url := destinationServer_Test + GenerateReplicationControllerNamespaceUrl(namespace)
	fmt.Print(url + "\n")
	PostUrl_test(url, byte)
}

func DeletePod(namespace string, name string) {
	url := destinationServer_Test + GeneratePodNamespaceUrl(namespace) + "/" + name
	DeleteUrl_test(url, nil)
}

func DeleteService(namespace string, name string) {
	url := destinationServer_Test + GenerateServiceListNamespaceUrl(namespace) + "/" + name
	DeleteUrl_test(url, nil)
}

func DeleteReplicationController(namespace string, name string) {
	url := destinationServer_Test + GenerateReplicationControllerNamespaceUrl(namespace) + "/" + name
	DeleteUrl_test(url, nil)
}

func GeneratePodBody(namespace string, image string, name string, resource classType.ResourceRequirements) []byte {
	//生成typeMedata
	var typeMedata classType.TypeMeta
	typeMedata.APIVersion = "v1"
	typeMedata.Kind = "Pod"
	//生成objectMedata
	var objectMedata classType.ObjectMeta
	objectMedata.Labels = make(map[string]string)
	objectMedata.Labels["name"] = name
	objectMedata.Namespace = namespace
	objectMedata.Name = name
	//生成spec.container
	var container classType.Container
	container.Name = name
	container.Image = image
	container.Resources = resource
	var containers [1]classType.Container
	containers[0] = container
	var pod classType.Pod
	pod.TypeMeta = typeMedata
	pod.ObjectMeta = objectMedata
	//将container拷到pod.spec.containers
	slice := []classType.Container{container}
	pod.Spec.Containers = slice
	b := jsonParse.JsonMarsha(pod)
	fmt.Print(string(b))
	return b

}

func GenerateServiceBody(name string, labelName string, namespace string, port int32) []byte {
	//生成typeMedata
	var typeMedata classType.TypeMeta
	typeMedata.APIVersion = "v1"
	typeMedata.Kind = "Service"

	//生成objectMedata
	var objectMedata classType.ObjectMeta
	objectMedata.Labels = make(map[string]string)
	objectMedata.Labels["name"] = name
	objectMedata.Namespace = namespace
	objectMedata.Name = name

	//生成Service spec
	var servicePort classType.ServicePort
	servicePort.Port = port
	slice := []classType.ServicePort{servicePort}
	var serviceSpec classType.ServiceSpec
	serviceSpec.Selector = make(map[string]string)
	serviceSpec.Selector["name"] = labelName
	serviceSpec.Ports = slice
	var service classType.Service
	service.ObjectMeta = objectMedata
	service.TypeMeta = typeMedata
	service.Spec = serviceSpec
	b := jsonParse.JsonMarsha(service)
	fmt.Print(string(b))
	return b
}

func GenerateReplicationcontrollerBody(namespace string, image string, name string, podName string, labelName string, replic int32, resource classType.ResourceRequirements) []byte {
	//生成typeMedata
	var typeMedata classType.TypeMeta
	typeMedata.APIVersion = "v1"
	typeMedata.Kind = "ReplicationController"
	//生成objectMedata
	var objectMedata classType.ObjectMeta
	objectMedata.Labels = make(map[string]string)
	objectMedata.Labels["name"] = labelName
	objectMedata.Namespace = namespace
	objectMedata.Name = name

	//生成PodObjectMedata
	var objectMedata2 classType.ObjectMeta
	objectMedata2.Labels = make(map[string]string)
	objectMedata2.Labels["name"] = labelName
	objectMedata2.Namespace = namespace
	objectMedata2.Name = podName

	//生成PodTemplateSpec.container
	var container classType.Container
	container.Name = name
	container.Image = image
	container.Resources = resource
	var containers [1]classType.Container
	containers[0] = container
	slice := []classType.Container{container}
	var podTemplateSpec classType.PodTemplateSpec
	podTemplateSpec.ObjectMeta = objectMedata2
	podTemplateSpec.Spec.Containers = slice

	var replicationControllerSpec classType.ReplicationControllerSpec
	replicationControllerSpec.Template = &podTemplateSpec

	replicationControllerSpec.Replicas = &replic
	replicationControllerSpec.Selector = make(map[string]string)
	replicationControllerSpec.Selector["name"] = labelName

	var replicationController classType.ReplicationController
	replicationController.Spec = replicationControllerSpec
	replicationController.ObjectMeta = objectMedata
	replicationController.TypeMeta = typeMedata

	b := jsonParse.JsonMarsha(replicationController)
	fmt.Print(string(b))
	return b

}



