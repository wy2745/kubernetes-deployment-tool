package locust

import (
	"os/exec"
	"fmt"
	"bufio"
	"os"
	"strings"
)

const (
	destination string = "http://120.26.120.30:30888"
	fileroot string = "/Users/panda/Documents/github/locustfile.py"
	//本地启动的指令，以后可能会使用master-slave模式
	startcommand string = "locust -f " + fileroot + " --host=" + destination
)

func LocustStart() {
	cmd := exec.Command("/bin/sh", "-c", startcommand)
	//cmd := exec.Command("/bin/sh", "-c", "ls")
	//out, err := cmd.Output()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(string(out))
	cmd.Start()
	fmt.Println("成功开启locust")

	//err = cmd.Wait()
	//err := cmd.Run()
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	_ = scanner.Text()
	cmd = exec.Command("/bin/sh", "-c", "lsof -i:8089")
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