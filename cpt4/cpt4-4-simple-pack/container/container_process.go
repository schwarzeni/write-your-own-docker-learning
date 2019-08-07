package container

import (
	log "github.com/Sirupsen/logrus"
	"github.com/schwarzeni/write-your-own-docker/cpt4/cpt4-4-simple-pack/config"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
)

func NewParentProcess(tty bool, volumn string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Errorf("New pipe error %v", err)
		return nil, nil
	}
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
					syscall.CLONE_NEWNET |
					syscall.CLONE_NEWPID |
					syscall.CLONE_NEWNS  |
					syscall.CLONE_NEWIPC,
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	// TODO 理解这里
	cmd.ExtraFiles = []*os.File{readPipe}

	// TODO: change working dir
	mntURL := path.Join(config.BASE_URL, "mnt")
	rootURL := path.Join(config.BASE_URL)
	NewWorkspace(rootURL, mntURL, volumn)
	cmd.Dir = mntURL

	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}

// 用于创建容器的文件系统
func NewWorkspace(rootURL string, mntURL string, volume string) {
	CreateReadOnlyLayer(rootURL)
	CreateWriteLayer(rootURL)
	CreateMountPoint(rootURL, mntURL)

	if volume != "" {
		volumeURLs := volumeUrlExtract(volume)
		if len(volumeURLs) == 2  && volumeURLs[0] != "" && volumeURLs[1] != "" {
			MountVolume(rootURL, mntURL, volumeURLs)
			log.Infof("mount volume: %q", volumeURLs)
		} else {
			log.Infof("Volume parameter input is not correct")
		}
	}

}

// 创建busybox文件夹，将busybox.tar解压到其目录下
// 作为容器的只读层
func CreateReadOnlyLayer(rootURL string) {
	busyboxURL := path.Join(rootURL, "busybox/")
	busyboxTarURL := path.Join(rootURL, "busybox.tar")

	exist, err := PathExists(busyboxURL)
	if err != nil {
		log.Infof("Fail to judge whether dir %s exists. %v", err)
	}

	if !exist {
		if err := os.Mkdir(busyboxURL, 0777); err != nil {
			log.Errorf("Mkdir dir %s for busybox error, %v", busyboxURL, err)
		}

		// 解压
		if _, err := exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
			log.Errorf("Untar dir %s error. %v", busyboxTarURL, err)
		}
	}
}

// 创建名为writeLayer的文件夹
// 作为容器唯一的可写层
func CreateWriteLayer(rootURL string) {
	writeURL := path.Join(rootURL, "writeLayer/")
	if err := os.Mkdir(writeURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s for writelayer error. %v", writeURL, err)
	}
}

// 首先创建mnt文件夹作为挂载点
// 然后把writeLayer目录和busybox目录mount到mnt目录下
func CreateMountPoint(rootURL string, mntURL string) {
	if err := os.Mkdir(mntURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s for mnt error. %v", mntURL, err)
	}
	// TODO: 体会一下
	dirs := "dirs=" + path.Join(rootURL, "writeLayer") + ":" + path.Join(rootURL, "busybox")
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("error when mount aufs. %v", err)
	}
}

// 在容器退出时 unmount mnt目录，删除mnt和writeLayer文件夹
func DeleteWorkspace(rootURL string, mntURL string, volume string) {
	volumeURLs := volumeUrlExtract(volume)
	if len(volumeURLs) == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
		DeleteMountPointWithVolume(rootURL, mntURL, volumeURLs)
	} else {
		DeleteMountPoint(mntURL)
	}
	DeleteWriteLayer(rootURL)
}

// unmount mnt
// delete mnt
func DeleteMountPoint(mntURL string) {
	cmd := exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("error when umount mnt %s. %v", mntURL, err)
	}
	if err := os.RemoveAll(mntURL); err != nil {
		log.Errorf("error when remove mnt folder %s. %v", mntURL, err)
	}
}

// delete writeLayer文件夹
func DeleteWriteLayer(rootURL string) {
	writeURL := path.Join(rootURL, "writeLayer")
	if err := os.RemoveAll(writeURL); err != nil {
		log.Errorf("error when remove writeLayer folder %s. %v", writeURL, err)
	}
}

// 判断文件路径是否存在
func PathExists(sourceURL string) (exist bool, err error) {
	_, err = os.Stat(sourceURL)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}



// 挂载 volume
func MountVolume(rootURL string, mntURL string, volumesURLs []string) {

	// 创建宿主机文件目录
	parentURL := volumesURLs[0]
	if err := os.Mkdir(parentURL, 0777); err != nil {
		log.Infof("Mkdir parent dir %s error. %v", parentURL, err)
	}

	// 在容器文件系统中创建挂载点
	containerURL := volumesURLs[1]
	containerVolumeURL := path.Join(mntURL, containerURL)
	if err := os.Mkdir(containerVolumeURL, 0777); err != nil {
		log.Infof("Mkdir container dir %s error. %v", containerVolumeURL, err)
	}

	// 将宿主机文件目录挂载到容器挂载点
	dirs := "dirs=" + parentURL
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumeURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Mount volume failed %v", err)
	}
}

// 删除 volume
func DeleteMountPointWithVolume(rootURL string, mntURL string, volumeURLs []string) {
	//  umount 掉
	containerURL := path.Join(mntURL, volumeURLs[1])
	cmd := exec.Command("umount", containerURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Unmount volumn failed. %v", err)
	}

	// 执行剩余的操作
	// umount掉整个容器的挂载点
	// 删除其挂载点
	DeleteMountPoint(mntURL)
}


// 解析用户传入的volume
func volumeUrlExtract(volume string) ([]string) {
	var volumeURLS []string
	volumeURLS = strings.Split(volume, ":")
	return volumeURLS
}

