package container

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"syscall"
)

// 子进程需要执行的内容
func RunContainerInitProgress(cmd string, args []string) (err error) {
	log.Infof("commamd %s", cmd)

	// TODO: 注意这里
	// https://github.com/xianlubird/mydocker/issues/41#issuecomment-478799767
	// systemd 加入linux之后, mount namespace 就变成 shared by default, 所以你必须显示
	//声明你要这个新的mount namespace独立。
	if err = syscall.Mount("", "/", "", syscall.MS_PRIVATE | syscall.MS_REC, ""); err != nil {
		return
	}

	// MS_NOEXEC 本文件系统不允许执行其他程序
	// MS_NOSUID 不允许 set-user-ID 和 set-group-ID
	// MS_NODEV  默认参数
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	if err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
		return
	}

	argv := []string{cmd}

	// 将当前进程的PID置为1
	// 调用了Kernel的 int execve(const char *filename, char *const argv[], char *const envp[])
	// 覆盖当前进程的镜像，数据和堆栈信息
	if err = syscall.Exec(cmd, argv, os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return
}
