package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"

	qrcode "github.com/skip2/go-qrcode"

	"github.com/kreuger97/hy2-installer/internal/install"
)

const usage = `hy2ctl - Hysteria 2 Server Manager

Usage:
  hy2ctl <command>

Commands:
  link      Show connection link and QR code
  status    Show service status
  start     Start the service
  stop      Stop the service
  restart   Restart the service
`

func main() {
	if os.Geteuid() != 0 {
		fmt.Fprintln(os.Stderr, "Error: requires root privileges")
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Print(usage)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "link":
		cmdLink()
	case "status":
		cmdStatus()
	case "start":
		cmdService("start")
	case "stop":
		cmdService("stop")
	case "restart":
		cmdService("restart")
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		fmt.Print(usage)
		os.Exit(1)
	}
}

func parseConfig() (port, password string) {
	data, err := os.ReadFile("/etc/hysteria/config.yaml")
	if err != nil {
		return "443", ""
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "listen:") {
			val := strings.TrimSpace(strings.TrimPrefix(line, "listen:"))
			val = strings.TrimPrefix(val, ":")
			if val != "" {
				port = val
			}
		}
		if strings.HasPrefix(line, "password:") {
			val := strings.TrimSpace(strings.TrimPrefix(line, "password:"))
			val = strings.Trim(val, "\"")
			if val != "" {
				password = val
			}
		}
	}

	if port == "" {
		port = "443"
	}
	return
}

func getOutboundIP() string {
	providers := []string{
		"curl -fsSL --connect-timeout 3 https://ipinfo.io/ip",
		"curl -fsSL --connect-timeout 3 https://api.ipify.org",
	}
	for _, p := range providers {
		out, err := exec.Command("bash", "-c", p).Output()
		if err == nil {
			ip := strings.TrimSpace(string(out))
			if ip != "" {
				return ip
			}
		}
	}
	out, err := exec.Command("hostname", "-I").Output()
	if err == nil {
		parts := strings.Fields(string(out))
		if len(parts) > 0 {
			return parts[0]
		}
	}
	return "<unknown>"
}

func buildLink(port, password, ip string) string {
	u := url.URL{
		Scheme: "hysteria2",
		User:   url.User(password),
		Host:   fmt.Sprintf("%s:%s", ip, port),
		RawQuery: url.Values{
			"insecure": {"1"},
			"alpn":     {"h3"},
			"sni":      {"hy2-server"},
		}.Encode(),
	}
	return u.String()
}

func cmdLink() {
	port, password := parseConfig()
	ip := getOutboundIP()
	link := buildLink(port, password, ip)

	fmt.Println()
	fmt.Println("  Connection Link")
	fmt.Println("  ─────────────────────────────────────────")
	fmt.Println()
	fmt.Printf("  %s\n", link)
	fmt.Println()

	qr, err := qrcode.New(link, qrcode.Medium)
	if err == nil {
		fmt.Println("  QR Code")
		fmt.Println("  ─────────────────────────────────────────")
		fmt.Println()
		for _, line := range strings.Split(qr.ToString(true), "\n") {
			fmt.Println("  " + line)
		}
	}
	fmt.Println()
}

func cmdStatus() {
	if !install.HysteriaInstalled() {
		fmt.Println("\n  Hysteria 2 is not installed.")
		os.Exit(1)
	}

	cmd := exec.Command("systemctl", "is-active", "hysteria-server")
	out, _ := cmd.Output()
	status := strings.TrimSpace(string(out))

	cmd = exec.Command("systemctl", "is-enabled", "hysteria-server")
	enOut, _ := cmd.Output()
	enabled := strings.TrimSpace(string(enOut))

	cmd = exec.Command("systemctl", "show", "hysteria-server", "--property=MainPID", "--value")
	pidOut, _ := cmd.Output()
	pid := strings.TrimSpace(string(pidOut))

	port, password := parseConfig()
	ip := getOutboundIP()

	fmt.Println()
	fmt.Println("  Hysteria 2 Server Status")
	fmt.Println("  ─────────────────────────────────────────")
	fmt.Println()

	if status == "active" {
		fmt.Println("  Status:   \033[32m● active\033[0m")
	} else {
		fmt.Printf("  Status:   \033[31m● %s\033[0m\n", status)
	}

	if enabled == "enabled" {
		fmt.Println("  Enabled:  yes")
	} else {
		fmt.Println("  Enabled:  no")
	}

	fmt.Printf("  PID:      %s\n", pid)
	fmt.Printf("  IP:       %s\n", ip)
	fmt.Printf("  Port:     %s\n", port)

	configPath := "/etc/hysteria/config.yaml"
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("  Config:   %s\n", configPath)
	} else {
		fmt.Printf("  Config:   not found\n")
	}

	uri := buildLink(port, password, ip)
	fmt.Printf("  Link:     %s\n", uri)
	fmt.Println()
}

func cmdService(action string) {
	if !install.HysteriaInstalled() {
		fmt.Println("\n  Hysteria 2 is not installed.")
		os.Exit(1)
	}

	cmd := exec.Command("systemctl", action, "hysteria-server")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "\n  Failed to %s service: %v\n\n", action, err)
		os.Exit(1)
	}
	fmt.Printf("\n  Service %sed successfully.\n\n", action)
}
