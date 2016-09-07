package workload

import (
	"strconv"
	"fmt"
	"bufio"
	"os"
	"strings"
	Request "../request"
	"time"
)

const (
	cpu string = "m"
	mem_G string = "Gi"
	mem_M string = "Mi"
	short_general string = "docker.io/zilinglius/workload:short-general"
	short_cpu_bound string = "docker.io/zilinglius/workload:short-cpu-bound"
)

type WorkLoad struct {
	cpuWorkLoad     string
	memWorkLoad     string
	cpuWorkLoad_int int
	memWorkLoad_int int
}

type WorkloadController struct {
	Total     *WorkLoad
	LongTerm  *WorkLoad
	ShortTerm *WorkLoad
}

func escapeMysqlQuery(path string) string {
	str := strings.Replace(path, "./", "\"./", 1)
	return str + "\""
}

func MissionRecord() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("^_^\n")
		fmt.Print("请输入总负载信息\n")
		fmt.Print("总负载Cpu(单位:m 整数):")
		var TotalCpu, TotalMem int
		var ShortCpuRate, ShortMemRate float64
		var err error
		var line string
		scanner.Scan()
		line = scanner.Text()

		for TotalCpu, err = strconv.Atoi(line); err != nil; {
			fmt.Print("必须输入整数\n")
			scanner.Scan()
			line = scanner.Text()
			TotalCpu, err = strconv.Atoi(line)

		}
		fmt.Print("\n总负载Mem(单位:Mi 整数):")
		scanner.Scan()
		line = scanner.Text()
		for TotalMem, err = strconv.Atoi(line); err != nil; err = nil {
			fmt.Print("必须输入整数\n")
			scanner.Scan()
			line = scanner.Text()
			TotalMem, err = strconv.Atoi(line)
		}
		fmt.Print("\n短时任务负载Cpu比例(小数):")
		scanner.Scan()
		line = scanner.Text()
		for ShortCpuRate, err = strconv.ParseFloat(line, 32); err != nil || ShortCpuRate > 1; err = nil {
			fmt.Print("必须输入小于等于1的小数\n")
			scanner.Scan()
			line = scanner.Text()
			ShortCpuRate, err = strconv.ParseFloat(line, 32)
		}
		fmt.Print("\n短时任务负载Mem比例(小数):")
		scanner.Scan()
		line = scanner.Text()
		for ShortMemRate, err = strconv.ParseFloat(line, 32); err != nil || ShortMemRate > 1; err = nil {
			fmt.Print("必须输入小于等于1的小数\n")
			scanner.Scan()
			line = scanner.Text()
			ShortMemRate, err = strconv.ParseFloat(line, 32)
		}
		var WL *WorkloadController
		WL = NewWorkLoadController(TotalCpu, TotalMem, float32(ShortCpuRate), float32(ShortMemRate))
		WorkLoadDisplay(WL)
		fmt.Println("部署长时任务")
		Request.CreatePod_test("default", short_general, "test1", WL.LongTerm.cpuWorkLoad, WL.LongTerm.cpuWorkLoad, WL.LongTerm.memWorkLoad, WL.LongTerm.memWorkLoad, "./home/ws 20000000")
		//为了测试，先删掉刚创建好的Pod
		Request.DeletePod("default", "test1")
		fmt.Println("部署短时任务")
		//resultChan := make(chan string)


		//每个cpu-bound占用10M内存，目前暂时用0.4个cpu
		//目前暂时考虑负载mem = 10*n  cpu = 400*n n是同一个数
		//TODO:根据负载，调整短期任务的内存和cpu占用，以保证能够通过整数数量的短期任务来满足负载
		memUse := 10
		cpuUse := 400

		boundNum := WL.ShortTerm.cpuWorkLoad_int / cpuUse
		var podNames []string
		//建立channel组，分别监听
		var resultChans [] chan string
		//生成并记录pod组的名称
		for index := 0; index < boundNum; index++ {
			podName := "test" + strconv.Itoa(index)
			append(podNames, podName)
			resultChan := make(chan string)
			append(resultChans, resultChan)
			go func() {
				for {
					podName := <-resultChan
					fmt.Println("正在删除旧pod: ", podName)
					Request.DeletePod("default", podName)
					fmt.Println("正在创建新pod: ", podName)
					Request.CreatePod_test("default", short_cpu_bound, podName, cpuUse, cpuUse, memUse, memUse, "./home/wsc 100 1000")

				}
			}()
		}

		//go func() {
		//	for {
		//
		//		if Request.PodComplete(Request.GetPodByNameAndNamespace("default", "test")) {
		//			fmt.Print("不存在，需要创建\n\n\n")
		//			resultChan <- "1"
		//			time.Sleep(time.Second * 3)
		//		}
		//
		//	}
		//}()
		//go func() {
		//	for {
		//		<-resultChan
		//		fmt.Println("正在删除旧pod")
		//		Request.DeletePod("default", "test")
		//		fmt.Println("正在创建新pod")
		//		Request.CreatePod_test("default", short_cpu_bound, "test", WL.ShortTerm.cpuWorkLoad, WL.ShortTerm.cpuWorkLoad, WL.ShortTerm.memWorkLoad, WL.ShortTerm.memWorkLoad, "./home/wsc 100 1000")
		//
		//	}
		//}()

		//构建一个新的检查pod组运行情况的go func，基本思想，建立多个线程，负责单一检查某个pod的运行状态，另外建一个线程负责收听pod的运行状态，重建pod
		for index, podName := range podNames {
			go func() {
				for {

					if Request.PodComplete(Request.GetPodByNameAndNamespace("default", podName)) {
						fmt.Print("不存在，需要创建\n\n\n")
						resultChans[index] <- podName
						time.Sleep(time.Second * 3)
					}

				}
			}()
		}

		fmt.Print("输入任何字符以结束:\n")
		scanner.Scan()
		line = scanner.Text()
		//t1.Stop()
		//Request.DeletePod("default", "test")
		for index := 0; index < boundNum; index++ {
			podName := "test" + strconv.Itoa(index)
			Request.DeletePod("default", podName)
		}
	}

}

