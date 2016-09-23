package workload

import (
	"strconv"
	"fmt"
	"bufio"
	"os"
	"strings"
	Request "../request"
	Locust "../locust"
	"../interf"
	"../ab"
	"time"
	"math/rand"
)

const (
	cpu string = "m"
	mem_G string = "Gi"
	mem_M string = "Mi"
	short_general string = "docker.io/zilinglius/workload:short-general"
	short_cpu_bound string = "docker.io/zilinglius/workload:short-cpu-bound"
	nginx_image string = "ymqytw/nginxhttps:1.5"
	shortTermWorkCompletion int32 = 1000
	cpu_Use = 400
	mem_Use = 10
	standardCpu_use = 400
	standardMem_use = 200
	longTermName = "ntest"
	longTermService = "nservice"
)

type WorkLoad struct {
	cpuWorkLoad     string
	memWorkLoad     string
	cpuWorkLoad_int int
	memWorkLoad_int int
}

type channel []chan string

type WorkloadController struct {
	Total     *WorkLoad
	LongTerm  *WorkLoad
	ShortTerm *WorkLoad
	Monitor   *channel
	Worker    *channel
	Num       int
	Name      []string
	NodeName  []string
}

func StartJobVersion() {
	var WL *WorkloadController
	var line string
	scanner := bufio.NewScanner(os.Stdin)
	//var WL *WorkloadController
	fmt.Println("^_^")
	fmt.Println("1.查看jobs状态")
	fmt.Println("2.设定WorkLoad参数并部署任务")
	fmt.Println("3.跑ab测试并记录测试结果到文件")
	fmt.Println("4.Locust")
	fmt.Println("5.停止所有任务")
	fmt.Println("6.退出")
	for {
		scanner.Scan()
		line = scanner.Text()
		switch line {
		case "1":
			//fmt.Println("1")
			GetJobStatus(WL)
			GetLongTermStatus()
		case "2":
			if WL != nil {
				WL = EndMission(longTermName, longTermService, "default", WL)
			}
			WL = record(scanner)
		//UploadLongTermService(longTermName, longTermService, 80, 30888)
		//fmt.Println("2")
		case "3":
			if WL != nil {
				fmt.Println("准备测试...")
				ab.Abtest()
				fmt.Println("测试完成，文件储存完成")
			}
		case "4":
			Locust.Locust(scanner)
		case "5":
			WL = EndMission(longTermName, longTermService, "default", WL)
		case "6":
			if WL != nil {
				WL = EndMission(longTermName, longTermService, "default", WL)
			}
			return
		//fmt.Println("3")
		}
		fmt.Println("1.查看jobs状态")
		fmt.Println("2.设定WorkLoad参数并部署任务")
		fmt.Println("3.跑ab测试并记录测试结果到文件")
		fmt.Println("4.Locust")
		fmt.Println("5.停止所有任务")
		fmt.Println("6.退出")

	}
}

func StartPodVersion() {
	var WL *WorkloadController
	var line string
	scanner := bufio.NewScanner(os.Stdin)
	mode := Request.Test
	//var WL *WorkloadController
	fmt.Println("^_^")
	fmt.Println("1.查看状态")
	fmt.Println("2.设定WorkLoad参数并部署任务")
	fmt.Println("3.退出")
	for {
		scanner.Scan()
		line = scanner.Text()
		switch line {
		case "1":
			GetPodStatus(WL, mode)
		case "2":
			if WL != nil {
				WL = WL.EndWorkLoad(mode)
			}
			WL = ScheduleTest(scanner, mode)
		case "3":
			if WL != nil {
				WL = WL.EndWorkLoad(mode)
			}
		case "4":
			if WL != nil {
				WL = WL.EndWorkLoad(mode)
			}
			return
		//fmt.Println("3")
		}
		fmt.Println("1.查看状态")
		fmt.Println("2.设定WorkLoad参数并部署任务")
		fmt.Println("3.退出")

	}
}

func escapeMysqlQuery(path string) string {
	str := strings.Replace(path, "./", "\"./", 1)
	return str + "\""
}

