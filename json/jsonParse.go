package jsonParse

import (
	"encoding/json"
	"strings"
	"net/http"
	"fmt"
	"io/ioutil"
	"../interf"
	"reflect"
	"../type137"
	"../type124"
)

func JsonNewDecoder(body []byte) *json.Decoder {
	dec := json.NewDecoder(strings.NewReader(string(body)))
	return dec

}

func JsonUnmarsha(body []byte, v interface{}) {
	var err error
	if err = json.Unmarshal(body, &v); err != nil {
	}
}
func JsonUnmarsha2(body []byte, v interf.Podinface) {
	var err error

	if err = json.Unmarshal(body, &v); err != nil {
	}
}

func JsonUnmarsha2Interface(resp *http.Response, v *interf.Podinface) interf.Podinface {

	if (resp != nil) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
			return *v
		}
		fmt.Println(string(body))
		JsonUnmarsha2(body, *v)
		fmt.Println(*v)
		//classType.PrintPod(v)
		return *v
	}
	return *v
}

func JsonUnmarshaPod(resp *http.Response, v *interf.Podinface) *interf.Podinface {

	if (resp != nil) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
			return v
		}
		switch reflect.TypeOf(*v).String() {
		case "type137.Pod":
			var v1 type137.Pod
			JsonUnmarsha(body, &v1)
			*v = v1
		case "type124.Pod":
			var v2 type124.Pod
			JsonUnmarsha(body, &v2)
			*v = v2
		}
		return v
	}
	return v
}
func JsonUnmarshaPodList(resp *http.Response, v *interf.PodListinface) *interf.PodListinface {

	if (resp != nil) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
			return v
		}
		switch reflect.TypeOf(*v).String() {
		case "type137.PodList":
			var v1 type137.PodList
			JsonUnmarsha(body, &v1)
			*v = v1
		case "type124.PodList":
			var v2 type124.PodList
			JsonUnmarsha(body, &v2)
			*v = v2
		}
		return v
	}
	return v
}
func JsonUnmarshaService(resp *http.Response, v *interf.Podinface) *interf.Podinface {

	if (resp != nil) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
			return v
		}
		switch reflect.TypeOf(*v).String() {
		case "type137.Pod":
			var v1 type137.Pod
			JsonUnmarsha(body, &v1)
			*v = v1
		case "type124.Pod":
			var v2 type124.Pod
			JsonUnmarsha(body, &v2)
			*v = v2
		}
		return v
	}
	return v
}

func JsonMarsha(v interface{}) []byte {
	b, err := json.Marshal(&v)
	if err != nil {
	}
	return b
}


