#!/usr/bin/env bash

#ls /var/lib/docker/aufs/diff
# image layer

# ls /var/lib/docker/aufs/layers
# 如何堆栈layer的metadata

# 实验准备工作
mkdir mnt container-layer
touch container-layer/container-layer.txt
echo "I am container layer" >> container-layer/container-layer.txt

for i in {1..4}
do
    dirname=image-layer${i}
    filename=${dirname}/image-layer${i}.txt

    mkdir $dirname
    touch ${filename}
    echo "I am image layer ${i}" >> ${filename}
done


# 进行挂载
# sudo mount -t aufs -o dirs=./container-layer:./image-layer4:./image-layer3:./image-layer2:./image-layer1 none ./mnt


#  cat /sys/fs/aufs/si_a236c06529c6ac94/* | less
# /home/parallels/Workspace/go/write-your-own-docker/cpt2/2_3_ufs/container-layer=rw
# /home/parallels/Workspace/go/write-your-own-docker/cpt2/2_3_ufs/image-layer4=ro
# /home/parallels/Workspace/go/write-your-own-docker/cpt2/2_3_ufs/image-layer3=ro
# /home/parallels/Workspace/go/write-your-own-docker/cpt2/2_3_ufs/image-layer2=ro
# /home/parallels/Workspace/go/write-your-own-docker/cpt2/2_3_ufs/image-layer1=ro
