package interf

type Podinface interface {
	GetName() string
	GetLabel(name string) string
	GetNamespace() string
	GetNodeName() string
	GetCpu() string
	GetMem() string
	GetStautsPhase() string
	GetContainerStatusesLen() int
	GetReady() bool
	SetTypeMeta(apiversion string, kind string)
	SetObjectMeta(labelname string, namespace string, name string)
	SetContainer(name string, image string, cpu string, mem string, command string, restart string)
	SetPod(apiversion string, kind string, name string, namespace string, image string, cpu string, mem string, command string) []byte
}
type TypeMetainface interface {
	SetTmProp(apiversion string, kind string)
}
type ObjectMetainface interface {
	SetOmProp(labelname string, namespace string, name string)
}
type Containerinface interface {

}

type PodListinface interface {
	GetItems() []*Podinface
}
type Resourceinface interface {
	SetResource(cpu_min string, cpu_max string, mem_max string, mem_min string)
}

func GetCpu(p Podinface) string {
	return p.GetCpu()
}
func GetMem(p Podinface) string {
	return p.GetMem()
}

func GetNamespace(p Podinface) string {
	return p.GetNamespace()
}

func GetLabel(p Podinface, name string) string {
	return p.GetLabel(name)
}
func SetPod(p Podinface, apiversion string, kind string, name string, namespace string, image string, cpu string, mem string, command string) []byte {
	return p.SetPod(apiversion, kind, name, namespace, image, cpu, mem, command)
}
func SetContainer(p Podinface, name string, image string, cpu string, mem string, command string, restart string) {
	p.SetContainer(name, image, cpu, mem, command, restart)
}

func SetTypeMeta(p Podinface, apiVersion string, kind string) {
	p.SetTypeMeta(apiVersion, kind)
}
func SetObjectMeta(p Podinface, labelname string, namespace string, name string) {
	p.SetObjectMeta(labelname, namespace, name)
}

func GetName(p Podinface) string {
	return p.GetName()
}

func GetNodeName(p Podinface) string {
	return p.GetNodeName()
}
func GetStautsPhase(p Podinface) string {
	return p.GetStautsPhase()
}
func GetContainerStatusesLen(p Podinface) int {
	return p.GetContainerStatusesLen()
}
func GetReady(p Podinface) bool {
	return p.GetReady()
}
func GetItems(pl PodListinface) []*Podinface {
	return pl.GetItems()
}
func SetResource(rs Resourceinface, cpu_min string, cpu_max string, mem_max string, mem_min string) {
	rs.SetResource(cpu_min, cpu_max, mem_max, mem_min)
}
func SetTmProp(tm TypeMetainface, apiversion string, kind string) {
	tm.SetTmProp(apiversion, kind)
}
func SetOmProp(ob ObjectMetainface, labelname string, namespace string, name string) {
	ob.SetOmProp(labelname, namespace, name)
}

