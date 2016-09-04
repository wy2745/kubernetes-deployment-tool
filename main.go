package main

import (
	"fmt"
	WorkLoadController "./workload"
	"net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		http.Redirect(w, req, "http://202.120.40.175:21101/users/all", http.StatusTemporaryRedirect)
	} else {
		http.Redirect(w, req, "http://202.120.40.175:21101/users", http.StatusTemporaryRedirect)
	}
}

func main() {
	http.HandleFunc("/hello", HelloServer)
	http.ListenAndServe(":12345", nil)


	//Request.GenetatePodBody("default", "mysql:5.7", "mysql----test")
	//Request.CreatePod_test("default", "mysql:5.7", "mysql----test", "300m", "400m", "1Gi", "2Gi")
	//Request.GetPodsOfNamespace_Test("default")
	//Request.CreateService_test("mysql-service", "mysql----test", "default", 3302)
	//Request.DeleteService("default", "mysql-service")
	//Request.DeletePod("default", "mysql----test")
	//Request.GenerateServiceBody("mysql-service", "mysql----test", "default", 3302)
	WorkLoadController.MissionRecord()
	//Request.GetAllService_Test("default")
	fmt.Println("tick tick stop")
	//Request.DeletePod("default", "mysql----test")
	//Request.GetAllReplicationcontrollers_Test()
	//Request.CreateReplicationController_test("default", "mysql:5.7", "mysql-rpc", "mysql-test", "mysql----test", 2)
	//Request.DeleteReplicationController("default", "mysql-rpc")


}

