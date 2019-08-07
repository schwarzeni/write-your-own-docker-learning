package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/schwarzeni/write-your-own-docker/cpt4/cpt4-4-simple-pack/config"
	"os/exec"
	"path"
)

// 路径写死的，此函数仅为演示
func commitContainer(imageName string) {
	// 将 mnt 中的内容打包
	mntURL := path.Join(config.BASE_URL, "mnt")
	imageTar := path.Join(config.BASE_URL, imageName + ".tar")
	log.Infof("compact the image %s into %s", imageName, imageTar)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntURL, ".").CombinedOutput(); err != nil {
		log.Errorf("Tar folder %s error. %v", mntURL, err)
	}
}
