package main

import (
	"fmt"
	"github.com/wy2745/kubernetes-deployment-tool/kubemark"
	"github.com/wy2745/kubernetes-deployment-tool/autoscale"
	"github.com/wy2745/kubernetes-deployment-tool/locust"
	"strconv"
	"os"
	"net/http"
	"time"
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
		tr := http.Transport{DisableKeepAlives:false}
		client := http.Client{Transport:&tr}
		kubemark.PodCreate(podNum, &client)
	case "-cpt":
		nodeNum, _ := strconv.Atoi(os.Args[2])
		count, _ := strconv.Atoi(os.Args[3])
		kubemark.CptHandler(nodeNum, count)
	case "-dp":
		var clients []http.Client
		for i := 0; i < 4800; i++ {
			tr := http.Transport{DisableKeepAlives:false}
			client := http.Client{Transport:&tr}
			clients = append(clients, client)
		}
		kubemark.PodDelete(clients)
	case "-cn":
		nodeNum, _ := strconv.Atoi(os.Args[2])
		kubemark.CnHandler(nodeNum)
	case "-dn":
		var clients []http.Client
		for i := 0; i < 320; i++ {
			tr := http.Transport{DisableKeepAlives:false}
			client := http.Client{Transport:&tr}
			clients = append(clients, client)
		}
		kubemark.DeleteNodev2(clients)
	case "-k":
		//nodeNum, _ := strconv.Atoi(os.Args[2])
		//kubemark.PodListTestV2(nodeNum)
		fmt.Println("成功退出")
	case "-t":
		time.Sleep(time.Second * 60)
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

