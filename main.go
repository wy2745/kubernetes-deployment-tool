package main

import (
	"fmt"
	"github.com/wy2745/kubernetes-deployment-tool/kubemark"
	"github.com/wy2745/kubernetes-deployment-tool/autoscale"
	"strconv"
	"os"
	"bufio"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	mode := os.Args[1]
	switch mode {
	case "-k":
		nodeNum, _ := strconv.Atoi(os.Args[2])
		kubemark.PodListTestV2(nodeNum)
		fmt.Println("成功退出")
	case "-t":
		fmt.Println("ok")
	case "-l":
		autoscale.BuildNginx()
		scanner.Scan()
		scanner.Text()
		autoscale.DestoryNginx()
	default:
		fmt.Println("参数输入错误")
	}

}

