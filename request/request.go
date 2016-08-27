package request

import (
	"net/http"
	jsonParse "../json"
	classType "../type"
	"io/ioutil"
	"fmt"
)

func InvokeGetReuqest(url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Print(err)
		return nil
	}
	return resp
}
func GetAllNode() {
	resp := InvokeGetReuqest(destinationServer + GetNodeList_GET)
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
func GetAllNamespace() {
	resp := InvokeGetReuqest(destinationServer + GetNamespaces_GET)
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