func GetJobStatus(WL *WorkloadController) {
	if WL == nil || len(WL.Name) == 0 {
		fmt.Println("目前没有在运行的任务")
		return
	}
	for _, jobName := range WL.Name {
		job := Request.GetJobByNameAndNamespace("default", jobName, Request.Caicloud)
		fmt.Println("JobName: ", job.Name)
		fmt.Println("    actice: ", job.Status.Active)
		fmt.Println("    succeeded: ", job.Status.Succeeded)
		fmt.Println("    goal: ", *job.Spec.Completions)
		fmt.Println("    failed: ", job.Status.Failed)
		if *job.Spec.Completions == job.Status.Succeeded {
			fmt.Println("Status: Finished.")
		} else {
			fmt.Println("Status: Not Finished.")
		}
		fmt.Println("total: ", job.Spec.Completions, " finished: ", job.Status.Succeeded)
	}
}
func GetPodStatus(WL *WorkloadController, mode string) {
	if WL == nil || len(WL.Name) == 0 {
		fmt.Println("目前没有在运行的Pod")
		return
	}
	for _, podName := range WL.Name {
		pod := Request.GetPodByNameAndNamespace("default", podName, mode)
		fmt.Println("PodName: ", interf.GetName(*pod))
		fmt.Println("    state: ", interf.GetStautsPhase(*pod))
	}
}

func GetLongTermStatus() {
	Pods := Request.GetPodsOfNamespace("default", Request.Caicloud)
	namespace, name := Request.FindPodByLabelName(longTermName, Pods)
	pod := Request.GetPodByNameAndNamespace(namespace, name, Request.Caicloud)
	fmt.Println("podName: ", interf.GetName(*pod))
	fmt.Println("status: ", interf.GetStautsPhase(*pod))
	Service := Request.GetServicesOfNamespaceAndName("default", longTermService, Request.Caicloud)
	fmt.Println("serviceName: ", Service.Name)
	for _, port := range Service.Spec.Ports {
		fmt.Println("expose port:", port.NodePort, "(container port:", port.Port, ") in each node")
	}
	fmt.Println("status: ", Service.Status)
}

func record(scanner *bufio.Scanner, ) *WorkloadController {
	fmt.Print("^_^\n")
	fmt.Print("请输入总负载信息\n")
	fmt.Print("总负载Cpu(单位:m 整数):")
	var TotalCpu, TotalMem int
	var ShortCpuRate, ShortMemRate float64
	var err error
	var line string
	var WL *WorkloadController
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
	WL = NewWorkLoadController(TotalCpu, TotalMem, float32(ShortCpuRate), float32(ShortMemRate))
	WorkLoadDisplay(WL)

	fmt.Println("部署短时任务")

	boundNum := WL.ShortTerm.cpuWorkLoad_int / cpu_Use
	fmt.Println("size: ", boundNum)

	UploadShortTermMissionV2(boundNum, WL)
	UploadLongTermMission(WL.LongTerm.cpuWorkLoad_int, WL.LongTerm.cpuWorkLoad_int, WL.LongTerm.memWorkLoad_int, WL.LongTerm.memWorkLoad_int)
	UploadLongTermService(longTermName, longTermService, 80, 30888)
	return WL
}

