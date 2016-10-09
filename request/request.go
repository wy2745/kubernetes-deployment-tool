package request

import (
	"net/http"
	jsonParse "../json"
	classType1 "../type124"
	classType2 "../type137"
	"io/ioutil"
	"fmt"
	"bytes"
	"strings"
	"../interf"
)

const (
	Test string = "test"
	Caicloud string = "caicloud"
)

func setBasicAuthOfCaicloud(r *http.Request) {
	r.SetBasicAuth(userName, password)
}

func InvokeRequest_Caicloud(method string, url string, body []byte) *http.Response {
	client := http.Client{}
	var req *http.Request
	var err error
	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	setBasicAuthOfCaicloud(req)
	if method != "GET" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := client.Do(req)
	//fmt.Println(resp.Header)
	//fmt.Println(resp.Status)
	//fmt.Println(resp.StatusCode)
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

func GetAllNode(mode string) {
	var resp *http.Response
	if mode == Test {
		resp = InvokeGetReuqest(destinationServer_Test + GetNodeList_GET)
	} else {
		resp = InvokeRequest_Caicloud("GET", destinationServer_Caicloud + GetNodeList_GET, nil)
	}
	if (resp != nil) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
			return
		}
		var v classType1.NodeList
		jsonParse.JsonUnmarsha(body, &v)
		//for _, item := range v.Items {
		//	classType.PrintNode(item)
		//}
	}
}

func GetAllNamespace(mode string) {
	var resp *http.Response
	if mode == Test {
		resp = InvokeGetReuqest(destinationServer_Test + GetNamespaces_GET)
	} else {
		resp = InvokeRequest_Caicloud("GET", destinationServer_Caicloud + GetNamespaces_GET, nil)
	}
	if (resp != nil) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
			return
		}
		var v classType1.NamespaceList
		jsonParse.JsonUnmarsha(body, &v)
		//for _, item := range v.Items {
		//	classType.PrintNamespace(item)
		//}
	}
}

func GetAllReplicationcontrollers(mode string) {
	var resp *http.Response
	if mode == Test {
		resp = InvokeGetReuqest(destinationServer_Test + ReadAllReplicationController_GET)
	} else {
		resp = InvokeRequest_Caicloud("GET", destinationServer_Caicloud + ReadAllReplicationController_GET, nil)
	}
	if (resp != nil) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
			return
		}
		var v classType1.ReplicationControllerList
		jsonParse.JsonUnmarsha(body, &v)
		//for _, item := range v.Items {
		//	classType.PrintReplicationController(item)
		//}
	}
}

func GetAllJobs(mode string) {
	var resp *http.Response
	if mode == Test {
		resp = InvokeGetReuqest(destinationServer_Test + ReadAllJob_GET)
	} else {
		resp = InvokeRequest_Caicloud("GET", destinationServer_Caicloud + ReadAllJob_GET, nil)
	}
	if (resp != nil) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
			return
		}
		var v classType1.JobList
		jsonParse.JsonUnmarsha(body, &v)
		fmt.Println(v.APIVersion)
		fmt.Println(v.SelfLink)
		fmt.Println(v.ResourceVersion)
		for _, item := range v.Items {
			classType1.PrintJob(item)
		}
	}
}

func GetJobOfNamespace(namespace string, mode string) {
	var resp *http.Response
	str := GenerateJobNamespaceUrl(namespace)
	if mode == Test {
		resp = InvokeGetReuqest(destinationServer_Test + str)
	} else {
		resp = InvokeRequest_Caicloud("GET", destinationServer_Caicloud + str, nil)
	}
	if (resp != nil) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
			return
		}
		var v classType1.JobList
		jsonParse.JsonUnmarsha(body, &v)
		for _, item := range v.Items {
			classType1.PrintJob(item)
		}
	}
}
func JobExist(namespace string, name string, mode string) bool {
	var resp *http.Response
	var v classType1.Job
	str := GenerateJobNameUrl(namespace, name)
	if mode == Test {
		resp = InvokeGetReuqest(destinationServer_Test + str)
	} else {
		resp = InvokeRequest_Caicloud("GET", destinationServer_Caicloud + str, nil)
	}
	if (resp != nil) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			return false
		}
		if (resp.Status == "404 Not Found") {
			return false
		}
		jsonParse.JsonUnmarsha(body, &v)

		return true
		//for _, item := range v.Items {
		//	classType.PrintJob(item)
		//}
	}
	return false
}

