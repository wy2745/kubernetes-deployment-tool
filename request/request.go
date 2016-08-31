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
			classType.PringNamespace(item)
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
func CreatePod_test(namespace string, image string, name string) {
	byte := GenetatePodBody(namespace, image, name)
	url := destinationServer_Test + GeneratePodListNamespaceUrl(namespace)
	fmt.Print(url + "\n")
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

func GenetatePodBody(namespace string, image string, name string) []byte {
	var typeMedata classType.TypeMeta
	typeMedata.APIVersion = "v1"
	typeMedata.Kind = "Pod"
	var objectMedata classType.ObjectMeta
	objectMedata.Labels = make(map[string]string)
	objectMedata.Labels["name"] = name
	objectMedata.Namespace = namespace
	objectMedata.Name = name
	var container classType.Container
	container.Name = name
	container.Image = image
	var containers [1]classType.Container
	containers[0] = container
	var pod classType.Pod
	pod.TypeMeta = typeMedata
	pod.ObjectMeta = objectMedata
	slice := []classType.Container{container}
	pod.Spec.Containers = slice
	b := jsonParse.JsonMarsha(pod)
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
//kind: Pod
//metadata:
//name: mysql
//labels:
//name: mysql
//spec:
//containers:
//- image: mysql:5.7
//name: mysql
//env:
//- name: MYSQL_ROOT_PASSWORD
//value: '123456'
//      ports:
//- containerPort: 3306
//hostPort: 3306
//name: mysql
//volumeMounts:
//- mountPath: /var/lib/mysql
//name: data
//volumes:
//- name: data
//hostPath:
//path: /home/administrator/data
