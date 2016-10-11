package kubemark

import (
	"github.com/wy2745/kubernetes-deployment-tool/request"
	"fmt"
	"io/ioutil"
	"github.com/wy2745/kubernetes-deployment-tool/json"
	classType "github.com/wy2745/kubernetes-deployment-tool/type137"
	"time"
	"strconv"
	"bufio"
	"os"
	"os/exec"
	"encoding/csv"
	"github.com/wy2745/kubernetes-deployment-tool/ab"
	"log"
)

const (
	replicationControllerName = "nginx"
	unit = int64(1000000)
)

func PodCreate(replic int32) int {
	var command []string
	startTime := time.Now()
	command = append(command, "/home/auto-reload-nginx.sh")
	body := request.GenerateReplicationcontrollerBodyV2("default", "ymqytw/nginxhttps:1.5", replicationControllerName, replic, command)
	url := request.KubemarkServer_Test + request.GenerateReplicationControllerNamespaceUrl("default")
	InvokeRequest("POST", url, body)
	for {
		url = request.KubemarkServer_Test + request.GeneratePodNamespaceUrl("default")
		resp := InvokeRequest("GET", url, nil)
		var count = int32(0)
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
					for _, pc := range pod.Status.Conditions {
						if pc.Type == "Ready" && pc.Status == "True" {
							count ++
						}
					}

				}
			}
			if count == replic {
				endTime := time.Now()
				return int((endTime.UnixNano() - startTime.UnixNano()) / unit)

			}
		}
	}
}