func GetJobByNameAndNamespace(namespace string, name string, mode string) classType1.Job {
	var resp *http.Response
	var v classType1.Job
	str := GenerateJobNameUrl(namespace, name)
	if mode == Test {
		resp = InvokeGetReuqest(destinationServer_Test + str)
	} else {
		resp = InvokeRequest_Caicloud("GET", destinationServer_Caicloud + str, nil)
	}
	if (resp != nil) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
			return v
		}
		jsonParse.JsonUnmarsha(body, &v)
		//for _, item := range v.Items {
		//	classType.PrintJob(item)
		//}
	}
	return v
}
func GetAllService(mode string) []classType1.Service {
	var v classType1.ServiceList
	var serviceList []classType1.Service
	var resp *http.Response
	if mode == Test {
		resp = InvokeGetReuqest(destinationServer_Test + ReadAllService_GET)
	} else {
		resp = InvokeRequest_Caicloud("GET", destinationServer_Caicloud + ReadAllService_GET, nil)
	}

	if (resp != nil) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
			return serviceList
		}

		jsonParse.JsonUnmarsha(body, &v)
		for _, item := range v.Items {
			classType1.PrintService(item)
		}
		serviceList = v.Items
		return serviceList
	}
	return serviceList
}
func GetServicesOfNamespace(namespace string, mode string) {
	var resp *http.Response
	str := GenerateServiceListNamespaceUrl(namespace)
	if mode == Test {
		resp = InvokeGetReuqest(destinationServer_Test + str)
	} else {
		resp = InvokeRequest_Caicloud("GET", destinationServer_Caicloud + str, nil)
	}

	if (resp != nil) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
			return
		}
		var v classType1.ServiceList
		jsonParse.JsonUnmarsha(body, &v)
		for _, item := range v.Items {
			classType1.PrintService(item)
		}
	}
}
func GetServicesOfNamespaceAndName(namespace string, name string, mode string) classType1.Service {
	var v classType1.Service
	var resp *http.Response
	str := GenerateServiceListNameUrl(namespace, name)
	if mode == Test {
		resp = InvokeGetReuqest(destinationServer_Test + str)
	} else {
		resp = InvokeRequest_Caicloud("GET", destinationServer_Caicloud + str, nil)
	}

	if (resp != nil) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
			return v
		}
		jsonParse.JsonUnmarsha(body, &v)
		return v
	}
	return v
}

func GetPodsOfNamespace(namespace string, mode string) []*interf.Podinface {
	var v = getPodListtypeByMode(mode)
	var resp *http.Response
	str := GeneratePodNamespaceUrl(namespace)
	if mode == Test {
		resp = InvokeGetReuqest(destinationServer_Test + str)
	} else {
		resp = InvokeRequest_Caicloud("GET", destinationServer_Caicloud + str, nil)
	}
	pl := jsonParse.JsonUnmarshaPodList(resp, &v)
	return interf.GetItems(*pl)

	//if (resp != nil) {
	//	defer resp.Body.Close()
	//	body, err := ioutil.ReadAll(resp.Body)
	//	if (err != nil) {
	//		fmt.Print(err)
	//		return pod
	//	}
	//	jsonParse.JsonUnmarsha(body, &v)
	//	pod = v.Items
	//	return pod
	//	//for _, item := range v.Items {
	//	//	classType.PrintPod(item)
	//	//}
	//}
}

func GetPodByNameAndNamespace(namespace string, name string, mode string) *interf.Podinface {

	v := getPodtypeByMode(mode)
	str := GeneratePodNameUrl(namespace, name)

	var resp *http.Response
	if mode == Test {
		resp = InvokeGetReuqest(destinationServer_Test + str)
	} else {
		resp = InvokeRequest_Caicloud("GET", destinationServer_Caicloud + str, nil)
	}

	return jsonParse.JsonUnmarshaPod(resp, &v)

	//if (resp != nil) {
	//	defer resp.Body.Close()
	//	body, err := ioutil.ReadAll(resp.Body)
	//	if (err != nil) {
	//		fmt.Print(err)
	//		return v
	//	}
	//	jsonParse.JsonUnmarsha(body, &v)
	//	//classType.PrintPod(v)
	//	return v
	//}
	//return v
}

