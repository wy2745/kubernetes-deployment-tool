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
)

const (
	destination string = "http://120.26.120.30:30888"
	fileroot string = "/Users/panda/Documents/github/locustfile.py"
	//本地启动的指令，以后可能会使用master-slave模式
	startcommand string = "locust -f " + fileroot + " --host=" + destination
	startTestUrl string = "http://localhost:8089/swarm"
	stopTestUrl string = "http://localhost:8089/stop"
	locust_count string = "locust_count"
	hatch_rate string = "hatch_rate"
	requestsUrl string = "http://localhost:8089/stats/requests/csv"
	distributionUrl string = "http://localhost:8089/stats/distribution/csv"
	exceptionsUrl string = "http://localhost:8089/exceptions/csv"
	requestRoot string = "/Users/panda/Documents/github/request.csv"
	distributionRoot string = "/Users/panda/Documents/github/idstribution.csv"
	exceptionsRoot string = "/Users/panda/Documents/github/exceptions.csv"
)

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
	FileTest()
}
func FileTest() {
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
			LocustTestStop()
			fmt.Println("测试完成，结果已经存起来了")
			return
		}
	}()
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