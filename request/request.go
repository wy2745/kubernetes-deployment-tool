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
	fmt.Print(v.Items)
	fmt.Print(v.APIVersion)
	fmt.Print(v.Kind)
	fmt.Print(v.SelfLink)
	fmt.Print(v.ListMeta.ResourceVersion)
	fmt.Print(v.ListMeta.SelfLink)
	fmt.Print(v.TypeMeta.APIVersion)
	fmt.Print(v.TypeMeta.Kind)
	//fmt.Print(obj)
	//for k,v := range obj.(map[string]interface{}) {
	//	fmt.Print(k)
	//	fmt.Print("---key-------\n")
	//	fmt.Print(v)
	//	fmt.Print("---value-------\n")
	//}
}

