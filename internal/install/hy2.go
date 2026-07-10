package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kreuger97/hy2tool/internal/config"
)

func HysteriaInstalled() bool {
	_, err := exec.LookPath("hysteria")
	return err == nil
}

func InstallHysteria() error {
	cmd := exec.Command("bash", "-c", "bash <(curl -fsSL https://get.hy2.sh/)")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func GenerateCert() error {
	dir := "/etc/ssl/hysteria"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	cmd := exec.Command("openssl", "req", "-x509", "-nodes", "-newkey", "rsa:4096",
		"-keyout", filepath.Join(dir, "key.pem"),
		"-out", filepath.Join(dir, "cert.pem"),
		"-days", "3650",
		"-subj", "/CN=hy2-server",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("openssl: %w", err)
	}

	for _, f := range []string{"key.pem", "cert.pem"} {
		if err := exec.Command("chmod", "644", filepath.Join(dir, f)).Run(); err != nil {
			return fmt.Errorf("chmod %s: %w", f, err)
		}
	}
	return nil
}

func WriteConfig(cfg config.Config) error {
	dir := "/etc/hysteria"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	yaml := fmt.Sprintf(`listen: :%s

auth:
  type: password
  password: %s

tls:
  cert: %s
  key: %s
`, cfg.Port, cfg.Password, cfg.CertPath, cfg.KeyPath)

	if cfg.MasqueradeEnabled {
		forceHTTPS := "false"
		if cfg.MasqueradeForceHTTPS {
			forceHTTPS = "true"
		}
		yaml += fmt.Sprintf(`
masquerade:
  type: proxy
  proxy:
    url: %s
    rewriteHost: true
  listenHTTP: %s
  listenHTTPS: %s
  forceHTTPS: %s
`, cfg.MasqueradeURL, cfg.MasqueradeHTTPPort, cfg.MasqueradeHTTPSPort, forceHTTPS)
	}

	if cfg.BandwidthUp != "" || cfg.BandwidthDown != "" {
		yaml += "\nbandwidth:\n"
		if cfg.BandwidthUp != "" {
			yaml += fmt.Sprintf("  up: %s\n", cfg.BandwidthUp)
		}
		if cfg.BandwidthDown != "" {
			yaml += fmt.Sprintf("  down: %s\n", cfg.BandwidthDown)
		}
	}

	return os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(yaml), 0644)
}

func StartService() error {
	cmds := []string{
		"systemctl daemon-reload",
		"systemctl start hysteria-server",
		"systemctl enable hysteria-server",
	}
	for _, c := range cmds {
		cmd := exec.Command("bash", "-c", c)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s: %w", c, err)
		}
	}
	return nil
}

func ConfigureFirewall(cfg config.Config) error {
	cmd := exec.Command("ufw", "allow", cfg.Port)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	if cfg.MasqueradeEnabled {
		if cfg.MasqueradeHTTPPort != "" {
			port := strings.TrimPrefix(cfg.MasqueradeHTTPPort, ":")
			cmd = exec.Command("ufw", "allow", port, "tcp")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return err
			}
		}
		if cfg.MasqueradeHTTPSPort != "" {
			port := strings.TrimPrefix(cfg.MasqueradeHTTPSPort, ":")
			cmd = exec.Command("ufw", "allow", port, "tcp")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return err
			}
		}
	}
	return nil
}
