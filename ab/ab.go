package ab

import (
	"os/exec"
)

const (
	root string = "/Users/panda/Desktop"
	node string = "http://202.120.40.178"
	port string = "1588"
	getfile string = "api/v1/pods"
	desroot string = "/home/administrator/test/"
	tmpurl string = node + ":" + port + "/" + getfile
)

func Abtest(name string, count string) {
	var err error
	str := "ab -n 10 -c 10 -e " + desroot + "record" + name + "-" + count + ".csv" + " " + tmpurl
	cmd := exec.Command("/bin/sh", "-c", str)
	_, err = cmd.Output()
	err = cmd.Start()
	err = cmd.Wait()
	if err == nil {

	}

}