func GetNodeByName(name string, mode string) (classType1.Node, string) {
	var v classType1.Node
	var str2 string
	str := GenerateNodeNameUrl(name)

	var resp *http.Response
	if mode == Test {
		resp = InvokeGetReuqest(destinationServer_Test + str)
	} else {
		resp = InvokeRequest_Caicloud("GET", destinationServer_Caicloud + str, nil)
	}

	if (resp != nil) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
			return v, str2
		}
		jsonParse.JsonUnmarsha(body, &v)
		str2 = classType1.PrintNodeResourceStatus(v)
		return v, str2
	}
	return v, str2
}

func PostUrl_test(url string, byte []byte) *http.Response {
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
	return resp
	//defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	//ioutil.ReadAll(resp.Body)
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))
}

func DeleteUrl_test(url string, byte []byte) {
	//fmt.Print(url)
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

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))
	ioutil.ReadAll(resp.Body)
}

func CreatePod(namespace string, image string, name string, cpu_min string, cpu_max string, mem_min string, mem_max string, command string, mode string) *interf.Podinface {

	//var v = getResourceByMode(mode)
	//interf.SetResource(v, cpu_min, cpu_max, mem_min, mem_max)
	//var resource classType1.ResourceRequirements
	//resource.Limits = make(map[classType1.ResourceName]string)
	//resource.Limits["cpu"] = cpu_max
	//resource.Limits["memory"] = mem_max
	//resource.Requests = make(map[classType1.ResourceName]string)
	//resource.Requests["cpu"] = cpu_min
	//resource.Requests["memory"] = mem_min
	byte := GeneratePodBody(namespace, image, name, cpu_max, mem_max, command, mode)
	var resp *http.Response
	var url string
	if mode == Test {
		url = destinationServer_Test + GeneratePodNamespaceUrl(namespace)
		resp = PostUrl_test(url, byte)
	} else {
		url = destinationServer_Caicloud + GeneratePodNamespaceUrl(namespace)
		resp = InvokeRequest_Caicloud("POST", url, byte)
	}
	var v2 = getPodtypeByMode(mode)
	return jsonParse.JsonUnmarshaPod(resp, &v2)
	//if (resp != nil) {
	//	defer resp.Body.Close()
	//	body, err := ioutil.ReadAll(resp.Body)
	//	if (err != nil) {
	//		fmt.Print(err)
	//		return v
	//	}
	//	jsonParse.JsonUnmarsha(body, &v)
	//	return v
	//}
	//return v

}
func CreateJob(namespace string, image string, name string, completion int32, cpu_min string, cpu_max string, mem_min string, mem_max string, command string, mode string) {
	var resource classType1.ResourceRequirements
	resource.Limits = make(map[classType1.ResourceName]string)
	resource.Limits["cpu"] = cpu_max
	resource.Limits["memory"] = mem_max
	resource.Requests = make(map[classType1.ResourceName]string)
	resource.Requests["cpu"] = cpu_min
	resource.Requests["memory"] = mem_min
	byte := GenerateJobBody(namespace, name, image, completion, command, resource, true)
	var url string
	if mode == Test {
		url = destinationServer_Test + GenerateJobNamespaceUrl(namespace)
		PostUrl_test(url, byte)
	} else {
		url = destinationServer_Caicloud + GenerateJobNamespaceUrl(namespace)
		InvokeRequest_Caicloud("POST", url, byte)
	}
}

func CreateJobWithoutResource(namespace string, image string, name string, completion int32, command string, mode string) {
	var resource classType1.ResourceRequirements
	byte := GenerateJobBody(namespace, name, image, completion, command, resource, false)
	var url string
	if mode == Test {
		url = destinationServer_Test + GenerateJobNamespaceUrl(namespace)
		PostUrl_test(url, byte)
	} else {
		url = destinationServer_Caicloud + GenerateJobNamespaceUrl(namespace)
		InvokeRequest_Caicloud("POST", url, byte)
	}
}

