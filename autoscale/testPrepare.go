package autoscale

import (
	"github.com/wy2745/kubernetes-deployment-tool/kubemark"
	"github.com/wy2745/kubernetes-deployment-tool/json"
	classType "github.com/wy2745/kubernetes-deployment-tool/type137"
	"fmt"
	"io/ioutil"
	"os/exec"
	"net/http"
	"time"
)

func BuildNginx(num int32, cpu string) {
	url := kubemark.DestinationServer_Test2 + kubemark.GenerateReplicationControllerNamespaceUrl("default")
	//fmt.Println(url)
	body := generateNginxReplic(num, cpu)
	//fmt.Println(string(body))
	tr := http.Transport{DisableKeepAlives:false}
	client := http.Client{Transport:&tr}
	resp := kubemark.InvokeRequestV2("POST", url, body, &client)
	if (resp != nil) {
		defer resp.Body.Close()
		var v classType.ReplicationController
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
		}
		jsonParse.JsonUnmarsha(body, &v)
		//fmt.Println(v)
	}

	url = kubemark.DestinationServer_Test2 + kubemark.GenerateServiceListNamespaceUrl("default")
	//fmt.Println(url)
	body = generateNginxsvc()
	//fmt.Println(string(body))
	resp = kubemark.InvokeRequestV2("POST", url, body, &client)
	if (resp != nil) {
		defer resp.Body.Close()
		var v classType.Service
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
		}
		jsonParse.JsonUnmarsha(body, &v)
		//fmt.Println(v)
	}

	for {
		count := int32(0)
		url := kubemark.DestinationServer_Test2 + kubemark.GeneratePodNamespaceUrl("default")
		resp := kubemark.InvokeRequestV2("GET", url, nil, &client)
		if (resp != nil) {
			defer resp.Body.Close()
			var v classType.PodList
			body, err := ioutil.ReadAll(resp.Body)
			if (err != nil) {
				fmt.Print(err)
			}
			jsonParse.JsonUnmarsha(body, &v)
			for _, pod := range v.Items {
				if pod.Labels["name"] == "nginx" && pod.Status.Phase == "Running" {
					for _, pc := range pod.Status.Conditions {
						if pc.Type == "Ready" && pc.Status == "True" {
							count ++
						}
					}

				}
			}
		}
		if count == num {
			cmd := exec.Command("/bin/sh", "-c", "~/go/src/github.com/wy2745/kubernetes-deployment-tool/autoscale.sh")
			cmd.Output()
			return
		}
	}

}
func BuildNginxForLB(num int32, cpu string) {

	url := kubemark.DestinationServer_Test2 + kubemark.GenerateReplicationControllerNamespaceUrl("default")
	//fmt.Println(url)
	body := generateNginxReplic(num, cpu)
	//fmt.Println(string(body))
	tr := http.Transport{DisableKeepAlives:false}
	client := http.Client{Transport:&tr}
	resp := kubemark.InvokeRequestV2("POST", url, body, &client)
	if (resp != nil) {
		defer resp.Body.Close()
		var v classType.ReplicationController
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
		}
		jsonParse.JsonUnmarsha(body, &v)
		//fmt.Println(v)
	}

	url = kubemark.DestinationServer_Test2 + kubemark.GenerateServiceListNamespaceUrl("default")
	//fmt.Println(url)
	body = generateNginxsvc()
	//fmt.Println(string(body))
	resp = kubemark.InvokeRequestV2("POST", url, body, &client)
	if (resp != nil) {
		defer resp.Body.Close()
		var v classType.Service
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
		}
		jsonParse.JsonUnmarsha(body, &v)
		//fmt.Println(v)
	}

	for {
		count := int32(0)
		url := kubemark.DestinationServer_Test2 + kubemark.GeneratePodNamespaceUrl("default")
		resp := kubemark.InvokeRequestV2("GET", url, nil, &client)
		if (resp != nil) {
			defer resp.Body.Close()
			var v classType.PodList
			body, err := ioutil.ReadAll(resp.Body)
			if (err != nil) {
				fmt.Print(err)
			}
			jsonParse.JsonUnmarsha(body, &v)
			for _, pod := range v.Items {
				if pod.Labels["name"] == "nginx" && pod.Status.Phase == "Running" {
					for _, pc := range pod.Status.Conditions {
						if pc.Type == "Ready" && pc.Status == "True" {
							count ++
						}
					}

				}
			}
		}
		fmt.Println(count)
		if count == num {
			cmd := exec.Command("/bin/sh", "-c", "~/go/src/github.com/wy2745/kubernetes-deployment-tool/loadBalanceTest.sh")
			cmd.Output()
			return
		}
	}

}

