package ab

import (
	"os/exec"
)

const (
	root string = "/Users/panda/Desktop"
	node string = "http://192.168.6.15"
	port string = "8080"
	getfile string = "api/v1/pods"
	getfile2 string = "api/v1/nodes"
	desroot string = "/home/administrator/test/ab/"
	tmpurl string = node + ":" + port + "/" + getfile2
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
func AbtestV2(name string, count string) {
	csvName := desroot + "record" + name + "-" + count + ".csv"
	gnpName := desroot + "record" + name + "-" + count + ".gnp"
	//fmt.Println("~/go/src/github.com/wy2745/kubernetes-deployment-tool/abTest.sh " + csvName + " " + gnpName + " " + tmpurl)
	cmd := exec.Command("/bin/sh", "-c", "~/go/src/github.com/wy2745/kubernetes-deployment-tool/abTest.sh " + csvName + " " + gnpName + " " + tmpurl)
	cmd.Output()
}