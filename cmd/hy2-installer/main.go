package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kreuger97/hy2-installer/internal/ctl"
	"github.com/kreuger97/hy2-installer/internal/ui"
)

func main() {
	if os.Geteuid() != 0 {
		fmt.Println("This installer requires root privileges. Please run with sudo.")
		os.Exit(1)
	}

	// No args = interactive installer
	if len(os.Args) < 2 {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		if err := ui.RunWithSignals(sigCh); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Subcommands
	switch os.Args[1] {
	case "link":
		ctl.CmdLink()
	case "status":
		ctl.CmdStatus()
	case "start":
		ctl.CmdService("start")
	case "stop":
		ctl.CmdService("stop")
	case "restart":
		ctl.CmdService("restart")
	case "help", "-h", "--help":
		fmt.Print(ctl.Usage)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		fmt.Print(ctl.Usage)
		os.Exit(1)
	}
}
