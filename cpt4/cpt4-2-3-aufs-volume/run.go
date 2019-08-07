package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/schwarzeni/write-your-own-docker/cpt4/cpt4-2-3-aufs-volume/cgroup"
	"github.com/schwarzeni/write-your-own-docker/cpt4/cpt4-2-3-aufs-volume/cgroup/subsystems"
	"github.com/schwarzeni/write-your-own-docker/cpt4/cpt4-2-3-aufs-volume/config"
	"github.com/schwarzeni/write-your-own-docker/cpt4/cpt4-2-3-aufs-volume/container"
	"os"
	"path"
	"strings"
)

// clone一个Namespace隔离的进程
// 在子进程中调用 /proc/self/exe，调用自己，执行init操作，初始化容器的资源
func Run(tty bool, cmdArr []string, res *subsystems.ResourceConfig, volumn string)  {
	parent, writePipe := container.NewParentProcess(tty, volumn)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	// 设置初始化参数
	cgroupManager := cgroup.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(cmdArr, writePipe)

	if err := parent.Wait(); err != nil {
		log.Fatal(err)
	}

	mntURL := path.Join(config.BASE_URL,  "mnt")
	rootURL := path.Join(config.BASE_URL)
	container.DeleteWorkspace(rootURL, mntURL, volumn)
	os.Exit(-1)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	// TODO 修复此处
	writePipe.WriteString(command)
	writePipe.Close()
}

