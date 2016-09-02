package main

import (
	Request "./request"
)

func main() {

	//Request.GenetatePodBody("default", "mysql:5.7", "mysql----test")
	//Request.CreatePod_test("default", "mysql:5.7", "mysql----test", "300m", "400m", "1Gi", "2Gi")
	//Request.GetPodsOfNamespace_Test("default")
	//Request.CreateService_test("mysql-service", "mysql----test", "default", 3302)
	//Request.DeleteService("default", "mysql-service")
	//Request.DeletePod("default", "mysql----test")
	//Request.GenerateServiceBody("mysql-service", "mysql----test", "default", 3302)
	Request.GetAllService_Test("default")
	//Request.GetAllReplicationcontrollers_Test()
	//Request.CreateReplicationController_test("default", "mysql:5.7", "mysql-rpc", "mysql-test", "mysql----test", 2)
	//Request.DeleteReplicationController("default", "mysql-rpc")
}