func CreateService(name string, label_name string, namespace string, port int32, nodeport int32, mode string) {
	byte := GenerateServiceBody(name, label_name, namespace, port, nodeport)

	var url string
	if mode == Test {
		url = destinationServer_Test + GenerateServiceListNamespaceUrl(namespace)
		PostUrl_test(url, byte)
	} else {
		url = destinationServer_Caicloud + GenerateServiceListNamespaceUrl(namespace)
		InvokeRequest_Caicloud("POST", url, byte)
	}

}

func CreateReplicationController(namespace string, image string, name string, podName string, labelName string, replic int32, cpu_min string, cpu_max string, mem_min string, mem_max string, ports map[int32]int32, mode string) {
	var resource classType1.ResourceRequirements
	resource.Limits = make(map[classType1.ResourceName]string)
	resource.Limits["cpu"] = cpu_max
	resource.Limits["memory"] = mem_max
	resource.Requests = make(map[classType1.ResourceName]string)
	resource.Requests["cpu"] = cpu_min
	resource.Requests["memory"] = mem_min
	var Ports []classType1.ContainerPort
	for key, value := range ports {
		var port classType1.ContainerPort
		port.ContainerPort = key
		port.HostPort = value
		Ports = append(Ports, port)
	}
	byte := GenerateReplicationcontrollerBody(namespace, image, name, podName, labelName, replic, resource, Ports)

	var url string
	if mode == Test {
		url = destinationServer_Test + GenerateReplicationControllerNamespaceUrl(namespace)
		PostUrl_test(url, byte)
	} else {
		url = destinationServer_Caicloud + GenerateReplicationControllerNamespaceUrl(namespace)
		InvokeRequest_Caicloud("POST", url, byte)
	}

}

func DeletePod(namespace string, name string, mode string) {
	var url string
	if mode == Test {
		url = destinationServer_Test + GeneratePodNameUrl(namespace, name)
		DeleteUrl_test(url, nil)
	} else {
		url = destinationServer_Caicloud + GeneratePodNameUrl(namespace, name)
		InvokeRequest_Caicloud("DELETE", url, nil)
	}

}

func DeleteJob(namespace string, name string, mode string) {
	var url string
	if mode == Test {
		url = destinationServer_Test + GenerateJobNameUrl(namespace, name)
		DeleteUrl_test(url, nil)
	} else {
		url = destinationServer_Caicloud + GenerateJobNameUrl(namespace, name)
		InvokeRequest_Caicloud("DELETE", url, nil)
	}
}

func DeleteService(namespace string, name string, mode string) {
	var url string
	if mode == Test {
		url = destinationServer_Test + GenerateServiceListNamespaceUrl(namespace) + "/" + name
		DeleteUrl_test(url, nil)
	} else {
		url = destinationServer_Caicloud + GenerateServiceListNamespaceUrl(namespace) + "/" + name
		InvokeRequest_Caicloud("DELETE", url, nil)
	}

}

func DeleteReplicationController(namespace string, name string, mode string) {
	var url string
	if mode == Test {
		url = destinationServer_Test + GenerateReplicationControllerNamespaceUrl(namespace) + "/" + name
		DeleteUrl_test(url, nil)
	} else {
		url = destinationServer_Caicloud + GenerateReplicationControllerNamespaceUrl(namespace) + "/" + name
		InvokeRequest_Caicloud("DELETE", url, nil)
	}

}
func GenerateJobBody(namespace string, name string, image string, completion int32, command string, resource classType1.ResourceRequirements, whetherResource bool) []byte {
	var typeMeta classType1.TypeMeta
	typeMeta.APIVersion = "batch/v1"
	typeMeta.Kind = "Job"
	var objectMeta classType1.ObjectMeta
	objectMeta.Name = name
	objectMeta.Namespace = namespace

	var jobSpec classType1.JobSpec
	jobSpec.Completions = &completion
	jobSpec.Template.Name = name
	var container classType1.Container
	container.Name = name
	container.Image = image
	if whetherResource == true {
		container.Resources = resource
	}

	container.Command = strings.Split(command, " ")
	var containers [1]classType1.Container
	containers[0] = container

	jobSpec.Template.Name = name
	slice := []classType1.Container{container}
	jobSpec.Template.Spec.Containers = slice
	var job classType1.Job
	job.TypeMeta = typeMeta
	job.ObjectMeta = objectMeta
	job.Spec = jobSpec
	job.Spec.Template.Spec.RestartPolicy = "Never"
	b := jsonParse.JsonMarsha(job)

	//fmt.Println("jaja")
	//fmt.Println(string(b))
	return b

}

