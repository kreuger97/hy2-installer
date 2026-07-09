package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kreuger97/hy2-installer/internal/ui"
)

func main() {
	if os.Geteuid() != 0 {
		fmt.Println("This installer requires root privileges. Please run with sudo.")
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	if err := ui.RunWithSignals(sigCh); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
