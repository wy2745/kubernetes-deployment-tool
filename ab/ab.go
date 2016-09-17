package ab

import (
	"os/exec"
	"fmt"
)

const (
	root string = "/Users/panda/Desktop"
	node string = "http://120.26.120.30"
	port string = "30888"
	getfile string = "index.html"
	desroot string = "/Users/panda/Desktop/result.cs"
	tmpurl string = node + ":" + port + "/" + getfile
)

func Abtest() {
	var err error
	str := "ab -n 10 -c 10 -e " + desroot + " " + tmpurl
	fmt.Println(str)
	cmd := exec.Command("/bin/sh", "-c", str)
	_, err = cmd.Output()
	err = cmd.Start()
	err = cmd.Wait()
	if err == nil {

	}

}