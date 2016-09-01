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
	str := GeneratePodListNamespaceUrl(namespace)
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

func CreatePod_test(namespace string, image string, name string) {
	byte := GeneratePodBody(namespace, image, name)
	url := destinationServer_Test + GeneratePodListNamespaceUrl(namespace)
	fmt.Print(url + "\n")
	PostUrl_test(url, byte)
	//req, err := http.NewRequest("POST", url, bytes.NewBuffer(byte))
	//if err != nil {
	//	panic(err)
	//}
	//req.Header.Set("Content-Type", "application/json")
	//
	//client := &http.Client{}
	//resp, err := client.Do(req)
	//if err != nil {
	//	panic(err)
	//}
	//defer resp.Body.Close()
	//
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))
}

func CreateService_test(name string, label_name string, namespace string, port int32) {
	byte := GenerateServiceBody(name, label_name, namespace, port)
	url := destinationServer_Test + GenerateServiceListNamespaceUrl(namespace)
	fmt.Print(url + "\n")
	PostUrl_test(url, byte)
	//req, err := http.NewRequest("POST", url, bytes.NewBuffer(byte))
	//if err != nil {
	//	panic(err)
	//}
	//req.Header.Set("Content-Type", "application/json")
	//
	//client := &http.Client{}
	//resp, err := client.Do(req)
	//if err != nil {
	//	panic(err)
	//}
	//defer resp.Body.Close()
	//
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))
}

func GeneratePodBody(namespace string, image string, name string) []byte {
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

func GenerateServiceBody(name string, label_name string, namespace string, port int32) []byte {
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
	serviceSpec.Selector["name"] = label_name
	serviceSpec.Ports = slice
	var service classType.Service
	service.ObjectMeta = objectMedata
	service.TypeMeta = typeMedata
	service.Spec = serviceSpec
	b := jsonParse.JsonMarsha(service)
	fmt.Print(string(b))
	return b
}

func Test() {
	var abc = "=-="
	var str = `"title":"haha"`
	var cde = `"ld":` + `"` + abc + `",` + str
	var dfg = `{` + cde + `}`
	fmt.Print(dfg)
}

//func GetRequest(url, parameter string) interface{} {
//
//	if (parameter == "") {
//		parameter = "/"
//	}
//	resp, err := http.Get(url + parameter)
//	if err != nil {
//		return nil
//	}
//	fmt.Print(url + parameter + "\n")
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	v := jsonParse.JsonUnmarsha(body)
//	return v
//}
//apiVersion: v1
//kind: Service
//metadata:
//labels:
//name: mysql
//name: mysql
//spec:
//ports:
//- port: 3306
//selector:
//name: mysql
