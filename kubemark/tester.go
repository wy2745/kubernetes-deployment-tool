package kubemark

import (
	"github.com/wy2745/kubernetes-deployment-tool/request"
	"fmt"
	"io/ioutil"
	"github.com/wy2745/kubernetes-deployment-tool/json"
	classType "github.com/wy2745/kubernetes-deployment-tool/type137"
	"time"
	"github.com/wy2745/kubernetes-deployment-tool/ab"
	"strconv"
	"bufio"
	"os"
	"encoding/csv"
)

const (
	replicationControllerName = "nginx"
	unit = int64(1000000)
)

func podCreate(replic int32) int {
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
				return int((endTime.UnixNano() - startTime.UnixNano()) / unit)

			}
		}
	}
}

func podDelete() int {
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
	return int((endTime.UnixNano() - startTime.UnixNano()) / unit)
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
	var rate = [6]int{3, 5, 10, 15, 20, 30}
	var nodeNum = getNodeNum()
	fmt.Println("node num：", nodeNum)
	fmt.Println("第1次测试")
	f, _ := os.Create("/home/administrator/test/" + strconv.Itoa(nodeNum) + "n.csv")
	defer f.Close()

	w := csv.NewWriter(f)
	w.Write([]string{" ", strconv.Itoa(nodeNum * rate[0]) + "C", strconv.Itoa(nodeNum * rate[0]) + "D", strconv.Itoa(nodeNum * rate[1]) + "C", strconv.Itoa(nodeNum * rate[1]) + "D", strconv.Itoa(nodeNum * rate[2]) + "C", strconv.Itoa(nodeNum * rate[2]) + "D", strconv.Itoa(nodeNum * rate[3]) + "C", strconv.Itoa(nodeNum * rate[3]) + "D", strconv.Itoa(nodeNum * rate[4]) + "C", strconv.Itoa(nodeNum * rate[4]) + "D", strconv.Itoa(nodeNum * rate[5]) + "C", strconv.Itoa(nodeNum * rate[5]) + "D"})

	var data [13]int
	data[0] = nodeNum
	for index, replic := range rate {
		data[2 * index + 1] = podCreate(int32(replic * nodeNum))
		fmt.Println("在", nodeNum, "个node上创建", replic * nodeNum, "个pod 使用了", data[2 * index + 1], "ms")
		//time.Sleep(time.Second * 3)
		nodeN := strconv.Itoa(nodeNum)
		podN := strconv.Itoa(replic * nodeNum)
		ab.Abtest(nodeN + "n" + podN + "p", "1")
		//time.Sleep(time.Second * 3)
		data[2 * index + 2] = podDelete()
		fmt.Println("在", nodeNum, "个node上删除", replic * nodeNum, "个pod 使用了", data[2 * index + 2], "ms")
	}
	w.Write([]string{strconv.Itoa(data[0]), strconv.Itoa(data[1]), strconv.Itoa(data[3]), strconv.Itoa(data[4]), strconv.Itoa(data[4]), strconv.Itoa(data[5]), strconv.Itoa(data[6]), strconv.Itoa(data[7]), strconv.Itoa(data[8]), strconv.Itoa(data[9]), strconv.Itoa(data[10]), strconv.Itoa(data[11]), strconv.Itoa(data[12])})
	time.Sleep(time.Second * 3)
	fmt.Println("第2次测试")
	for index, replic := range rate {
		data[2 * index + 1] = podCreate(int32(replic * nodeNum))
		fmt.Println("在", nodeNum, "个node上创建", replic * nodeNum, "个pod 使用了", data[2 * index + 1], "ms")
		//time.Sleep(time.Second * 3)
		nodeN := strconv.Itoa(nodeNum)
		podN := strconv.Itoa(replic * nodeNum)
		ab.Abtest(nodeN + "n" + podN + "p", "2")
		//time.Sleep(time.Second * 3)
		data[2 * index + 2] = podDelete()
		fmt.Println("在", nodeNum, "个node上删除", replic * nodeNum, "个pod 使用了", data[2 * index + 2], "ms")
	}
	w.Write([]string{strconv.Itoa(data[0]), strconv.Itoa(data[1]), strconv.Itoa(data[3]), strconv.Itoa(data[4]), strconv.Itoa(data[4]), strconv.Itoa(data[5]), strconv.Itoa(data[6]), strconv.Itoa(data[7]), strconv.Itoa(data[8]), strconv.Itoa(data[9]), strconv.Itoa(data[10]), strconv.Itoa(data[11]), strconv.Itoa(data[12])})
	time.Sleep(time.Second * 3)
	fmt.Println("第3次测试")
	for index, replic := range rate {
		data[2 * index + 1] = podCreate(int32(replic * nodeNum))
		fmt.Println("在", nodeNum, "个node上创建", replic * nodeNum, "个pod 使用了", data[2 * index + 1], "ms")
		//time.Sleep(time.Second * 3)
		nodeN := strconv.Itoa(nodeNum)
		podN := strconv.Itoa(replic * nodeNum)
		ab.Abtest(nodeN + "n" + podN + "p", "3")
		//time.Sleep(time.Second * 3)
		data[2 * index + 2] = podDelete()
		fmt.Println("在", nodeNum, "个node上删除", replic * nodeNum, "个pod 使用了", data[2 * index + 2], "ms")
	}
	w.Write([]string{strconv.Itoa(data[0]), strconv.Itoa(data[1]), strconv.Itoa(data[3]), strconv.Itoa(data[4]), strconv.Itoa(data[4]), strconv.Itoa(data[5]), strconv.Itoa(data[6]), strconv.Itoa(data[7]), strconv.Itoa(data[8]), strconv.Itoa(data[9]), strconv.Itoa(data[10]), strconv.Itoa(data[11]), strconv.Itoa(data[12]), })
	w.Flush()
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