func PodDelete() int {
	fmt.Println("开始删除")
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

	for {
		url = request.KubemarkServer_Test + request.GeneratePodNamespaceUrl("default")
		resp = InvokeRequest("GET", url, nil)
		if (resp != nil) {
			defer resp.Body.Close()
			var v classType.PodList
			body, err := ioutil.ReadAll(resp.Body)
			if (err != nil) {
				fmt.Print(err)
			}
			jsonParse.JsonUnmarsha(body, &v)
			if len(v.Items) == 0 {
				break
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

func waitallNodeReady(nodeNum int) bool {
	for {
		if getNodeNum() == nodeNum {
			return true
		}
	}
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

	for {
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
			if len(v.Items) == 0 {
				return
			}
		}
	}
}

func PodListTest() {

	PodDelete()

	var nodenum = [7]int{5, 10, 20, 40, 80, 160, 320}
	for i := 0; i < 7; i++ {

		nodeNum := nodenum[i]
		ChangeNode(nodeNum)

		waitallNodeReady(nodeNum)
		time.Sleep(time.Second * 5)

		var rate = [6]int{3, 5, 10, 15, 20, 30}
		fmt.Println("node num：", nodeNum)
		fmt.Println("第1次测试")
		f, _ := os.Create("/home/administrator/test/" + strconv.Itoa(nodeNum) + "n.csv")
		defer f.Close()

		w := csv.NewWriter(f)
		w.Write([]string{" ", strconv.Itoa(nodeNum * rate[0]) + "C", strconv.Itoa(nodeNum * rate[0]) + "D", strconv.Itoa(nodeNum * rate[1]) + "C", strconv.Itoa(nodeNum * rate[1]) + "D", strconv.Itoa(nodeNum * rate[2]) + "C", strconv.Itoa(nodeNum * rate[2]) + "D", strconv.Itoa(nodeNum * rate[3]) + "C", strconv.Itoa(nodeNum * rate[3]) + "D", strconv.Itoa(nodeNum * rate[4]) + "C", strconv.Itoa(nodeNum * rate[4]) + "D", strconv.Itoa(nodeNum * rate[5]) + "C", strconv.Itoa(nodeNum * rate[5]) + "D"})

		var data [13]int
		var data1 [13]int
		var data2 [13]int
		data[0] = nodeNum
		data1[0] = nodeNum
		data2[0] = nodeNum
		for index, replic := range rate {
			data[2 * index + 1] = PodCreate(int32(replic * nodeNum))
			fmt.Println("在", nodeNum, "个node上创建", replic * nodeNum, "个pod 使用了", data[2 * index + 1], "ms")
			nodeN := strconv.Itoa(nodeNum)
			podN := strconv.Itoa(replic * nodeNum)
			//ab.Abtest(nodeN + "n" + podN + "p", "1")
			ab.AbtestV2(nodeN + "n" + podN + "p", "1")
			time.Sleep(time.Second * 3)
			data[2 * index + 2] = PodDelete()
			fmt.Println("在", nodeNum, "个node上删除", replic * nodeNum, "个pod 使用了", data[2 * index + 2], "ms")
			data1[2 * index + 1] = PodCreate(int32(replic * nodeNum))
			fmt.Println("在", nodeNum, "个node上创建", replic * nodeNum, "个pod 使用了", data1[2 * index + 1], "ms")
			//ab.Abtest(nodeN + "n" + podN + "p", "1")
			ab.AbtestV2(nodeN + "n" + podN + "p", "1")
			time.Sleep(time.Second * 3)
			data1[2 * index + 2] = PodDelete()
			fmt.Println("在", nodeNum, "个node上删除", replic * nodeNum, "个pod 使用了", data1[2 * index + 2], "ms")
			data2[2 * index + 1] = PodCreate(int32(replic * nodeNum))
			fmt.Println("在", nodeNum, "个node上创建", replic * nodeNum, "个pod 使用了", data2[2 * index + 1], "ms")
			//ab.Abtest(nodeN + "n" + podN + "p", "1")
			ab.AbtestV2(nodeN + "n" + podN + "p", "1")
			time.Sleep(time.Second * 10)
			data2[2 * index + 2] = PodDelete()
			fmt.Println("在", nodeNum, "个node上删除", replic * nodeNum, "个pod 使用了", data2[2 * index + 2], "ms")
		}
		w.Write([]string{strconv.Itoa(data[0]), strconv.Itoa(data[1]), strconv.Itoa(data[2]), strconv.Itoa(data[3]), strconv.Itoa(data[4]), strconv.Itoa(data[5]), strconv.Itoa(data[6]), strconv.Itoa(data[7]), strconv.Itoa(data[8]), strconv.Itoa(data[9]), strconv.Itoa(data[10]), strconv.Itoa(data[11]), strconv.Itoa(data[12])})
		w.Write([]string{strconv.Itoa(data1[0]), strconv.Itoa(data1[1]), strconv.Itoa(data1[2]), strconv.Itoa(data1[3]), strconv.Itoa(data1[4]), strconv.Itoa(data1[5]), strconv.Itoa(data1[6]), strconv.Itoa(data1[7]), strconv.Itoa(data1[8]), strconv.Itoa(data1[9]), strconv.Itoa(data1[10]), strconv.Itoa(data1[11]), strconv.Itoa(data1[12])})
		w.Write([]string{strconv.Itoa(data2[0]), strconv.Itoa(data2[1]), strconv.Itoa(data2[2]), strconv.Itoa(data2[3]), strconv.Itoa(data2[4]), strconv.Itoa(data2[5]), strconv.Itoa(data2[6]), strconv.Itoa(data2[7]), strconv.Itoa(data2[8]), strconv.Itoa(data2[9]), strconv.Itoa(data2[10]), strconv.Itoa(data2[11]), strconv.Itoa(data2[12])})
		//time.Sleep(time.Second * 3)
		//fmt.Println("第2次测试")
		//for index, replic := range rate {
		//	data[2 * index + 1] = podCreate(int32(replic * nodeNum))
		//	fmt.Println("在", nodeNum, "个node上创建", replic * nodeNum, "个pod 使用了", data[2 * index + 1], "ms")
		//	time.Sleep(time.Second * 3)
		//	nodeN := strconv.Itoa(nodeNum)
		//	podN := strconv.Itoa(replic * nodeNum)
		//	ab.Abtest(nodeN + "n" + podN + "p", "2")
		//	time.Sleep(time.Second * 3)
		//	data[2 * index + 2] = podDelete()
		//	fmt.Println("在", nodeNum, "个node上删除", replic * nodeNum, "个pod 使用了", data[2 * index + 2], "ms")
		//}
		//w.Write([]string{strconv.Itoa(data[0]), strconv.Itoa(data[1]), strconv.Itoa(data[2]), strconv.Itoa(data[3]), strconv.Itoa(data[4]), strconv.Itoa(data[5]), strconv.Itoa(data[6]), strconv.Itoa(data[7]), strconv.Itoa(data[8]), strconv.Itoa(data[9]), strconv.Itoa(data[10]), strconv.Itoa(data[11]), strconv.Itoa(data[12])})
		//time.Sleep(time.Second * 3)
		//fmt.Println("第3次测试")
		//for index, replic := range rate {
		//	data[2 * index + 1] = podCreate(int32(replic * nodeNum))
		//	fmt.Println("在", nodeNum, "个node上创建", replic * nodeNum, "个pod 使用了", data[2 * index + 1], "ms")
		//	time.Sleep(time.Second * 3)
		//	nodeN := strconv.Itoa(nodeNum)
		//	podN := strconv.Itoa(replic * nodeNum)
		//	ab.Abtest(nodeN + "n" + podN + "p", "3")
		//	time.Sleep(time.Second * 3)
		//	data[2 * index + 2] = podDelete()
		//	fmt.Println("在", nodeNum, "个node上删除", replic * nodeNum, "个pod 使用了", data[2 * index + 2], "ms")
		//}
		//w.Write([]string{strconv.Itoa(data[0]), strconv.Itoa(data[1]), strconv.Itoa(data[2]), strconv.Itoa(data[3]), strconv.Itoa(data[4]), strconv.Itoa(data[5]), strconv.Itoa(data[6]), strconv.Itoa(data[7]), strconv.Itoa(data[8]), strconv.Itoa(data[9]), strconv.Itoa(data[10]), strconv.Itoa(data[11]), strconv.Itoa(data[12]), })
		w.Flush()
		if err := w.Error(); err != nil {
			log.Fatal(err)
		}
	}

}

func PodListTestV2(nodeNum int) {

	PodDelete()

	ChangeNode(nodeNum)

	waitallNodeReady(nodeNum)
	time.Sleep(time.Second * 5)

	var rate = [6]int{3, 5, 10, 15, 20, 30}
	fmt.Println("node num：", nodeNum)
	fmt.Println("第1次测试")
	f, _ := os.Create("/home/administrator/test/" + strconv.Itoa(nodeNum) + "n.csv")
	defer f.Close()

	w := csv.NewWriter(f)
	w.Write([]string{" ", strconv.Itoa(nodeNum * rate[0]) + "C", strconv.Itoa(nodeNum * rate[0]) + "D", strconv.Itoa(nodeNum * rate[1]) + "C", strconv.Itoa(nodeNum * rate[1]) + "D", strconv.Itoa(nodeNum * rate[2]) + "C", strconv.Itoa(nodeNum * rate[2]) + "D", strconv.Itoa(nodeNum * rate[3]) + "C", strconv.Itoa(nodeNum * rate[3]) + "D", strconv.Itoa(nodeNum * rate[4]) + "C", strconv.Itoa(nodeNum * rate[4]) + "D", strconv.Itoa(nodeNum * rate[5]) + "C", strconv.Itoa(nodeNum * rate[5]) + "D"})

	var data [13]int
	var data1 [13]int
	var data2 [13]int
	data[0] = nodeNum
	data1[0] = nodeNum
	data2[0] = nodeNum
	for index, replic := range rate {
		data[2 * index + 1] = PodCreate(int32(replic * nodeNum))
		fmt.Println("在", nodeNum, "个node上创建", replic * nodeNum, "个pod 使用了", data[2 * index + 1], "ms")
		nodeN := strconv.Itoa(nodeNum)
		podN := strconv.Itoa(replic * nodeNum)
		ab.Abtest(nodeN + "n" + podN + "p", "1")
		time.Sleep(time.Second * 3)
		data[2 * index + 2] = PodDelete()
		fmt.Println("在", nodeNum, "个node上删除", replic * nodeNum, "个pod 使用了", data[2 * index + 2], "ms")
		data1[2 * index + 1] = PodCreate(int32(replic * nodeNum))
		fmt.Println("在", nodeNum, "个node上创建", replic * nodeNum, "个pod 使用了", data1[2 * index + 1], "ms")
		ab.Abtest(nodeN + "n" + podN + "p", "1")
		time.Sleep(time.Second * 3)
		data1[2 * index + 2] = PodDelete()
		fmt.Println("在", nodeNum, "个node上删除", replic * nodeNum, "个pod 使用了", data1[2 * index + 2], "ms")
		data2[2 * index + 1] = PodCreate(int32(replic * nodeNum))
		fmt.Println("在", nodeNum, "个node上创建", replic * nodeNum, "个pod 使用了", data2[2 * index + 1], "ms")
		ab.Abtest(nodeN + "n" + podN + "p", "1")
		time.Sleep(time.Second * 10)
		data2[2 * index + 2] = PodDelete()
		fmt.Println("在", nodeNum, "个node上删除", replic * nodeNum, "个pod 使用了", data2[2 * index + 2], "ms")
	}
	w.Write([]string{strconv.Itoa(data[0]), strconv.Itoa(data[1]), strconv.Itoa(data[2]), strconv.Itoa(data[3]), strconv.Itoa(data[4]), strconv.Itoa(data[5]), strconv.Itoa(data[6]), strconv.Itoa(data[7]), strconv.Itoa(data[8]), strconv.Itoa(data[9]), strconv.Itoa(data[10]), strconv.Itoa(data[11]), strconv.Itoa(data[12])})
	w.Write([]string{strconv.Itoa(data1[0]), strconv.Itoa(data1[1]), strconv.Itoa(data1[2]), strconv.Itoa(data1[3]), strconv.Itoa(data1[4]), strconv.Itoa(data1[5]), strconv.Itoa(data1[6]), strconv.Itoa(data1[7]), strconv.Itoa(data1[8]), strconv.Itoa(data1[9]), strconv.Itoa(data1[10]), strconv.Itoa(data1[11]), strconv.Itoa(data1[12])})
	w.Write([]string{strconv.Itoa(data2[0]), strconv.Itoa(data2[1]), strconv.Itoa(data2[2]), strconv.Itoa(data2[3]), strconv.Itoa(data2[4]), strconv.Itoa(data2[5]), strconv.Itoa(data2[6]), strconv.Itoa(data2[7]), strconv.Itoa(data2[8]), strconv.Itoa(data2[9]), strconv.Itoa(data2[10]), strconv.Itoa(data2[11]), strconv.Itoa(data2[12])})
	//time.Sleep(time.Second * 3)
	//fmt.Println("第2次测试")
	//for index, replic := range rate {
	//	data[2 * index + 1] = podCreate(int32(replic * nodeNum))
	//	fmt.Println("在", nodeNum, "个node上创建", replic * nodeNum, "个pod 使用了", data[2 * index + 1], "ms")
	//	time.Sleep(time.Second * 3)
	//	nodeN := strconv.Itoa(nodeNum)
	//	podN := strconv.Itoa(replic * nodeNum)
	//	ab.Abtest(nodeN + "n" + podN + "p", "2")
	//	time.Sleep(time.Second * 3)
	//	data[2 * index + 2] = podDelete()
	//	fmt.Println("在", nodeNum, "个node上删除", replic * nodeNum, "个pod 使用了", data[2 * index + 2], "ms")
	//}
	//w.Write([]string{strconv.Itoa(data[0]), strconv.Itoa(data[1]), strconv.Itoa(data[2]), strconv.Itoa(data[3]), strconv.Itoa(data[4]), strconv.Itoa(data[5]), strconv.Itoa(data[6]), strconv.Itoa(data[7]), strconv.Itoa(data[8]), strconv.Itoa(data[9]), strconv.Itoa(data[10]), strconv.Itoa(data[11]), strconv.Itoa(data[12])})
	//time.Sleep(time.Second * 3)
	//fmt.Println("第3次测试")
	//for index, replic := range rate {
	//	data[2 * index + 1] = podCreate(int32(replic * nodeNum))
	//	fmt.Println("在", nodeNum, "个node上创建", replic * nodeNum, "个pod 使用了", data[2 * index + 1], "ms")
	//	time.Sleep(time.Second * 3)
	//	nodeN := strconv.Itoa(nodeNum)
	//	podN := strconv.Itoa(replic * nodeNum)
	//	ab.Abtest(nodeN + "n" + podN + "p", "3")
	//	time.Sleep(time.Second * 3)
	//	data[2 * index + 2] = podDelete()
	//	fmt.Println("在", nodeNum, "个node上删除", replic * nodeNum, "个pod 使用了", data[2 * index + 2], "ms")
	//}
	//w.Write([]string{strconv.Itoa(data[0]), strconv.Itoa(data[1]), strconv.Itoa(data[2]), strconv.Itoa(data[3]), strconv.Itoa(data[4]), strconv.Itoa(data[5]), strconv.Itoa(data[6]), strconv.Itoa(data[7]), strconv.Itoa(data[8]), strconv.Itoa(data[9]), strconv.Itoa(data[10]), strconv.Itoa(data[11]), strconv.Itoa(data[12]), })
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

}
func CnHandler(nodeNum int) {
	ChangeNode(nodeNum)

	waitallNodeReady(nodeNum)
	time.Sleep(time.Second * 3)
}
func AbHandler(nodeNum int, replic int, count int) {
	nodeN := strconv.Itoa(nodeNum)
	podN := strconv.Itoa(replic * nodeNum)
	couN := strconv.Itoa(count)
	ab.Abtest(nodeN + "n" + podN + "p", couN)
}

func CptHandler(nodeNum int) {
	var rate = [6]int{3, 5, 10, 15, 20, 30}
	fmt.Println("node num：", nodeNum)
	fmt.Println("开始测试")
	f, _ := os.Create("/home/administrator/test/" + strconv.Itoa(nodeNum) + "n.csv")
	defer f.Close()

	w := csv.NewWriter(f)
	w.Write([]string{" ", strconv.Itoa(nodeNum * rate[0]) + "C", strconv.Itoa(nodeNum * rate[0]) + "D", strconv.Itoa(nodeNum * rate[1]) + "C", strconv.Itoa(nodeNum * rate[1]) + "D", strconv.Itoa(nodeNum * rate[2]) + "C", strconv.Itoa(nodeNum * rate[2]) + "D", strconv.Itoa(nodeNum * rate[3]) + "C", strconv.Itoa(nodeNum * rate[3]) + "D", strconv.Itoa(nodeNum * rate[4]) + "C", strconv.Itoa(nodeNum * rate[4]) + "D", strconv.Itoa(nodeNum * rate[5]) + "C", strconv.Itoa(nodeNum * rate[5]) + "D"})

	var data [13]int
	var data1 [13]int
	var data2 [13]int
	data[0] = nodeNum
	data1[0] = nodeNum
	data2[0] = nodeNum
	for index, replic := range rate {
		data[2 * index + 1] = PodCreate(int32(replic * nodeNum))
		fmt.Println("在", nodeNum, "个node上创建", replic * nodeNum, "个pod 使用了", data[2 * index + 1], "ms")
		data[2 * index + 2] = PodDelete()
		fmt.Println("在", nodeNum, "个node上删除", replic * nodeNum, "个pod 使用了", data[2 * index + 2], "ms")
		data1[2 * index + 1] = PodCreate(int32(replic * nodeNum))
		fmt.Println("在", nodeNum, "个node上创建", replic * nodeNum, "个pod 使用了", data1[2 * index + 1], "ms")
		data1[2 * index + 2] = PodDelete()
		fmt.Println("在", nodeNum, "个node上删除", replic * nodeNum, "个pod 使用了", data1[2 * index + 2], "ms")
		data2[2 * index + 1] = PodCreate(int32(replic * nodeNum))
		fmt.Println("在", nodeNum, "个node上创建", replic * nodeNum, "个pod 使用了", data2[2 * index + 1], "ms")
		data2[2 * index + 2] = PodDelete()
		fmt.Println("在", nodeNum, "个node上删除", replic * nodeNum, "个pod 使用了", data2[2 * index + 2], "ms")
	}
	w.Write([]string{strconv.Itoa(data[0]), strconv.Itoa(data[1]), strconv.Itoa(data[2]), strconv.Itoa(data[3]), strconv.Itoa(data[4]), strconv.Itoa(data[5]), strconv.Itoa(data[6]), strconv.Itoa(data[7]), strconv.Itoa(data[8]), strconv.Itoa(data[9]), strconv.Itoa(data[10]), strconv.Itoa(data[11]), strconv.Itoa(data[12])})
	w.Write([]string{strconv.Itoa(data1[0]), strconv.Itoa(data1[1]), strconv.Itoa(data1[2]), strconv.Itoa(data1[3]), strconv.Itoa(data1[4]), strconv.Itoa(data1[5]), strconv.Itoa(data1[6]), strconv.Itoa(data1[7]), strconv.Itoa(data1[8]), strconv.Itoa(data1[9]), strconv.Itoa(data1[10]), strconv.Itoa(data1[11]), strconv.Itoa(data1[12])})
	w.Write([]string{strconv.Itoa(data2[0]), strconv.Itoa(data2[1]), strconv.Itoa(data2[2]), strconv.Itoa(data2[3]), strconv.Itoa(data2[4]), strconv.Itoa(data2[5]), strconv.Itoa(data2[6]), strconv.Itoa(data2[7]), strconv.Itoa(data2[8]), strconv.Itoa(data2[9]), strconv.Itoa(data2[10]), strconv.Itoa(data2[11]), strconv.Itoa(data2[12])})
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

func ChangeNode(num int) {
	DeleteNodev2()
	cmd := exec.Command("/bin/sh", "-c", "~/go/src/github.com/wy2745/kubernetes-deployment-tool/changeNode.sh " + strconv.Itoa(num))
	cmd.Output()
}

//func AbscriptTest() {
//	cmd := exec.Command("/bin/sh", "-c", "~/go/src/github.com/wy2745/kubernetes-deployment-tool/abTest.sh /home/administrator/test/ab/abc.csv /home/administrator/test/ab/abc.gnp http://www.baidu.com/")
//	cmd.Output()
//}

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
			PodDelete()
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
//func Csvfunc() {
//	f, _ := os.Create("/Users/panda/Desktop/" + "5n.csv")
//	var rate = [6]int{3, 5, 10, 15, 20, 30}
//	defer f.Close()
//
//	w := csv.NewWriter(f)
//	w.Write([]string{strconv.Itoa(rate[0]) + "C", strconv.Itoa(rate[0]) + "C"})
//	w.Flush()
//	if err := w.Error(); err != nil {
//		log.Fatal(err)
//	}
//}