func NewWorkLoad(cpu string, mem string, cpu_int int, mem_int int) *WorkLoad {
	w := WorkLoad{
		cpuWorkLoad:cpu,
		memWorkLoad:mem,
		cpuWorkLoad_int:cpu_int,
		memWorkLoad_int:mem_int,
	}
	return &w
}

func NewWorkLoadController(TotalCpu int, TotalMem int, ShortCpuRate float32, ShortMemRate float32) *WorkloadController {
	shortCpu_int := int(float32(TotalCpu) * ShortCpuRate)
	shortMem_int := int(float32(TotalMem) * ShortMemRate)
	shortCpu := strconv.Itoa(shortCpu_int)
	shortMem := strconv.Itoa(shortMem_int)
	longCpu := strconv.Itoa(TotalCpu - shortCpu_int)
	longMem := strconv.Itoa(TotalMem - shortMem_int)
	s := WorkloadController{
		Total : NewWorkLoad(strconv.Itoa(TotalCpu) + cpu, strconv.Itoa(TotalMem) + mem_M, TotalCpu, TotalMem),
		LongTerm:NewWorkLoad(longCpu + cpu, longMem + mem_M, TotalCpu - shortCpu_int, TotalMem - shortMem_int),
		ShortTerm:NewWorkLoad(shortCpu + cpu, shortMem + mem_M, shortCpu_int, shortMem_int),
	}

	return &s
}

func WorkLoadDisplay(wl *WorkloadController) {
	if wl == nil {
		fmt.Print("对象为空")
		return
	}
	fmt.Print("总负载为:\n", "Cpu: ", wl.Total.cpuWorkLoad, "\n", "Mem: ", wl.Total.memWorkLoad, "\n")
	fmt.Print("长时任务负载为:\n", "Cpu: ", wl.LongTerm.cpuWorkLoad, "\n", "Mem: ", wl.LongTerm.memWorkLoad, "\n")
	fmt.Print("短时任务负载为:\n", "Cpu: ", wl.ShortTerm.cpuWorkLoad, "\n", "Mem: ", wl.ShortTerm.memWorkLoad, "\n\n")
}

