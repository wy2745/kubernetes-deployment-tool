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
	"bufio"
	"os"
)

const (
	replicationControllerName = "nginx"
)

func podCreate(replic int32) int64 {
	var command []string
	startTime := time.Now()
	command = append(command, "/home/auto-reload-nginx.sh")
	body := request.GenerateReplicationcontrollerBodyV2("default", "ymqytw/nginxhttps:1.5", replicationControllerName, replic, command)
	url := request.KubemarkServer_Test + request.GenerateReplicationControllerNamespaceUrl("default")
	InvokeRequest("POST", url, body)
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
			if count == int(replic) {
				endTime := time.Now()
				return endTime.Unix() - startTime.Unix()

			}
		}
	}
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
				url = request.KubemarkServer_Test + request.GeneratePodNameUrl("default", pod.Name)
				InvokeRequest("DELETE", url, nil)
			}
		}
	}
	endTime := time.Now()
	return endTime.Unix() - startTime.Unix()
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

func DeleteNodev2() {
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
		for _, node := range v.Items {
			url = request.KubemarkServer_Test + request.GenerateNodeNameUrl(node.Name)
			InvokeRequest("DELETE", url, nil)
		}
	}
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
	fmt.Println("node num：", nodeNum)
	fmt.Println("第1次测试")
	for _, replic := range rate {
		fmt.Println("在", nodeNum, "个node上创建", replic * nodeNum, "个pod 使用了", podCreate(int32(replic * nodeNum)), "s")
		nodeN := strconv.Itoa(nodeNum)
		podN := strconv.Itoa(replic * nodeNum)
		ab.Abtest(nodeN + "n" + podN + "p", "1")
		time.Sleep(time.Second * 10)
		fmt.Println("在", nodeNum, "个node上删除", replic * nodeNum, "个pod 使用了", podDelete(), "s")
	}
	time.Sleep(time.Second * 10)
	fmt.Println("第2次测试")
	for _, replic := range rate {
		fmt.Println("在", nodeNum, "个node上创建", replic * nodeNum, "个pod 使用了", podCreate(int32(replic * nodeNum)), "s")
		nodeN := strconv.Itoa(nodeNum)
		podN := strconv.Itoa(replic * nodeNum)
		ab.Abtest(nodeN + "n" + podN + "p", "2")
		time.Sleep(time.Second * 10)
		fmt.Println("在", nodeNum, "个node上删除", replic * nodeNum, "个pod 使用了", podDelete(), "s")
	}

	time.Sleep(time.Second * 10)
	fmt.Println("第3次测试")
	for _, replic := range rate {
		fmt.Println("在", nodeNum, "个node上创建", replic * nodeNum, "个pod 使用了", podCreate(int32(replic * nodeNum)), "s")
		nodeN := strconv.Itoa(nodeNum)
		podN := strconv.Itoa(replic * nodeNum)
		ab.Abtest(nodeN + "n" + podN + "p", "3")
		time.Sleep(time.Second * 10)
		fmt.Println("在", nodeNum, "个node上删除", replic * nodeNum, "个pod 使用了", podDelete(), "s")
	}

}

func Test() {
	var line string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("^_^")
	fmt.Println("1.删除node")
	fmt.Println("2.删除pod")
	fmt.Println("3.跑测试")
	fmt.Println("4.退出")
	for {
		scanner.Scan()
		line = scanner.Text()
		switch line {
		case "1":
			DeleteNodev2()
		case "2":
			podDelete()
		case "3":
			PodListTest()
		case "4":
			return
		}
		fmt.Println("1.删除node")
		fmt.Println("2.删除pod")
		fmt.Println("3.跑测试")
		fmt.Println("4.退出")

	}
}

