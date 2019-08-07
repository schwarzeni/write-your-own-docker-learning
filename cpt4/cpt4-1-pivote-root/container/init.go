package container

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// 子进程需要执行的内容
func RunContainerInitProgress() (err error) {
	cmdArray := readUserCommand() // 获取参数
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("Run container get user command error, cmdArray is nil")
	}

	//log.Infof("commamd %s", cmd)

	// 进行挂载操作
	setUpMount()

	path, err := exec.LookPath(cmdArray[0]) // 寻找命令的绝对路径
	if err != nil {
		log.Errorf("Exec loop path error %v", err)
		return err
	}
	log.Infof("Find path %s", path)

	// 将当前进程的PID置为1
	// 调用了Kernel的 int execve(const char *filename, char *const argv[], char *const envp[])
	// 覆盖当前进程的镜像，数据和堆栈信息
	if err = syscall.Exec(path, cmdArray, os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return
}

// 从管道中读取参数
func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	if msg, err := ioutil.ReadAll(pipe); err!= nil {
		log.Errorf("init read pipe error %v", err)
		return nil
	} else {
		msgStr := string(msg)
		return strings.Split(msgStr, " ")
	}
}

// 设置挂载点
func setUpMount() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("Get current location error %v", err)
		return
	}
	log.Infof("Current location is %s", pwd)

	if err = syscall.Mount("", "/", "", syscall.MS_PRIVATE | syscall.MS_REC, ""); err != nil {
		return
	}

	// 改变root
	if err = pivotRoot(pwd); err != nil {
		return
	}

	// mount proc
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	if err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
		return
	}
	if err = syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755"); err != nil {
		return
	}
}

func pivotRoot(root string) (err error) {
	// 为了使当前root的 老root 和 新root 不在同一个文件系统下，需要将root重新mount一次
	// bind mount  将相同的内容换一个挂载点
	if err = syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mount rootfs to itself error: %v", err)
	}

	// 创建 rootfs/.pivot_root 存储到 old_root 中
	pivotDir := filepath.Join(root, ".pivot_root")
	if err = os.Mkdir(pivotDir, 0777); err != nil {
		return
	}

	// 将 pivot_root 挂载到新的rootfs，现在老的 old_root 是挂载在 rootfs/.pivot_root
	// 挂载点现在依然可以在mount命令中看到
	if err = syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}

	// 修改当前的工作目录到根目录
	if err = syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")

	// unmount rootfs/.pivot_root
	if err = syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}

	// 删除临时文件夹
	return os.Remove(pivotDir)
}
