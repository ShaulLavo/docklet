package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mholt/archiver/v3"
)

func main() {

	if len(os.Args) <= 1 {
		showHelp()
		return
	}

	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	switch os.Args[1] {
	case "run":
		runCmd.Parse(os.Args[2:])
		arg := runCmd.Args()
		run(arg)
	case "child":
		runCmd.Parse(os.Args[2:])
		arg := runCmd.Args()
		child(arg)
	case "build":
		buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
		tag := buildCmd.String("tag", "", "Name of docklet image")
		path := buildCmd.String("path", "", "Path to DockletFile")
		buildCmd.Parse(os.Args[2:])
		build(*tag, *path)
	case "extract":
		extractCmd := flag.NewFlagSet("extract", flag.ExitOnError)
		extractCmd.Parse(os.Args[2:])
		arg := extractCmd.Args()[0]
		//extract(arg,"ubuntu-rootfs")
		err := archiver.Unarchive(arg, "ubuntu-rootfs")
		if err != nil {
			fmt.Printf("Error extracting archive %s", arg)
			os.Exit(1)
		}
	default:
		fmt.Printf("invalid subcommand %s", os.Args[1])
		os.Exit(1)
	}

}

func showHelp() {
	fmt.Println("Usage: [command] [options]")
	fmt.Println()
	fmt.Println("Available commands:")
	fmt.Println("  run       Execute a command in a new isolated environment.")
	fmt.Println("  build     Build a container image from a DockletFile.")
	fmt.Println("  extract    Extract a .tar.gz archive into the specified destination directory, defaulting to 'ubuntu-rootfs'.")
	fmt.Println()
	fmt.Println("Use '[command] --help' for more information about a specific command.")
}
