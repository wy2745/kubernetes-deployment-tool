package request

import (
	"net/http"
	jsonParse "../json"
	classType "../type"
	"io/ioutil"
	"fmt"
	"bytes"
	"strings"
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
	fmt.Println(resp.Header)
	fmt.Println(resp.Status)
	fmt.Println(resp.StatusCode)
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
		var v classType.NodeList
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
		var v classType.NamespaceList
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
		var v classType.ReplicationControllerList
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
		var v classType.JobList
		jsonParse.JsonUnmarsha(body, &v)
		fmt.Println(v.APIVersion)
		fmt.Println(v.SelfLink)
		fmt.Println(v.ResourceVersion)
		for _, item := range v.Items {
			classType.PrintJob(item)
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
		var v classType.JobList
		jsonParse.JsonUnmarsha(body, &v)
		for _, item := range v.Items {
			classType.PrintJob(item)
		}
	}
}
func GetJobByNameAndNamespace(namespace string, name string, mode string) {
	var resp *http.Response
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
			return
		}
		var v classType.JobList
		jsonParse.JsonUnmarsha(body, &v)
		for _, item := range v.Items {
			classType.PrintJob(item)
		}
	}
}
func GetAllService(mode string) {
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
			return
		}
		var v classType.ServiceList
		jsonParse.JsonUnmarsha(body, &v)
		for _, item := range v.Items {
			classType.PrintService(item)
		}
	}
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
		var v classType.ServiceList
		jsonParse.JsonUnmarsha(body, &v)
		for _, item := range v.Items {
			classType.PrintService(item)
		}
	}
}

func GetPodsOfNamespace(namespace string, mode string) {
	var resp *http.Response
	str := GeneratePodNamespaceUrl(namespace)
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
		var v classType.PodList
		jsonParse.JsonUnmarsha(body, &v)
		//for _, item := range v.Items {
		//	classType.PrintPod(item)
		//}
	}
}

func GetPodByNameAndNamespace(namespace string, name string, mode string) classType.Pod {
	var v classType.Pod
	str := GeneratePodNameUrl(namespace, name)

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
			return v
		}
		jsonParse.JsonUnmarsha(body, &v)
		//classType.PrintPod(v)
		return v
	}
	return v
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

func CreatePod(namespace string, image string, name string, cpu_min string, cpu_max string, mem_min string, mem_max string, command string, mode string) {
	var resource classType.ResourceRequirements
	resource.Limits = make(map[classType.ResourceName]string)
	resource.Limits["cpu"] = cpu_max
	resource.Limits["memory"] = mem_max
	resource.Requests = make(map[classType.ResourceName]string)
	resource.Requests["cpu"] = cpu_min
	resource.Requests["memory"] = mem_min
	byte := GeneratePodBody(namespace, image, name, resource, command)
	var url string
	if mode == Test {
		url = destinationServer_Test + GeneratePodNamespaceUrl(namespace)
		PostUrl_test(url, byte)
	} else {
		url = destinationServer_Caicloud + GeneratePodNamespaceUrl(namespace)
		InvokeRequest_Caicloud("POST", url, byte)
	}

	//fmt.Print(url + "\n")

}
func CreateJob(namespace string, image string, name string, completion int32, parallelism int32, cpu_min string, cpu_max string, mem_min string, mem_max string, command string, mode string) {
	var resource classType.ResourceRequirements
	resource.Limits = make(map[classType.ResourceName]string)
	resource.Limits["cpu"] = cpu_max
	resource.Limits["memory"] = mem_max
	resource.Requests = make(map[classType.ResourceName]string)
	resource.Requests["cpu"] = cpu_min
	resource.Requests["memory"] = mem_min
	byte := GenerateJobBody(namespace, name, image, completion, parallelism, command, resource)
	var url string
	if mode == Test {
		url = destinationServer_Test + GenerateJobNamespaceUrl(namespace)
		PostUrl_test(url, byte)
	} else {
		url = destinationServer_Caicloud + GenerateJobNamespaceUrl(namespace)
		InvokeRequest_Caicloud("POST", url, byte)
	}

	//fmt.Print(url + "\n")
	//fmt.Println(string(byte))

}

