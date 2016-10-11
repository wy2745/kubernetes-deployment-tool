package ab

import (
	"os/exec"
)

const (
	root string = "/Users/panda/Desktop"
	node string = "http://192.168.6.15"
	port string = "8080"
	getfile string = "api/v1/pods"
	desroot string = "/home/administrator/test/ab/"
	tmpurl string = node + ":" + port + "/" + getfile
)

func Abtest(name string, count string) {
	var err error
	str := "ab -k -n 2000 -c 10 -e " + desroot + "record" + name + "-" + count + ".csv" + " -g " + desroot + "record" + name + "-" + count + ".gnp " + tmpurl
	//fmt.Println(str)
	cmd := exec.Command("/bin/sh", "-c", str)
	_, err = cmd.Output()
	err = cmd.Start()
	err = cmd.Wait()
	if err == nil {

	}

}