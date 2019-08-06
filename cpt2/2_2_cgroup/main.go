package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"syscall"
)

const cgroupMemoryHierarchMount = "/sys/fs/cgroup/memory"

func main() {
	if os.Args[0] == "/proc/self/exe" {
		log.Printf("current pid %d\n", syscall.Getpid())
		cmd := exec.Command("sh", "-c", `stress --vm-bytes 200m --vm-keep -m 1`)
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal("stress has error", err)
		}
	}

	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("%v\n", cmd.Process.Pid)
		cgName := "testmemorylimit"
		var err error

		// 区分fork
		if err = os.Mkdir(path.Join(cgroupMemoryHierarchMount, cgName), 0755); err != nil {
			if !strings.Contains(fmt.Sprintln(err), "file exists") {
				log.Fatal("mkdir error", err)
			} else {
				log.Println("mkdir warning ", err)
			}
		}
		if err = ioutil.WriteFile(path.Join(cgroupMemoryHierarchMount, cgName, "tasks"),
			[]byte(strconv.Itoa(cmd.Process.Pid)), 0644); err != nil {
			log.Fatal("write to tasks error", err)
		}

		//bs := make([]byte, 4)
		//binary.LittleEndian.PutUint32(bs, 1)
		// TODO: 老是写不进去
		fspath := path.Join(cgroupMemoryHierarchMount, cgName, "memory.oom_control")
		if err= exec.Command("echo", "1 >>"+ fspath).Run(); err != nil {
			log.Fatal("write to memory.oom_control error", err)
		}


		if err = ioutil.WriteFile(path.Join(cgroupMemoryHierarchMount, cgName, "memory.limit_in_bytes"),
			[]byte("100m"), 0644); err != nil {
			log.Fatal("write to memory.limit_in_bytes error", err)
		}
	}
	if _, err := cmd.Process.Wait(); err != nil {
		log.Fatal("cmd.Process.Wait error", err)
	}
}