func MissionRecorder() {
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
		//fmt.Println("部署长时任务")
		//Request.CreatePod("default", short_general, "test1", WL.LongTerm.cpuWorkLoad, WL.LongTerm.cpuWorkLoad, WL.LongTerm.memWorkLoad, WL.LongTerm.memWorkLoad, "./home/ws 20000000", Request.Caicloud)
		////为了测试，先删掉刚创建好的Pod
		//Request.DeletePod("default", "test1", Request.Caicloud)
		fmt.Println("部署短时任务")
		//resultChan := make(chan string)


		//每个cpu-bound占用10M内存，目前暂时用0.4个cpu
		//目前暂时考虑负载mem = 10*n  cpu = 400*n n是同一个数
		//TODO:根据负载，调整短期任务的内存和cpu占用，以保证能够通过整数数量的短期任务来满足负载
		memUse := 10
		cpuUse := 400

		boundNum := WL.ShortTerm.cpuWorkLoad_int / cpuUse
		fmt.Println("size: ", boundNum)

		qw, qm := UploadShortTermMission(boundNum, cpuUse, cpuUse, memUse, memUse)
		fmt.Print("输入任何字符以结束:\n")
		scanner.Scan()
		line = scanner.Text()
		EndShortTermWorkLoad(qw, qm, boundNum)
		//t1.Stop()
		//Request.DeletePod("default", "test")

	}
}
func UploadShortTermMissionV2(JobNum int, WL *WorkloadController) {
	WL.Num = JobNum
	var jobNames []string
	//建立channel组，分别监听
	var resultChans [] chan string
	var quitWorkLoaders channel
	var quitMoinitors channel
	//生成并记录job组的名称
	for index := 0; index < JobNum; index++ {
		var ind int
		var pon string
		ind = index
		jobName := "test" + strconv.Itoa(ind)
		fmt.Println(jobName)
		jobNames = append(jobNames, jobName)
		resultChan := make(chan string)
		resultChans = append(resultChans, resultChan)
		quitWorkLoader := make(chan string)
		quitWorkLoaders = append(quitWorkLoaders, quitWorkLoader)
		quitMoinitor := make(chan string)
		quitMoinitors = append(quitMoinitors, quitMoinitor)

		go func() {
			for {
				select {
				case pon = <-resultChan:
				//fmt.Println("正在删除旧pod: ", pon)
					Request.DeleteJob("default", pon, Request.Caicloud)
				//fmt.Println("正在创建新pod: ", pon)
				//Request.CreateJob("default", short_cpu_bound, pon, 1000, strconv.Itoa(cpuMin) + cpu, strconv.Itoa(cpuMax) + cpu, strconv.Itoa(memMin) + mem_M, strconv.Itoa(memMax) + mem_M, "./home/wsc 100 1000", Request.Caicloud)
					Request.CreateJobWithoutResource("default", short_cpu_bound, pon, shortTermWorkCompletion, "./home/wsc 100 1000", Request.Caicloud)
				case <-quitWorkLoader:
					fmt.Println("进程杀死", ind)
					return
				default:
					continue
				}
			}
		}()
	}
	WL.Name = jobNames

	//构建一个新的检查job组运行情况的go func，基本思想，建立多个线程，负责单一检查某个job的运行状态，另外建一个线程负责收听job的运行状态，重建job
	for index, jobName := range jobNames {
		var ind int
		var jon string
		jon = jobName
		ind = index
		go func() {
			for {
				select {
				case <-quitMoinitors[ind]:
					fmt.Println("进程杀死", ind)
					return
				default:
					if !Request.JobExist("default", jon, Request.Caicloud) || Request.JobComplete(Request.GetJobByNameAndNamespace("default", jon, Request.Caicloud)) {
						//fmt.Print("不存在，需要创建\n\n\n")
						resultChans[ind] <- jon
						time.Sleep(time.Second * 3)
					}
				}
			}
		}()
	}
	WL.Monitor = &quitMoinitors
	WL.Worker = &quitWorkLoaders
}

func UploadLongTermService(labelName string, name string, port int32, nodeport int32) {
	Request.CreateService(name, labelName, "default", port, nodeport, Request.Caicloud)
	fmt.Println("Service ", name, " created!")
}

func UploadLongTermMission(cpuMin int, cpuMax int, memMin int, memMax int) {
	var ports map[int32]int32
	ports = make(map[int32]int32)
	ports[int32(443)] = int32(8843)
	ports[int32(80)] = int32(8088)
	Request.CreateReplicationController("default", nginx_image, longTermName, longTermName, longTermName, int32(1), strconv.Itoa(cpuMin) + cpu, strconv.Itoa(cpuMax) + cpu, strconv.Itoa(memMin) + mem_M, strconv.Itoa(memMax) + mem_M, ports, Request.Caicloud)
	fmt.Println("长期任务nginx创建...")
}

func EndMission(podName string, serviceName string, namespace string, WL *WorkloadController) *WorkloadController {
	EndShortTermWorkLoadV2(WL)
	EndLongTermMission(longTermName, longTermService, "default")
	return nil
}

func EndLongTermMission(podName string, serviceName string, namespace string) {
	Request.DeleteService(namespace, serviceName, Request.Caicloud)
	Request.DeleteReplicationController(namespace, podName, Request.Caicloud)
	pods := Request.GetPodsOfNamespace("default", Request.Caicloud)
	names, name := Request.FindPodByLabelName(longTermName, pods)
	Request.DeletePod(names, name, Request.Caicloud)

	fmt.Println("长期任务删除成功")
}

