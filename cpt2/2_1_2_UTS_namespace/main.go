// p 11
// 模拟出UTS Namespace
// 拥有独立的hostname
// 利用 `pstree -pl` 查看进程列表 为首的进程PID为1
// $$ 当前PID 7906  -->  sudo readlink /proc/7906/ns/uts -->  uts[4026532246]
// $$ 父进程PID 7902 --> sudo readlink /proc/7902/ns/uts -->  uts[4026531838]
// 修改hostname -b nzy --> 外部不受影响
package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
