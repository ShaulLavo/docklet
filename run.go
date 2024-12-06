package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

type limits struct {
	Limits []limit
}

type limit struct {
	name  string
	path  string
	param string
	value []byte
}

func run(args []string) {
	fmt.Printf("Running %v \n", args)

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	//cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	err := cmd.Run()
	if err != nil {
		fmt.Println("error running child:", err)
		panic(err)
	}
}

func child(args []string) {
	fmt.Printf("Running from proc in namespace %v \n", args)

	err := cgroup()
	if err != nil {
		fmt.Println("error setting cgroup:", err)
		panic(err)
	}

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = syscall.Sethostname([]byte("docklet"))
	if err != nil {
		fmt.Println("error setting hostname:", err)
		panic(err)
	}

	if err = syscall.Chroot("ubuntu-rootfs/"); err != nil {
		fmt.Println("error changing root:", err)
		panic(err)
	}

	if err = syscall.Chdir("/"); err != nil {
		fmt.Println("error changing directory:", err)
		panic(err)
	}

	if err = syscall.Mount("proc", "proc", "proc", 0, ""); err != nil {
		fmt.Println("error mounting proc:", err)
		panic(err)
	}

	err = cmd.Run()
	if err != nil {
		fmt.Println("error running proc:", err)

		panic(err)
	}

	if err = syscall.Unmount("proc", 0); err != nil {
		fmt.Println("error unmounting proc:", err)
		panic(err)
	}
}

func cgroup() error {

	cgrouplimits := limits{
		[]limit{
			{
				"pids",
				"/sys/fs/cgroup/pids/docklet",
				"pids.max",
				[]byte("20"),
			},
			{
				"memory",
				"/sys/fs/cgroup/memory/docklet",
				"memory.limit_in_bytes",
				[]byte("1000000"),
			},
			{
				"cpu",
				"/sys/fs/cgroup/cpu/docklet",
				"cpu.shares",
				[]byte("512"),
			},
		},
	}

	for _, l := range cgrouplimits.Limits {
		fmt.Println(filepath.Join(l.path, l.param))
		//os.Mkdir
		os.Mkdir(l.path, 0755)
		//Create cgroup limit
		err := os.WriteFile(filepath.Join(l.path, l.param), l.value, 0700)
		if err != nil {
			return err
		}
		//ioutil.WriteFile (add proc to cgroup)
		err = os.WriteFile(filepath.Join(l.path, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)
		if err != nil {
			return err
		}
		//ioutil.WriteFile notify_on-release
		err = os.WriteFile(filepath.Join(l.path, "notify_on_release"), []byte("1"), 0700)
		if err != nil {
			return err
		}
	}

	//pids := filepath.Join(cgroup, "pids")
	//os.Mkdir(filepath.Join(pids, "docklet"), 0755)

	//err := ioutil.WriteFile(filepath.Join(pids, "docklet/pids.max"), []byte("20"), 0700)
	//if err != nil {
	//	return err
	//}

	//err = ioutil.WriteFile(filepath.Join(pids, "docklet/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)
	//if err != nil {
	//	return err
	//}

	//mem := filepath.Join(cgroup, "memory")
	//os.Mkdir(filepath.Join(mem, "docklet"), 0755)

	//err = ioutil.WriteFile(filepath.Join(mem, "docklet/memory.limit_in_bytes"), []byte("1000000"), 0700)
	//if err != nil {
	//	return err
	//}

	//err = ioutil.WriteFile(filepath.Join(mem, "docklet/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)
	//if err != nil {
	//	return nil
	//}

	//cpu := filepath.Join(cgroup, "cpu")
	//os.Mkdir(filepath.Join(cpu, "docklet"), 0755)

	//err = ioutil.WriteFile(filepath.Join(cpu, "docklet/cpu.shares"), []byte("512"), 0700)
	//if err != nil {
	//	return err
	//}

	//err = ioutil.WriteFile(filepath.Join(cpu, "docklet/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)
	//if err != nil {
	//	return err
	//}

	//err = ioutil.WriteFile(filepath.Join(pids, "docklet/notify_on_release"), []byte("1"), 0700)
	//if err != nil {
	//	return err
	//}

	//err = ioutil.WriteFile(filepath.Join(mem, "docklet/notify_on_release"), []byte("1"), 0700)
	//if err != nil {
	//	return err
	//}

	//err = ioutil.WriteFile(filepath.Join(cpu, "docklet/notify_on_release"), []byte("1"), 0700)
	//if err != nil {
	//	return err
	//}

	return nil
}
