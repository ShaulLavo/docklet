package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Subcommand is required")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		runCmd := flag.NewFlagSet("run", flag.ExitOnError)
		runCmd.Parse(os.Args[2:])
		run(runCmd.Args())
	case "build":
		buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
		tag := buildCmd.String("tag", "latest", "Name of the container image")
		path := buildCmd.String("path", ".", "Path to the Dockletfile")
		buildCmd.Parse(os.Args[2:])
		build(*tag, *path)
	case "child":
		runCmd := flag.NewFlagSet("child", flag.ExitOnError)
		runCmd.Parse(os.Args[2:])
		child(runCmd.Args())
	default:
		fmt.Printf("Invalid subcommand: %s\n", os.Args[1])
	}
}

func run(args []string) {
	fmt.Printf("Running command: %v\n", args)

	// Setting up a command to re-execute itself but with "child" subcommand
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, args...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	const CLONE_NEWPID = syscall.CLONE_NEWPID // New PID namespace
	const CLONE_NEWUTS = syscall.CLONE_NEWUTS // New UTS namespace (hostname)
	cmd.SysProcAttr = &syscall.SysProcAttr{Cloneflags: CLONE_NEWUTS | CLONE_NEWPID}

	handleError(cmd.Run(), "Error running the command")
}

func build(tag, path string) {
	fmt.Printf("Building container image with tag: %s from path: %s\n", tag, path)
	//TODO Add your container build logic here.
}

func child(args []string) {
	fmt.Printf("Running in container namespace with arguments: %v\n", args)

	// Command to execute within the new namespace
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set hostname for the new UTS namespace
	handleError(syscall.Sethostname([]byte("docklet")), "Error setting hostname")

	// Change root to the new filesystem (chroot)
	handleError(syscall.Chroot("ubuntu-rootfs/"), "Error changing root")
	handleError(syscall.Chdir("/"), "Error changing directory")

	// Mount proc filesystem for process information
	handleError(syscall.Mount("proc", "proc", "proc", 0, ""), "Error mounting proc")

	// Run the command
	handleError(cmd.Run(), "Error running the command in child namespace")

	// Unmount proc filesystem before exiting
	handleError(syscall.Unmount("/proc", 0), "Error unmounting /proc")
}

func handleError(err error, message string) {
	if err != nil {
		fmt.Printf("%s: %v\n", message, err)
		os.Exit(1)
	}
}
