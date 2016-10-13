package locust

import (
	"os/exec"
	"fmt"
	"strings"
	"bufio"
	"bytes"
	"net/http"
	"os"
	"io"
	"time"
	"strconv"
	"github.com/wy2745/kubernetes-deployment-tool/kubemark"
	"github.com/wy2745/kubernetes-deployment-tool/json"
	classType "github.com/wy2745/kubernetes-deployment-tool/type137"
	"io/ioutil"
	"encoding/csv"
	"log"
)

const (
	destination string = "http://120.26.120.30:30888"
	nodePortUrl string = "192.168.6.22:30080"
	apiProxyUrl string = "http://192.168.6.10:8080/api/v1/proxy/namespaces/default/services/nginx-svc:80"
	apiserviceUrl string = "http://192.168.6.10:8080"

	fileroot string = "/Users/panda/Documents/github/locustfile.py"
	//本地启动的指令，以后可能会使用master-slave模式
	startcommand string = "locust -f " + fileroot + " --host=" + destination
	destinationUrl string = "http://192.168.6.31:8089"
	startTestUrl string = destinationUrl + "/swarm"
	stopTestUrl string = destinationUrl + "/stop"
	locust_count string = "locust_count"
	hatch_rate string = "hatch_rate"
	requestsUrl string = destinationUrl + "/stats/requests/csv"
	requestStatUrl string = destinationUrl + "/stats/requests"
	distributionUrl string = destinationUrl + "/stats/distribution/csv"
	exceptionsUrl string = destinationUrl + "/exceptions/csv"
	destinationRoot string = "/home/administrator/test/locust/"
	requestRoot string = "/Users/panda/Documents/github/request.csv"
	distributionRoot string = "/Users/panda/Documents/github/distribution.csv"
	exceptionsRoot string = "/Users/panda/Documents/github/exceptions.csv"
)

func generateRequestName(fileName string) string {
	return destinationRoot + "request" + fileName + ".csv"
}
func generateDistributionName(fileName string) string {
	return destinationRoot + "distribution" + fileName + ".csv"
}
func generateExceptionsName(fileName string) string {
	return destinationRoot + "exceptions" + fileName + ".csv"
}
func getReplic() int {
	url := apiserviceUrl + kubemark.GenerateReplicationControllerNameUrl("default", "nginx")
	tr := http.Transport{DisableKeepAlives:false}
	client := http.Client{Transport:&tr}
	resp := kubemark.InvokeRequestV2("GET", url, nil, &client)
	if (resp != nil) {
		defer resp.Body.Close()
		var v classType.ReplicationController
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
		}
		jsonParse.JsonUnmarsha(body, &v)
		return int(v.Status.Replicas)
	}
	return 0
}

func Locust(scanner *bufio.Scanner) {
	fmt.Println("1.开启Locust")
	fmt.Println("2.关闭Locust")
	fmt.Println("3.设置压力参数并开始测试")
	fmt.Println("4.停止测试")
	fmt.Println("5.返回上一页")
	var line string
	for {
		scanner.Scan()
		line = scanner.Text()
		switch line {
		case "1":
			LocustStart()
		case "2":
			LocustClose()
		case "3":
			LocustParaSet(scanner)
		case "4":
			LocustTestStop()

		case "5":
			return
		}
		fmt.Println("1.开启Locust")
		fmt.Println("2.关闭Locust")
		fmt.Println("3.设置压力参数并开始测试")
		fmt.Println("4.停止测试")
		fmt.Println("5.返回上一页")
	}
}
func locustOpened() bool {
	cmd := exec.Command("/bin/sh", "-c", "lsof -i:8089")
	out, _ := cmd.Output()
	if len(out) == 0 || !strings.Contains(string(out), "Python") {
		return false
	}
	return true
}

func LocustStart() {

	if (!locustOpened()) {
		cmd := exec.Command("/bin/sh", "-c", startcommand)
		cmd.Start()
		fmt.Println("成功开启locust")
		return
	}
	fmt.Println("端口8089被占用或Locust已开启")
}
func swarmBody(str1 string, str2 string) string {
	return locust_count + "=" + str1 + "&" + hatch_rate + "=" + str2
}
func LocustTestStop() {
	http.Get(stopTestUrl)
}
func fileTest() {
	resp, err := http.Get(requestsUrl)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	out, err := os.Create(requestRoot)
	if err != nil {
		// panic?
	}
	defer out.Close()
	io.Copy(out, resp.Body)

	resp, err = http.Get(distributionUrl)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	out, err = os.Create(distributionRoot)
	if err != nil {
		// panic?
	}
	defer out.Close()
	io.Copy(out, resp.Body)

	resp, err = http.Get(exceptionsUrl)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	out, err = os.Create(exceptionsRoot)
	if err != nil {
		// panic?
	}
	defer out.Close()
	io.Copy(out, resp.Body)
}
func fileTestV2(fileName string) {
	resp, err := http.Get(requestsUrl)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	out, err := os.Create(generateRequestName(fileName))
	if err != nil {
		// panic?
	}
	defer out.Close()
	io.Copy(out, resp.Body)

	resp, err = http.Get(distributionUrl)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	out, err = os.Create(generateDistributionName(fileName))
	if err != nil {
		// panic?
	}
	defer out.Close()
	io.Copy(out, resp.Body)

	resp, err = http.Get(exceptionsUrl)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	out, err = os.Create(generateExceptionsName(fileName))
	if err != nil {
		// panic?
	}
	defer out.Close()
	io.Copy(out, resp.Body)
}

