package subsystems

// 传递资源配置信息的结构体
type ResourceConfig struct {
	MemoryLimit string // 内存限制
	CpuShare    string // CPU时间片权重
	CpuSet      string // CPU核心数
}

// 将cgroup抽象成path，原因是cgroup为hierarchy的路径 -->  虚拟文件系统中的虚拟路径
type Subsystem interface {
	Name() string
	Set(path string, res *ResourceConfig) error // 设置某个cgroup在这个subsystem中的资源限制
	Apply(path string, pid int) error           // 添加进程至某个cgroup中
	Remove(path string) error
}

var (
	SubsystemsIns = []Subsystem{
		&CpusetSubSystem{},
		&MemorySubSystem{},
		&CpuSubSystem{},
	}
)
