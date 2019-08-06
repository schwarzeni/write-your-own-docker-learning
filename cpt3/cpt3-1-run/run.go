package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/schwarzeni/write-your-own-docker/cpt3/cpt3-1-run/container"
	"os"
)

// clone一个Namespace隔离的进程
// 在子进程中调用 /proc/self/exe，调用自己，执行init操作，初始化容器的资源
func Run(tty bool, cmd string)  {
	parent := container.NewParentProcess(tty, cmd)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	if err := parent.Wait(); err != nil {
		log.Fatal(err)
	}
	os.Exit(-1)
}