func DestoryNginx() {
	replicName := "nginx"
	svcName := "nginx-svc"

	tr := http.Transport{DisableKeepAlives:false}
	client := http.Client{Transport:&tr}
	url := kubemark.DestinationServer_Test2 + kubemark.GenerateServiceListNameUrl("default", svcName)
	kubemark.InvokeRequestV2("DELETE", url, nil, &client)

	url = kubemark.DestinationServer_Test2 + kubemark.GenerateReplicationControllerNameUrl("default", replicName)
	kubemark.InvokeRequestV2("DELETE", url, nil, &client)

	url = kubemark.DestinationServer_Test2 + kubemark.GeneratePodNamespaceUrl("default")
	resp := kubemark.InvokeRequestV2("GET", url, nil, &client)
	if (resp != nil) {
		defer resp.Body.Close()
		var v classType.PodList
		body, err := ioutil.ReadAll(resp.Body)
		if (err != nil) {
			fmt.Print(err)
		}
		jsonParse.JsonUnmarsha(body, &v)
		var clients []http.Client
		for i := 0; i < len(v.Items); i++ {
			tr := http.Transport{DisableKeepAlives:false}
			clienttmp := http.Client{Transport:&tr}
			clients = append(clients, clienttmp)
		}
		for index, pod := range v.Items {
			if pod.Labels["name"] == replicName {
				url = kubemark.DestinationServer_Test2 + kubemark.GeneratePodNameUrl("default", pod.Name)
				kubemark.InvokeRequestV2("DELETE", url, nil, &(clients[index]))
			}
		}
		for {
			url = kubemark.DestinationServer_Test2 + kubemark.GeneratePodNamespaceUrl("default")
			resp = kubemark.InvokeRequestV2("GET", url, nil, &client)
			if (resp != nil) {
				defer resp.Body.Close()
				var v classType.PodList
				body, err := ioutil.ReadAll(resp.Body)
				if (err != nil) {
					fmt.Print(err)
				}
				jsonParse.JsonUnmarsha(body, &v)
				var complete = true
				for _, pod := range v.Items {
					if pod.Labels["name"] == replicName {
						complete = false
						break
					}
				}
				fmt.Println(complete)
				if complete == true {
					return
				} else {
					time.Sleep(time.Second * 1)
				}
			}
		}
	}
}

