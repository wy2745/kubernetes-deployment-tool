package autoscale

import (
	"github.com/wy2745/kubernetes-deployment-tool/kubemark"
	"github.com/wy2745/kubernetes-deployment-tool/json"
	classType "github.com/wy2745/kubernetes-deployment-tool/type137"
	"fmt"
)

func BuildNginx() {
	url := kubemark.DestinationServer_Test + kubemark.GenerateReplicationControllerNamespaceUrl("default")
}

func generateNginxReplic(replic int32) []byte {
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
	volumn.Name = "secret-volume"
	volumn.Secret.SecretName = "nginxsecret"
	var volumn2 classType.Volume
	volumn2.Name = "configmap-volume"
	volumn2.Secret.SecretName = "nginxconfigmap"
	vslice := []classType.Volume{volumn, volumn2}

	//生成PodTemplateSpec.container
	var Ports []classType.ContainerPort
	var port classType.ContainerPort
	port.ContainerPort = 443
	Ports = append(Ports, port)
	port.ContainerPort = 80
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
	httpgetaction.Port.IntVal = 80
	probe.Handler.HTTPGet = &httpgetaction
	probe.InitialDelaySeconds = int32(30)
	probe.TimeoutSeconds = int32(1)

	var container classType.Container
	container.Name = containerName
	container.Image = image
	container.Ports = Ports
	container.VolumeMounts = vmslice
	container.LivenessProbe = &probe

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
	replicationController.Status.Replicas = replic

	b := jsonParse.JsonMarsha(replicationController)
	fmt.Print(string(b))
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
	servicePort.Port = int32(80)
	servicePort.NodePort = int32(10080)
	var servicePort2 classType.ServicePort
	servicePort2.Port = int32(443)
	servicePort2.NodePort = int32(10443)
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

