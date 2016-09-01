package main

import (
	Request "./request"
)

func main() {

	//Request.GenetatePodBody("default", "mysql:5.7", "mysql----test")
	//Request.CreatePod_test("default", "mysql:5.7", "mysql----test")
	Request.CreateService_test("mysql-service", "mysql----test", "default", 3302)
	//Request.GenerateServiceBody("mysql-service", "mysql----test", "default", 3302)
	Request.GetAllService_Test("default")
}

