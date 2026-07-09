package main

import (
	"fmt"
	"os"

	"github.com/kreuger97/hy2-installer/internal/ui"
)

func main() {
	if os.Geteuid() != 0 {
		fmt.Println("This installer requires root privileges. Please run with sudo.")
		os.Exit(1)
	}

	if err := ui.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
