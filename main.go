package main

import "fmt"
import "net/http"
import "io/ioutil"

func main() {
	resp,err := http.Get("http://202.120.40.177:16380/api/v1/pods")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Print(string(body))
}