func GeneratePodBody(namespace string, image string, name string, cpu string, mem string, command string, mode string) []byte {
	var pod = getPodtypeByMode(mode)
	b := interf.SetPod(pod, "v1", "Pod", name, namespace, image, cpu, mem, command)
	//interf.SetTypeMeta(&pod, "v1", "Pod")
	//interf.SetObjectMeta(&pod, name, namespace, name)
	//interf.SetContainer(&pod, name, image, cpu, mem, command, "Always")
	//var typeMedata classType1.TypeMeta
	//typeMedata.APIVersion = "v1"
	//typeMedata.Kind = "Pod"
	//生成objectMedata
	//var objectMedata classType1.ObjectMeta
	//objectMedata.Labels = make(map[string]string)
	//objectMedata.Labels["name"] = name
	//objectMedata.Namespace = namespace
	//objectMedata.Name = name
	//生成spec.container
	//var container classType1.Container
	//container.Name = name
	//container.Image = image
	//container.Resources = *resource
	//container.Command = strings.Split(command, " ")
	//var containers [1]classType1.Container
	//containers[0] = container
	////
	////pod.TypeMeta = typeMedata
	////pod.ObjectMeta = objectMedata
	////将container拷到pod.spec.containers
	//slice := []classType1.Container{container}
	//pod.Spec.Containers = slice
	//pod.Spec.RestartPolicy = "Always"
	//b := jsonParse.JsonMarsha(pod)
	//fmt.Print(string(b))
	return b

}

func GenerateServiceBody(name string, labelName string, namespace string, port int32, nodePort int32) []byte {
	//生成typeMedata
	var typeMedata classType1.TypeMeta
	typeMedata.APIVersion = "v1"
	typeMedata.Kind = "Service"

	//生成objectMedata
	var objectMedata classType1.ObjectMeta
	objectMedata.Labels = make(map[string]string)
	objectMedata.Labels["name"] = name
	objectMedata.Namespace = namespace
	objectMedata.Name = name

	//生成Service spec
	var servicePort classType1.ServicePort
	servicePort.Port = port
	servicePort.NodePort = nodePort
	slice := []classType1.ServicePort{servicePort}
	var serviceSpec classType1.ServiceSpec
	serviceSpec.Selector = make(map[string]string)
	serviceSpec.Selector["name"] = labelName
	serviceSpec.Ports = slice
	serviceSpec.Type = classType1.ServiceTypeNodePort
	var service classType1.Service
	service.ObjectMeta = objectMedata
	service.TypeMeta = typeMedata
	service.Spec = serviceSpec
	b := jsonParse.JsonMarsha(service)
	//fmt.Print(string(b))
	return b
}

