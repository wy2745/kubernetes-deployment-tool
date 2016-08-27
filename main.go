package main

import (
	//"fmt"
	//"net/http"
	//"io/ioutil"
	//"encoding/json"
	Request "./request"

	//jsonParse "./json"
	//"fmt"
)

func main() {
	//resp,err := http.Get("http://202.120.40.177:16380/api/v1/nodes")
	//if err != nil {
	//	return
	//}
	//defer resp.Body.Close()
	//
	////使用NewDecoder的方法读取Json,缺点:只能挨个挨个token地读取
	//body, err := ioutil.ReadAll(resp.Body)
	//var objs interface{}
	//if err := json.Unmarshal(body, &objs); err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Print(objs)
	Request.GetAllNamespace()

}