func LocustParaSet(scanner *bufio.Scanner) {
	var line string
	fmt.Println("输入模拟用户数")
	scanner.Scan()
	line = scanner.Text()
	locust_count := line
	fmt.Println("输入每秒模拟用户数量")
	scanner.Scan()
	line = scanner.Text()
	hatch_rate := line
	fmt.Println("输入测试时间，单位为s，如果为0，则测试一直运行")
	scanner.Scan()
	line = scanner.Text()
	duration := line

	client := http.Client{}
	fmt.Println(swarmBody(locust_count, hatch_rate))
	body := []byte(swarmBody(locust_count, hatch_rate))
	req, _ := http.NewRequest("POST", startTestUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client.Do(req)

	go func() {
		tim, _ := strconv.Atoi(duration)
		if tim == 0 {
			return
		} else {
			time.Sleep(time.Duration(tim) * time.Second)
			//LocustTestStop()
			fmt.Println("测试完成，结果已经存起来了")
			return
		}
	}()
}

func getUserNum() int {
	tr := http.Transport{DisableKeepAlives:false}
	client := http.Client{Transport:&tr}
	req, _ := http.NewRequest("GET", requestStatUrl, nil)
	resp, _ := client.Do(req)
	if (resp != nil) {
		defer resp.Body.Close()
		var v requestStat
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
		}
		jsonParse.JsonUnmarsha(body, &v)
		return v.User_count
	}
	return 0
}

func LocustTest(locust_count string, hatch_rate string) {
	LocustTestStop()
	time.Sleep(time.Second * 60)
	fmt.Println("准备启动测试")

	startReplic := getReplic()

	tr := http.Transport{DisableKeepAlives:false}
	client := http.Client{Transport:&tr}
	fmt.Println(swarmBody(locust_count, hatch_rate))
	body := []byte(swarmBody(locust_count, hatch_rate))
	req, _ := http.NewRequest("POST", startTestUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client.Do(req)
	fmt.Println("测试启动,参数:", locust_count + "个user" + hatch_rate + "个并发")
	usercount, _ := strconv.Atoi(locust_count)
	for {
		if getUserNum() == usercount {
			break
		}
		time.Sleep(time.Second * 2)
	}
	fmt.Println("测试将持续120s...")
	time.Sleep(time.Second * 30)
	fmt.Println("还有90s...")
	time.Sleep(time.Second * 30)
	fmt.Println("还有60s...")
	time.Sleep(time.Second * 30)
	fmt.Println("还有30s...")
	time.Sleep(time.Second * 30)

	fmt.Println("本次测试结束，准备收集并存储数据...")

	endReplic := getReplic()
	LocustTestStop()

	f, _ := os.Create("/home/administrator/test/locust/replic/" + locust_count + "C" + hatch_rate + "H.csv")
	defer f.Close()

	w := csv.NewWriter(f)
	w.Write([]string{"replic-start", "replic-end"})
	w.Write([]string{strconv.Itoa(startReplic), strconv.Itoa(endReplic)})

	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	fileTestV2(locust_count + "C" + hatch_rate + "H")
	fmt.Println("数据存储成功!")
}
func LocustClose() {
	if (!locustOpened()) {
		fmt.Println("Locust 没有开启")
		return
	}
	cmd := exec.Command("/bin/sh", "-c", "lsof -i:8089")
	out, _ := cmd.Output()
	str := strings.Split(string(out), "Python")
	for ind, ss := range str {
		if (ss[0] == byte('C')) {
			sss := strings.Split(str[ind + 1], " ")
			for _, s := range sss {
				if (s != "") {
					cmd = exec.Command("/bin/sh", "-c", "kill " + s)
					out, _ = cmd.Output()
					fmt.Println("成功关闭")
					break
				}
			}
		}

	}
}