func generateNginxReplic(replic int32, cpu string) []byte {
	var name = "nginx"
	var image = "ymqytw/nginxhttps:1.5"
	var containerName = "nginxhttps"
	var labelName = "nginx"
	//生成typeMedata
	var typeMedata classType.TypeMeta
	typeMedata.APIVersion = "v1"
	typeMedata.Kind = "ReplicationController"
	//生成objectMedata
	var objectMedata classType.ObjectMeta
	objectMedata.Labels = make(map[string]string)
	objectMedata.Labels["name"] = labelName
	objectMedata.Name = name
	objectMedata.Namespace = "default"

	//生成PodObjectMedata
	var objectMedata2 classType.ObjectMeta
	objectMedata2.Labels = make(map[string]string)
	objectMedata2.Labels["name"] = labelName
	objectMedata2.Namespace = "default"

	//volumns
	var volumn classType.Volume
	var secret classType.SecretVolumeSource
	secret.SecretName = "nginxsecret"
	volumn.Name = "secret-volume"
	volumn.Secret = &secret
	var volumn2 classType.Volume
	var configmap classType.ConfigMapVolumeSource
	configmap.Name = "nginxconfigmap"
	volumn2.Name = "configmap-volume"
	volumn2.ConfigMap = &configmap
	vslice := []classType.Volume{volumn, volumn2}

	//生成PodTemplateSpec.container
	var Ports []classType.ContainerPort
	var port classType.ContainerPort
	port.ContainerPort = int32(443)
	port.HostPort = int32(4443)
	Ports = append(Ports, port)
	port.ContainerPort = int32(80)
	port.HostPort = int32(8888)
	Ports = append(Ports, port)

	var volumnMount classType.VolumeMount
	volumnMount.Name = "secret-volume"
	volumnMount.MountPath = "/etc/nginx/ssl"
	var volumnMount2 classType.VolumeMount
	volumnMount2.Name = "configmap-volume"
	volumnMount2.MountPath = "/etc/nginx/conf.d"
	vmslice := []classType.VolumeMount{volumnMount, volumnMount2}

	var probe classType.Probe
	var httpgetaction classType.HTTPGetAction
	httpgetaction.Path = "/index.html"
	httpgetaction.Port = int32(80)
	probe.Handler.HTTPGet = &httpgetaction
	probe.InitialDelaySeconds = int32(30)
	probe.TimeoutSeconds = int32(1)

	var resource classType.ResourceRequirements
	resource.Limits = make(map[classType.ResourceName]string)
	resource.Limits["cpu"] = cpu
	resource.Limits["memory"] = "200M"
	resource.Requests = make(map[classType.ResourceName]string)
	resource.Requests["cpu"] = cpu
	resource.Requests["memory"] = "200M"

	var container classType.Container
	container.Name = containerName
	container.Image = image
	container.Ports = Ports
	container.VolumeMounts = vmslice
	container.LivenessProbe = &probe
	container.Resources = resource

	var command []string
	command = append(command, "/home/auto-reload-nginx.sh")
	container.Command = command

	slice := []classType.Container{container}
	var podTemplateSpec classType.PodTemplateSpec
	podTemplateSpec.ObjectMeta = objectMedata2
	podTemplateSpec.Spec.Containers = slice
	podTemplateSpec.Spec.Volumes = vslice

	var replicationControllerSpec classType.ReplicationControllerSpec
	replicationControllerSpec.Template = &podTemplateSpec

	replicationControllerSpec.Replicas = &replic
	replicationControllerSpec.Selector = make(map[string]string)
	replicationControllerSpec.Selector["name"] = labelName

	var replicationController classType.ReplicationController
	replicationController.Spec = replicationControllerSpec
	replicationController.ObjectMeta = objectMedata
	replicationController.TypeMeta = typeMedata
	//replicationController.Status.Replicas = replic

	b := jsonParse.JsonMarsha(replicationController)
	//fmt.Print(string(b))
	return b
}

func generateNginxsvc() []byte {
	var name = "nginx"
	var namespace = "default"
	var typeMedata classType.TypeMeta
	typeMedata.APIVersion = "v1"
	typeMedata.Kind = "Service"

	//生成objectMedata
	var objectMedata classType.ObjectMeta
	objectMedata.Labels = make(map[string]string)
	objectMedata.Labels["name"] = name
	objectMedata.Namespace = namespace
	objectMedata.Name = "nginx-svc"

	//生成Service spec
	var servicePort classType.ServicePort
	servicePort.Name = "http"
	servicePort.Port = int32(80)
	servicePort.NodePort = int32(30080)
	var servicePort2 classType.ServicePort
	servicePort2.Name = "https"
	servicePort2.Port = int32(443)
	servicePort2.NodePort = int32(30443)
	slice := []classType.ServicePort{servicePort, servicePort2}
	var serviceSpec classType.ServiceSpec
	serviceSpec.Selector = make(map[string]string)
	serviceSpec.Selector["name"] = name
	serviceSpec.Ports = slice
	serviceSpec.Type = classType.ServiceTypeNodePort
	var service classType.Service
	service.ObjectMeta = objectMedata
	service.TypeMeta = typeMedata
	service.Spec = serviceSpec
	b := jsonParse.JsonMarsha(service)
	//fmt.Print(string(b))
	return b
}



