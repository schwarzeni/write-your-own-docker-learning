package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/schwarzeni/write-your-own-docker/cpt3/cpt3-1-run/container"
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
	},
	Action: func(ctx *cli.Context) error{
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		cmd := ctx.Args().Get(0)
		tty := ctx.Bool("ti")
		Run(tty, cmd)
		return nil
	},
}

// 容器内进程初始化时的指令
var initCommand = cli.Command{
	Name: "init",
	Usage: "Do not call it outside",
	Action: func(ctx *cli.Context) error{
		log.Infof("init come on")
		cmd := ctx.Args().Get(0)
		log.Infof("command %s", cmd)
		err := container.RunContainerInitProgress(cmd, nil)
		return err
	},
}
