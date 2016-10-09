package kubemark

import (
	"../request"
	"fmt"
	"io/ioutil"
	"../json"
	classType "../type137"
	"time"
	"../ab"
	"strconv"
)

const (
	replicationControllerName = "nginx"
)

func podCreate(replic int32) chan int64 {
	var command []string
	startTime := time.Now()
	command = append(command, "/home/auto-reload-nginx.sh")
	body := request.GenerateReplicationcontrollerBodyV2("default", "ymqytw/nginxhttps:1.5", replicationControllerName, replic, command)
	url := request.KubemarkServer_Test + request.GenerateReplicationControllerNamespaceUrl("defalut")
	InvokeRequest("POST", url, body)
	timeChan := make(chan int64)
	go func() {
		for {
			url = request.KubemarkServer_Test + request.GeneratePodNamespaceUrl("default")
			resp := InvokeRequest("GET", url, nil)
			var count = 0
			if (resp != nil) {
				defer resp.Body.Close()
				var v classType.PodList
				body, err := ioutil.ReadAll(resp.Body)
				if (err != nil) {
					fmt.Print(err)
				}
				jsonParse.JsonUnmarsha(body, &v)
				for _, pod := range v.Items {
					if pod.Labels["name"] == replicationControllerName && pod.Status.Phase == "Running" {
						count ++
					}
				}
				if count == replic {
					endTime := time.Now()
					timeChan <- endTime.Second() - startTime.Second()
					return
				}
			}
		}
	}(startTime, replic)
	return timeChan
}

func podDelete() int64 {
	startTime := time.Now()
	url := request.KubemarkServer_Test + request.GenerateReplicationControllerNameUrl("default", replicationControllerName)
	InvokeRequest("DELETE", url, nil)

	url = request.KubemarkServer_Test + request.GeneratePodNamespaceUrl("default")
	resp := InvokeRequest("GET", url, nil)
	if (resp != nil) {
		defer resp.Body.Close()
		var v classType.PodList
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
		}
		jsonParse.JsonUnmarsha(body, &v)
		for _, pod := range v.Items {
			if pod.Labels["name"] == replicationControllerName {
				url = request.GeneratePodNameUrl("default", pod.Name)
				InvokeRequest("DELETE", url, nil)
			}
		}
	}
	endTime := time.Now()
	return endTime.Second() - startTime.Second()
}

func getNodeNum() int {
	url := request.KubemarkServer_Test + request.GenerateNodeUrl()
	resp := InvokeRequest("GET", url, nil)
	if (resp != nil) {
		defer resp.Body.Close()
		var v classType.NodeList
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
		}
		jsonParse.JsonUnmarsha(body, &v)
		return len(v.Items)
	}
	return 0
}

func PodListTest() {
	var rate []int
	rate = append(rate, 3)
	rate = append(rate, 5)
	rate = append(rate, 10)
	rate = append(rate, 15)
	rate = append(rate, 20)
	rate = append(rate, 30)
	var nodeNum = getNodeNum()
	for _, replic := range rate {
		timeChan := podCreate(int32(replic * nodeNum))
		createTime := <-timeChan
		fmt.Println("在", nodeNum, "个node上创建", replic * nodeNum, "个pod 使用了", createTime, "s")
		nodeN := strconv.Itoa(nodeNum)
		podN := strconv.Itoa(replic * nodeNum)
		ab.Abtest(nodeN + "n" + podN + "p")
		fmt.Println("在", nodeNum, "个node上删除", replic * nodeNum, "个pod 使用了", podDelete(), "s")

	}
}
