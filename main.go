package main

import (
	"fmt"
	"github.com/wy2745/kubernetes-deployment-tool/kubemark"
	"github.com/wy2745/kubernetes-deployment-tool/autoscale"
	"github.com/wy2745/kubernetes-deployment-tool/locust"
	"strconv"
	"os"
)

func main() {

	mode := os.Args[1]
	switch mode {
	case "-ab":
		nodeNum, _ := strconv.Atoi(os.Args[2])
		replic, _ := strconv.Atoi(os.Args[3])
		count, _ := strconv.Atoi(os.Args[4])
		kubemark.AbHandler(nodeNum, replic, count)
	case "-cp":
		nodeNum, _ := strconv.Atoi(os.Args[2])
		replic, _ := strconv.Atoi(os.Args[3])
		podNum := int32(nodeNum * replic)
		kubemark.PodCreate(podNum)
	case "-cpt":
		nodeNum, _ := strconv.Atoi(os.Args[2])
		kubemark.CptHandler(nodeNum)
	case "-dp":
		kubemark.PodDelete()
	case "-cn":
		nodeNum, _ := strconv.Atoi(os.Args[2])
		kubemark.CnHandler(nodeNum)
	case "-dn":
		kubemark.DeleteNodev2()
	case "-k":
		nodeNum, _ := strconv.Atoi(os.Args[2])
		kubemark.PodListTestV2(nodeNum)
		fmt.Println("成功退出")
	case "-t":
		fmt.Println("ok")
	case "-l":
		replic, _ := strconv.Atoi(os.Args[2])
		autoscale.BuildNginx(int32(replic))
		locust.LocustTest(os.Args[3], os.Args[4])

		autoscale.DestoryNginx()
	default:
		fmt.Println("参数输入错误")
	}

}

