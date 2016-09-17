package workload

import (
	"strconv"
	"fmt"
	"bufio"
	"os"
	"strings"
	Request "../request"
	"../ab"
	"time"
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
	Total      *WorkLoad
	LongTerm   *WorkLoad
	ShortTerm  *WorkLoad
	JobMonitor *channel
	JobWorker  *channel
	JobNum     int
	JobName    []string
}

func Start() {
	var WL *WorkloadController
	var line string
	scanner := bufio.NewScanner(os.Stdin)
	//var WL *WorkloadController
	fmt.Println("^_^")
	fmt.Println("1.查看jobs状态")
	fmt.Println("2.设定WorkLoad参数并部署任务")
	fmt.Println("3.跑测试并记录测试结果到文件")
	fmt.Println("4.停止所有任务")
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
			WL = EndMission(longTermName, longTermService, "default", WL)
		//fmt.Println("3")
		}
		fmt.Println("^_^")
		fmt.Println("1.查看jobs状态")
		fmt.Println("2.设定WorkLoad参数并部署任务")
		fmt.Println("3.跑测试并记录测试结果到文件")
		fmt.Println("4.停止所有任务")
	}
}

func escapeMysqlQuery(path string) string {
	str := strings.Replace(path, "./", "\"./", 1)
	return str + "\""
}

func GetJobStatus(WL *WorkloadController) {
	if WL == nil || len(WL.JobName) == 0 {
		fmt.Println("目前没有在运行的任务")
		return
	}
	for _, jobName := range WL.JobName {
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

func GetLongTermStatus() {
	Pods := Request.GetPodsOfNamespace("default", Request.Caicloud)
	namespace, name := Request.FindPodByLabelName(longTermName, &Pods)
	pod := Request.GetPodByNameAndNamespace(namespace, name, Request.Caicloud)
	fmt.Println("podName: ", pod.Name)
	fmt.Println("status: ", pod.Status.Phase)
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
	//fmt.Print("输入任何字符以结束:\n")
	//scanner.Scan()
	//line = scanner.Text()
	//EndShortTermWorkLoadV2(WL)
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
	WL.JobNum = JobNum
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
	WL.JobName = jobNames
	fmt.Println("haha:             ", WL.JobName)

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
	WL.JobMonitor = &quitMoinitors
	WL.JobWorker = &quitWorkLoaders
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
	names, name := Request.FindPodByLabelName(longTermName, &pods)
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

func EndShortTermWorkLoadV2(WL *WorkloadController) {
	for index := 0; index < WL.JobNum; index++ {
		(*WL.JobMonitor)[index] <- "q"
		(*WL.JobWorker)[index] <- "q"
		jobName := "test" + strconv.Itoa(index)
		Request.DeleteJob("default", jobName, Request.Caicloud)
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
						if Request.PodComplete(Request.GetPodByNameAndNamespace("default", pon, Request.Caicloud)) {
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

