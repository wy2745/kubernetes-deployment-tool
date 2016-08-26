package request

import(
	"net/http"
	jsonParse "../json"
	"io/ioutil"
	"fmt"
)

func GetRequest(url ,parameter string) {

	if(parameter == ""){
		parameter = "/"
	}
	resp,err := http.Get(url+parameter)
	if err != nil {
		return
	}
	fmt.Print(url+parameter+"\n")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	v := jsonParse.JsonUnmarsha(body)
	fmt.Print("kind: ")
	fmt.Print(v.Kind)
	fmt.Print("\n")
	fmt.Print("nodeInfo: ")
	fmt.Print(v.NodeInfo)
	fmt.Print("\n")
	fmt.Print("apiversion: ")
	fmt.Print(v.ApiVersion)
	fmt.Print("\n")
	fmt.Print("da: ")
	fmt.Print(v.DaemonEndpoints)
	fmt.Print("\n")
	fmt.Print("ite: ")
	fmt.Print(v.Items)
	fmt.Print("\n")
	fmt.Print("me: ")
	fmt.Print(v.Metadata)
	fmt.Print("\n")
	//fmt.Print(obj)
	//for k,v := range obj.(map[string]interface{}) {
	//	fmt.Print(k)
	//	fmt.Print("---key-------\n")
	//	fmt.Print(v)
	//	fmt.Print("---value-------\n")
	//}
}