func CreateService(name string, label_name string, namespace string, port int32, mode string) {
	byte := GenerateServiceBody(name, label_name, namespace, port)

	var url string
	if mode == Test {
		url = destinationServer_Test + GenerateServiceListNamespaceUrl(namespace)
		PostUrl_test(url, byte)
	} else {
		url = destinationServer_Caicloud + GenerateServiceListNamespaceUrl(namespace)
		InvokeRequest_Caicloud("POST", url, byte)
	}

}

func CreateReplicationController(namespace string, image string, name string, podName string, labelName string, replic int32, cpu_min string, cpu_max string, mem_min string, mem_max string, mode string) {
	var resource classType.ResourceRequirements
	resource.Limits = make(map[classType.ResourceName]string)
	resource.Limits["cpu"] = cpu_max
	resource.Limits["memory"] = mem_max
	resource.Requests = make(map[classType.ResourceName]string)
	resource.Requests["cpu"] = cpu_min
	resource.Requests["memory"] = mem_min
	byte := GenerateReplicationcontrollerBody(namespace, image, name, podName, labelName, replic, resource)

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
func GenerateJobBody(namespace string, name string, image string, completion int32, Parallelism int32, command string, resource classType.ResourceRequirements) []byte {
	var typeMeta classType.TypeMeta
	typeMeta.APIVersion = "batch/v1"
	typeMeta.Kind = "Job"
	var objectMeta classType.ObjectMeta
	objectMeta.Name = name
	objectMeta.Namespace = namespace

	var jobSpec classType.JobSpec
	jobSpec.Parallelism = &Parallelism
	jobSpec.Completions = &completion
	jobSpec.Template.Name = name
	var container classType.Container
	container.Name = name
	container.Image = image
	container.Resources = resource
	container.Command = strings.Split(command, " ")
	var containers [1]classType.Container
	containers[0] = container

	jobSpec.Template.Name = name
	slice := []classType.Container{container}
	jobSpec.Template.Spec.Containers = slice
	var job classType.Job
	job.TypeMeta = typeMeta
	job.ObjectMeta = objectMeta
	job.Spec = jobSpec
	job.Spec.Template.Spec.RestartPolicy = "Never"
	b := jsonParse.JsonMarsha(job)

	//fmt.Println("jaja")
	//fmt.Println(string(b))
	return b

}

func GeneratePodBody(namespace string, image string, name string, resource classType.ResourceRequirements, command string) []byte {
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
	container.Command = strings.Split(command, " ")
	var containers [1]classType.Container
	containers[0] = container
	var pod classType.Pod
	pod.TypeMeta = typeMedata
	pod.ObjectMeta = objectMedata
	//将container拷到pod.spec.containers
	slice := []classType.Container{container}
	pod.Spec.Containers = slice
	pod.Spec.RestartPolicy = "Always"
	b := jsonParse.JsonMarsha(pod)
	//fmt.Print(string(b))
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
	//fmt.Print(string(b))
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
	//fmt.Print(string(b))
	return b

}
func PodComplete(pod classType.Pod) bool {
	//fmt.Println("pod寻找结果：", len(pod.Name) == 0)
	if len(pod.Status.ContainerStatuses) != 0 {
		//fmt.Println("ready:", pod.Status.ContainerStatuses[0].Ready)
		//fmt.Println("status:", pod.Status.Phase)
		//fmt.Println("结果:", len(pod.Name) == 0 || (pod.Status.ContainerStatuses[0].Ready == false && pod.Status.Phase == "Running"))
		return len(pod.Name) == 0 || (pod.Status.ContainerStatuses[0].Ready == false && pod.Status.Phase == "Running")
	}
	return len(pod.Name) == 0
}