func UploadShortTermMission(JobNum int, cpuMin int, cpuMax int, memMin int, memMax int) (*[]chan string, *[]chan string) {
	var jobNames []string
	//建立channel组，分别监听
	var resultChans [] chan string
	var quitWorkLoaders []chan string
	var quitMoinitors [] chan string
	//生成并记录job组的名称
	for index := 0; index < JobNum; index++ {
		var ind int
		var pon string
		ind = index
		jobName := "test" + strconv.Itoa(ind)
		fmt.Println(jobName)
		jobNames = append(jobNames, jobName)
		resultChan := make(chan string)
		resultChans = append(resultChans, resultChan)
		quitWorkLoader := make(chan string)
		quitWorkLoaders = append(quitWorkLoaders, quitWorkLoader)
		quitMoinitor := make(chan string)
		quitMoinitors = append(quitMoinitors, quitMoinitor)

		go func() {
			for {
				select {
				case pon = <-resultChan:
				//fmt.Println("正在删除旧pod: ", pon)
					Request.DeleteJob("default", pon, Request.Caicloud)
				//fmt.Println("正在创建新pod: ", pon)
				//Request.CreateJob("default", short_cpu_bound, pon, 1000, strconv.Itoa(cpuMin) + cpu, strconv.Itoa(cpuMax) + cpu, strconv.Itoa(memMin) + mem_M, strconv.Itoa(memMax) + mem_M, "./home/wsc 100 1000", Request.Caicloud)
					Request.CreateJobWithoutResource("default", short_cpu_bound, pon, shortTermWorkCompletion, "./home/wsc 100 1000", Request.Caicloud)
				case <-quitWorkLoader:
					fmt.Println("进程杀死", ind)
					return
				default:
					continue
				}
			}
		}()
	}

	//构建一个新的检查job组运行情况的go func，基本思想，建立多个线程，负责单一检查某个job的运行状态，另外建一个线程负责收听job的运行状态，重建job
	for index, jobName := range jobNames {
		var ind int
		var jon string
		jon = jobName
		ind = index
		go func() {
			for {
				select {
				case <-quitMoinitors[ind]:
					fmt.Println("进程杀死", ind)
					return
				default:
					if !Request.JobExist("default", jon, Request.Caicloud) || Request.JobComplete(Request.GetJobByNameAndNamespace("default", jon, Request.Caicloud)) {
						//fmt.Print("不存在，需要创建\n\n\n")
						resultChans[ind] <- jon
						time.Sleep(time.Second * 3)
					}
				}
			}
		}()
	}
	return &quitWorkLoaders, &quitMoinitors
}
func UploadShortTermPodMission(PodNum int, WL *WorkloadController) {
	WL.Num = PodNum
	var podNames []string
	//建立channel组，分别监听
	var resultChans [] chan string
	var quitWorkLoaders channel
	var quitMoinitors channel
	//生成并记录pod组的名称
	for index := 0; index < PodNum; index++ {
		var ind int
		var pon string
		ind = index
		podName := "test" + strconv.Itoa(ind)
		fmt.Println(podName)
		podNames = append(podNames, podName)
		resultChan := make(chan string)
		resultChans = append(resultChans, resultChan)
		quitWorkLoader := make(chan string)
		quitWorkLoaders = append(quitWorkLoaders, quitWorkLoader)
		quitMoinitor := make(chan string)
		quitMoinitors = append(quitMoinitors, quitMoinitor)

		go func() {
			for {
				select {
				case pon = <-resultChan:
					fmt.Println("正在删除旧pod: ", pon)
					Request.DeletePod("default", pon, Request.Caicloud)
					fmt.Println("正在创建新pod: ", pon)
					Request.CreatePod("default", short_cpu_bound, pon, strconv.Itoa(cpu_Use) + cpu, strconv.Itoa(cpu_Use) + cpu, strconv.Itoa(mem_Use) + mem_M, strconv.Itoa(mem_Use) + mem_M, "./home/wsc 100 1000", Request.Caicloud)
				case <-quitWorkLoader:
					fmt.Println("进程杀死", ind)
					return
				default:
					continue
				}
			}
		}()
	}
	WL.Name = podNames

	//构建一个新的检查pod组运行情况的go func，基本思想，建立多个线程，负责单一检查某个pod的运行状态，另外建一个线程负责收听pod的运行状态，重建pod
	for index, podName := range podNames {
		var ind int
		var pon string
		pon = podName
		ind = index
		go func() {
			for {
				select {
				case <-quitMoinitors[ind]:
					fmt.Println("进程杀死", ind)
					return
				default:
					if Request.PodComplete(*Request.GetPodByNameAndNamespace("default", pon, Request.Caicloud)) {
						fmt.Print("不存在，需要创建\n\n\n")
						resultChans[ind] <- pon
						time.Sleep(time.Second * 3)
					}
				}
			}
		}()
	}
	WL.Monitor = &quitMoinitors
	WL.Worker = &quitWorkLoaders

}

