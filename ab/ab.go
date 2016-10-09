package ab

import (
	"os/exec"
	"fmt"
)

const (
	root string = "/Users/panda/Desktop"
	node string = "http://202.120.40.177"
	port string = "17080"
	getfile string = "api/v1/pods"
	desroot string = "/Users/panda/Desktop/"
	tmpurl string = node + ":" + port + "/" + getfile
)

func Abtest(name string) {
	var err error
	str := "ab -n 10 -c 10 -e " + desroot + "record" + name + ".csv" + " " + tmpurl
	fmt.Println(str)
	cmd := exec.Command("/bin/sh", "-c", str)
	_, err = cmd.Output()
	err = cmd.Start()
	err = cmd.Wait()
	if err == nil {

	}

}