package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/schwarzeni/write-your-own-docker/cpt4/cpt4-2-3-aufs-volume/cgroup/subsystems"
	"github.com/schwarzeni/write-your-own-docker/cpt4/cpt4-2-3-aufs-volume/container"
	"github.com/urfave/cli"
	)

// 父进程的运行指令
var runCommand = cli.Command{
	Name: "run",
	Usage: "mydocker run -ti [command]",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name: "ti",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name: "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name: "cpushare",
			Usage: "cpushare limit",
		},
		cli.StringFlag{
			Name: "cpuset",
			Usage: "cpuset limit",
		},
		cli.StringFlag{
			Name: "v",
			Usage: "volume",
		},
	},
	Action: func(ctx *cli.Context) error{
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}

		var cmdArr []string
		for _, arg := range ctx.Args() {
			cmdArr = append(cmdArr, arg)
		}

		tty := ctx.Bool("ti")
		resConf := &subsystems.ResourceConfig{
			MemoryLimit: ctx.String("m"),
			CpuSet: ctx.String("cpuset"),
			CpuShare: ctx.String("cpushare"),
		}
		volumn := ctx.String("v")
		Run(tty, cmdArr, resConf, volumn)
		return nil
	},
}

// 容器内进程初始化时的指令
var initCommand = cli.Command{
	Name: "init",
	Usage: "Do not call it outside",
	Action: func(ctx *cli.Context) error{
		log.Infof("init come on")
		err := container.RunContainerInitProgress()
		return err
	},
}