func EndShortTermWorkLoadV2(WL *WorkloadController) {
	for index := 0; index < WL.Num; index++ {
		(*WL.Monitor)[index] <- "q"
		(*WL.Worker)[index] <- "q"
		jobName := "test" + strconv.Itoa(index)
		Request.DeleteJob("default", jobName, Request.Caicloud)
	}
}

func EndShortTermPodWorkLoad(WL *WorkloadController) {
	for index := 0; index < WL.Num; index++ {
		(*WL.Monitor)[index] <- "q"
		(*WL.Worker)[index] <- "q"
		podName := "test" + strconv.Itoa(index)
		Request.DeletePod("default", podName, Request.Caicloud)
	}
}
func EndShortTermWorkLoad(quitWorkLoaders *[]chan string, quitMoinitors *[]chan string, boundNum int) {
	for index := 0; index < boundNum; index++ {
		(*quitMoinitors)[index] <- "q"
		(*quitWorkLoaders)[index] <- "q"
		jobName := "test" + strconv.Itoa(index)
		Request.DeleteJob("default", jobName, Request.Caicloud)
	}
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
		//fmt.Println("部署长时任务")
		//Request.CreatePod("default", short_general, "test1", WL.LongTerm.cpuWorkLoad, WL.LongTerm.cpuWorkLoad, WL.LongTerm.memWorkLoad, WL.LongTerm.memWorkLoad, "./home/ws 20000000", Request.Caicloud)
		////为了测试，先删掉刚创建好的Pod
		//Request.DeletePod("default", "test1", Request.Caicloud)
		fmt.Println("部署短时任务")
		//resultChan := make(chan string)


		//每个cpu-bound占用10M内存，目前暂时用0.4个cpu
		//目前暂时考虑负载mem = 10*n  cpu = 400*n n是同一个数
		//TODO:根据负载，调整短期任务的内存和cpu占用，以保证能够通过整数数量的短期任务来满足负载
		memUse := 10
		cpuUse := 400

		boundNum := WL.ShortTerm.cpuWorkLoad_int / cpuUse
		fmt.Println("size: ", boundNum)
		var podNames []string
		//建立channel组，分别监听
		var resultChans [] chan string
		var quitWorkLoaders []chan string
		var quitMoinitors [] chan string
		//生成并记录pod组的名称
		for index := 0; index < boundNum; index++ {
			var ind int
			var pon string
			ind = index
			podName := "test" + strconv.Itoa(ind)
			fmt.Println(podName)
			podNames = append(podNames, podName)
			resultChan := make(chan string)
			resultChans = append(resultChans, resultChan)
			quitWorkLoader := make(chan string)
			quitWorkLoaders = append(quitWorkLoaders, quitWorkLoader)
			quitMoinitor := make(chan string)
			quitMoinitors = append(quitMoinitors, quitMoinitor)

			go func() {
				for {
					select {
					case pon = <-resultChan:
						fmt.Println("正在删除旧pod: ", pon)
						Request.DeletePod("default", pon, Request.Caicloud)
						fmt.Println("正在创建新pod: ", pon)
						Request.CreatePod("default", short_cpu_bound, pon, strconv.Itoa(cpuUse) + cpu, strconv.Itoa(cpuUse) + cpu, strconv.Itoa(memUse) + mem_M, strconv.Itoa(memUse) + mem_M, "./home/wsc 100 1000", Request.Caicloud)
					case <-quitWorkLoader:
						fmt.Println("进程杀死", ind)
						return
					default:
						continue
					}
				}
			}()
		}

		//构建一个新的检查pod组运行情况的go func，基本思想，建立多个线程，负责单一检查某个pod的运行状态，另外建一个线程负责收听pod的运行状态，重建pod
		for index, podName := range podNames {
			var ind int
			var pon string
			pon = podName
			ind = index
			go func() {
				for {
					select {
					case <-quitMoinitors[ind]:
						fmt.Println("进程杀死", ind)
						return
					default:
						if Request.PodComplete(*Request.GetPodByNameAndNamespace("default", pon, Request.Caicloud)) {
							fmt.Print("不存在，需要创建\n\n\n")
							resultChans[ind] <- pon
							time.Sleep(time.Second * 3)
						}
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
			quitMoinitors[index] <- "q"
			quitWorkLoaders[index] <- "q"
			podName := "test" + strconv.Itoa(index)
			Request.DeletePod("default", podName, Request.Caicloud)
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

func ScheduleTest(scanner *bufio.Scanner, mode string) *WorkloadController {
	fmt.Print("^_^\n")
	fmt.Print("请输入总负载信息\n")
	fmt.Print("总负载Cpu(单位:m 整数):")
	var TotalCpu, TotalMem int
	var err error
	var line string
	var WL *WorkloadController
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
	WL = NewWorkLoadController(TotalCpu, TotalMem, float32(1), float32(1))
	WorkLoadDisplay(WL)

	fmt.Println("部署连续任务")

	//cpuBound := WL.ShortTerm.cpuWorkLoad_int / standardCpu_use
	//memBound := WL.ShortTerm.memWorkLoad_int / standardMem_use
	//var boundNum int
	//if cpuBound >= memBound {
	//	boundNum = cpuBound
	//} else {
	//	boundNum = memBound
	//}
	//fmt.Println("实际调整负载为，cpu: ", boundNum * standardCpu_use, cpu, " mem: ", boundNum * standardMem_use, mem_M)
	//
	//fmt.Println("size: ", boundNum)
	//WL.Num = boundNum
	WL.UploadWorkLoad(mode)
	fmt.Println("搞定")
	scanner.Scan()
	line = scanner.Text()

	log(WL, mode)
	return WL
}
func contain(nodeNames *[]string, nodeName string) bool {
	for _, node := range *nodeNames {
		if node == nodeName {
			return true
		}
	}
	return false
}

func (Wl *WorkloadController) UploadWorkLoad(mode string) {
	if Wl == nil {
		return
	}
	allocMap := allocLoad(Wl.ShortTerm.cpuWorkLoad_int, Wl.ShortTerm.memWorkLoad_int)
	count := 0
	for key, value := range allocMap {
		for key2, value2 := range value {
			for i := 0; i < value2; i++ {
				podName := "test" + strconv.Itoa(count)
				Wl.Name = append(Wl.Name, podName)
				Request.CreatePod("default", short_cpu_bound, podName, strconv.Itoa(key) + cpu, strconv.Itoa(key) + cpu, strconv.Itoa(key2) + mem_M, strconv.Itoa(key2) + mem_M, "./home/wsc 100 1000", mode)
				//time.Sleep(time.Second * 2)
				count++
			}
		}
	}
	Wl.Num = count

	//for index := 0; index < Wl.Num; index++ {
	//	podName := "test" + strconv.Itoa(index)
	//	Wl.Name = append(Wl.Name, podName)
	//	Request.CreatePod("default", short_cpu_bound, podName, strconv.Itoa(standardCpu_use) + cpu, strconv.Itoa(standardCpu_use) + cpu, strconv.Itoa(standardMem_use) + mem_M, strconv.Itoa(standardMem_use) + mem_M, "./home/wsc 100 1000", mode)
	//	time.Sleep(time.Second * 2)
	//}
}

func log(Wl *WorkloadController, mode string) {
	inputFile, inputError := os.Create("/Users/panda/Desktop/record.txt")
	if inputError != nil {
		fmt.Printf("An error occurred on opening the inputfile\n" +
			"Does the file exist?\n" +
			"Have you got acces to it?\n")
		return // exit the function on error
	}
	defer inputFile.Close()
	var cpusum int = 0
	var memsum int = 0
	for index := 0; index < Wl.Num; index++ {
		pod := Request.GetPodByNameAndNamespace("default", Wl.Name[index], mode)
		cpusum += getCpuInt(interf.GetCpu(*pod))
		memsum += getMemInt(interf.GetMem(*pod))
		//inputFile.WriteString("PodName: " + interf.GetName(*pod) + "被部署到了node " + interf.GetNodeName(*pod) + "上")
		inputFile.WriteString("PodName: " + interf.GetName(*pod))
		if ( Wl.NodeName != nil && contain(&Wl.NodeName, interf.GetNodeName(*pod))) {
			inputFile.WriteString("为旧node")
			str := "pod 负载为" + " cpu: " + strconv.Itoa(cpusum) + cpu + " memory: " + strconv.Itoa(memsum) + mem_M + "\n"
			inputFile.WriteString(str)
			//_, str2 := Request.GetNodeByName(interf.GetNodeName(*pod), mode)
			//inputFile.WriteString(str2)
			//fmt.Println("为旧node，node状态如下")
		} else {
			Wl.NodeName = append(Wl.NodeName, interf.GetNodeName(*pod))
			//fmt.Println("为新node，node状态如下")
			Request.GetNodeByName(interf.GetNodeName(*pod), mode)
			inputFile.WriteString("为新node")
			str := "pod 负载为" + " cpu: " + strconv.Itoa(cpusum) + cpu + " memory: " + strconv.Itoa(memsum) + mem_M + "\n"
			inputFile.WriteString(str)
			//_, str2 := Request.GetNodeByName(interf.GetNodeName(*pod), mode)
			//inputFile.WriteString(str2)
		}
	}
}

func (Wl *WorkloadController) EndWorkLoad(mode string) *WorkloadController {
	if Wl == nil {
		return nil
	}
	for _, podName := range Wl.Name {
		Request.DeletePod("default", podName, mode)
		fmt.Println("Pod ", podName, "已删除")
	}
	fmt.Println("workLoad 删除完毕")
	return nil
}

func allocLoad(totalCpu int, totalMem int) map[int](map[int]int) {
	var result map[int](map[int]int) = make(map[int](map[int]int))

	for {
		var result1 map[int](map[int]int) = make(map[int](map[int]int))
		var cpu = totalCpu
		var mem = totalMem
		for i := 100; i <= 400; i += 100 {
			maxCpu := cpu / i
			s1 := rand.NewSource(time.Now().UnixNano())
			r1 := rand.New(s1)
			if maxCpu <= 0 {
				break
			}
			ran1 := r1.Intn(maxCpu)
			if ran1 == 0 {
				continue
			}
			result2, memuse, num := allocMem(i * ran1 * totalMem / totalCpu, ran1)
			result1[i] = result2
			mem -= memuse
			cpu -= num * i
			if mem <= 0 || cpu <= 0 {
				break
			}
		}
		if mem <= 0 || cpu <= 0 {
			result = result1
			break
		}
	}

	fmt.Println("分配情况如下")
	fmt.Println("原负载是")
	fmt.Println(" cpu: ", totalCpu)
	fmt.Println(" mem: ", totalMem)
	ncpu := 0
	nmem := 0
	for key, value := range result {
		for key2, value2 := range value {
			fmt.Println("cpu(", key, ") mem(", key2, ")数量:", value2)
			ncpu += key * value2
			nmem += key2 * value2
		}
	}
	fmt.Println("均衡后负载是")
	fmt.Println(" cpu: ", ncpu)
	fmt.Println(" mem: ", nmem)
	return result
}
func allocMem(mem int, num int) (map[int]int, int, int) {
	memtmp := mem
	numtmp := num
	memuse := 0
	for {
		memtmp = mem
		numtmp = num
		memuse = 0
		result := make(map[int]int)
		fmt.Println(num, "  ", numtmp, " ", memuse)
		for i := 50; i <= 200; i += 50 {
			s1 := rand.NewSource(time.Now().UnixNano())
			r1 := rand.New(s1)
			ran1 := r1.Intn(numtmp) + 1
			fmt.Println("ran", ran1)
			numtmp -= ran1
			fmt.Println("numtmp", numtmp)
			memuse += ran1 * i
			memtmp -= ran1 * i
			result[i] = ran1
			if numtmp <= 0 || memtmp <= 0 {
				return result, memuse, num - numtmp
			}
		}
	}
}
func getCpuInt(cpus string) int {
	str := strings.Split(cpus, cpu)
	in, _ := strconv.Atoi(str[0])
	return in
}
func getMemInt(mems string) int {
	str := strings.Split(mems, mem_M)
	in, _ := strconv.Atoi(str[0])
	return in
}