func GenerateReplicationcontrollerBody(namespace string, image string, name string, podName string, labelName string, replic int32, resource classType1.ResourceRequirements, ports []classType1.ContainerPort) []byte {
	//生成typeMedata
	var typeMedata classType1.TypeMeta
	typeMedata.APIVersion = "v1"
	typeMedata.Kind = "ReplicationController"
	//生成objectMedata
	var objectMedata classType1.ObjectMeta
	objectMedata.Labels = make(map[string]string)
	objectMedata.Labels["name"] = labelName
	objectMedata.Namespace = namespace
	objectMedata.Name = name

	//生成PodObjectMedata
	var objectMedata2 classType1.ObjectMeta
	objectMedata2.Labels = make(map[string]string)
	objectMedata2.Labels["name"] = labelName
	objectMedata2.Namespace = namespace
	objectMedata2.Name = podName

	//生成PodTemplateSpec.container
	var container classType1.Container
	container.Name = name
	container.Image = image
	container.Resources = resource
	container.Ports = ports
	var containers [1]classType1.Container
	containers[0] = container
	slice := []classType1.Container{container}
	var podTemplateSpec classType1.PodTemplateSpec
	podTemplateSpec.ObjectMeta = objectMedata2
	podTemplateSpec.Spec.Containers = slice

	var replicationControllerSpec classType1.ReplicationControllerSpec
	replicationControllerSpec.Template = &podTemplateSpec

	replicationControllerSpec.Replicas = &replic
	replicationControllerSpec.Selector = make(map[string]string)
	replicationControllerSpec.Selector["name"] = labelName

	var replicationController classType1.ReplicationController
	replicationController.Spec = replicationControllerSpec
	replicationController.ObjectMeta = objectMedata
	replicationController.TypeMeta = typeMedata

	b := jsonParse.JsonMarsha(replicationController)
	//fmt.Print(string(b))
	return b

}
func GenerateReplicationcontrollerBodyV2(namespace string, image string, name string, replic int32, command []string) []byte {
	//生成typeMedata
	var typeMedata classType1.TypeMeta
	typeMedata.APIVersion = "v1"
	typeMedata.Kind = "ReplicationController"
	//生成objectMedata
	var objectMedata classType1.ObjectMeta
	objectMedata.Labels = make(map[string]string)
	objectMedata.Labels["name"] = name
	objectMedata.Namespace = namespace
	objectMedata.Name = name

	//生成PodObjectMedata
	var objectMedata2 classType1.ObjectMeta
	objectMedata2.Labels = make(map[string]string)
	objectMedata2.Labels["name"] = name
	objectMedata2.Namespace = namespace

	//生成PodTemplateSpec.container
	var container classType1.Container
	container.Name = name
	container.Image = image
	container.Command = command
	var containers [1]classType1.Container
	containers[0] = container
	slice := []classType1.Container{container}
	var podTemplateSpec classType1.PodTemplateSpec
	podTemplateSpec.ObjectMeta = objectMedata2
	podTemplateSpec.Spec.Containers = slice

	var replicationControllerSpec classType1.ReplicationControllerSpec
	replicationControllerSpec.Template = &podTemplateSpec

	replicationControllerSpec.Replicas = &replic
	replicationControllerSpec.Selector = make(map[string]string)
	replicationControllerSpec.Selector["name"] = name

	var replicationController classType1.ReplicationController
	replicationController.Spec = replicationControllerSpec
	replicationController.ObjectMeta = objectMedata
	replicationController.TypeMeta = typeMedata

	b := jsonParse.JsonMarsha(replicationController)
	//fmt.Print(string(b))
	return b

}
func PodComplete(pod interf.Podinface) bool {
	//fmt.Println("pod寻找结果：", len(pod.Name) == 0)
	if pod.GetContainerStatusesLen() != 0 {
		//fmt.Println("ready:", pod.Status.ContainerStatuses[0].Ready)
		//fmt.Println("status:", pod.Status.Phase)
		//fmt.Println("结果:", len(pod.Name) == 0 || (pod.Status.ContainerStatuses[0].Ready == false && pod.Status.Phase == "Running"))
		return len(pod.GetName()) == 0 || (pod.GetReady() == false && pod.GetStautsPhase() == "Running")
	}
	return len(pod.GetName()) == 0
}
func JobComplete(job classType1.Job) bool {
	return *job.Spec.Completions == job.Status.Succeeded
}

func FindPodByLabelName(name string, pods []*interf.Podinface) (string, string) {
	for _, pod := range pods {
		if (*pod).GetLabel(name) == name {
			return (*pod).GetNamespace(), (*pod).GetName()
		}
	}
	return " ", " "
}
func getPodtypeByMode(mode string) interf.Podinface {

	var v1 classType1.Pod
	var v2 classType2.Pod
	if mode == Caicloud {
		return v1
	} else {
		return v2
	}
	return v2
}
func getPodListtypeByMode(mode string) interf.PodListinface {
	var v1 classType1.PodList
	var v2 classType2.PodList
	if mode == Caicloud {
		return v1
	} else {
		return v2
	}
	return v2
}
func getResourceByMode(mode string) interf.Resourceinface {
	var v1 classType1.ResourceRequirements
	var v2 classType2.ResourceRequirements
	if mode == Caicloud {
		return v1
	} else {
		return v2
	}
	return v2
}
func getTypemetaByMode(mode string) interf.TypeMetainface {
	var v1 classType1.TypeMeta
	var v2 classType2.TypeMeta
	if mode == Caicloud {
		return v1
	} else {
		return v2
	}
	return v2
}
func getObjectmetaByMode(mode string) interf.ObjectMetainface {
	var v1 classType1.ObjectMeta
	var v2 classType2.ObjectMeta
	if mode == Caicloud {
		return v1
	} else {
		return v2
	}
	return v2
}




