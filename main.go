package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func run() {
	fmt.Printf("pid: %d\n", os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
		//		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWUSER,
		//		UidMappings: []syscall.SysProcIDMap{{
		//			ContainerID: 1000,
		//			HostID:      0,
		//			Size:        1,
		//		}},
		//		GidMappings: []syscall.SysProcIDMap{{
		//			ContainerID: 1001,
		//			HostID:      0,
		//			Size:        1,
		//		}},
	}

	must(cmd.Run())
}

func child() {
	fmt.Printf("Running %v with arguments %v as %d\n", os.Args[2], os.Args[3:], os.Getpid())

	syscall.Chroot("/home/btoll/projects/containers-from-scratch/rootfs")
	os.Chdir("/")

	syscall.Sethostname([]byte("container"))

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("Unrecognized command")
	}
}
