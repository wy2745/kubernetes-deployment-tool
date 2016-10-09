package main

import (
	"fmt"
	"./ab"
	//WorkLoadController "./workload"
	//"./request"
	//"./locust"
	//"./kubemark"
	//"bufio"
	//"os"
)

func main() {

	//Request.GenetatePodBody("default", "mysql:5.7", "mysql----test")
	//Request.CreatePod_test("default", "mysql:5.7", "mysql----test", "300m", "400m", "1Gi", "2Gi")
	//Request.GetPodsOfNamespace_Test("default")
	//Request.CreateService("mysql-service", "mysql----test", "default", 3302, Request.Caicloud)
	//Request.DeleteService("default", "mysql-service", Request.Caicloud)
	//Request.DeletePod("default", "test0", Request.Caicloud)
	//Request.GenerateServiceBody("mysql-service", "mysql----test", "default", 3302)
	//WorkLoadController.MissionRecorder()
	//WorkLoadController.Start()
	//WorkLoadController.StartPodVersion()
	//Request.GetAllService_Test("default")
	//fmt.Println(request.GetPodsOfNamespace("default", request.Test))
	//Request.GetAllService("default", Request.Caicloud)
	//Request.GetJobOfNamespace("default", Request.Caicloud)
	//Request.CreateJob("default", "docker.io/zilinglius/workload:short-general", "wsg", int32(5), "400m", "500m", "10Mi", "20Mi", "./home/wsg 200000000", Request.Caicloud)
	//Request.GetJobOfNamespace("default", Request.Caicloud)
	//Request.DeleteJob("default", "wsg", Request.Caicloud)
	//ab.Abtest()
	//locust.LocustStart()
	//locust.LocustClose()
	//locust.FileTest()
	//kubemark.GenerateReplicationcontrollerBody(10)
	//kubemark.BuildNode(int32(5))
	//scanner := bufio.NewScanner(os.Stdin)
	//scanner.Scan()
	//scanner.Text()
	//kubemark.DeleteNode()
	ab.Abtest()
	fmt.Println("成功退出")
	//Request.DeletePod("default", "mysql----test")
	//Request.GetAllReplicationcontrollers_Test()
	//Request.CreateReplicationController_test("default", "mysql:5.7", "mysql-rpc", "mysql-test", "mysql----test", 2)
	//Request.DeleteReplicationController("default", "mysql-rpc")


}

