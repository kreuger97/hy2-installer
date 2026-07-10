package install

import (
	"os"
	"strings"
	"testing"

	"github.com/kreuger97/hy2tool/internal/config"
)

func TestGenerateCert(t *testing.T) {
	if err := GenerateCert(); err != nil {
		t.Fatalf("GenerateCert failed: %v", err)
	}
	for _, f := range []string{"/etc/ssl/hysteria/cert.pem", "/etc/ssl/hysteria/key.pem"} {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Fatalf("%s not created", f)
		}
	}
	t.Log("certificates created")
}

func TestWriteConfig(t *testing.T) {
	cfg := config.Default()
	cfg.Port = "8443"
	cfg.Password = "test-password"

	if err := WriteConfig(cfg); err != nil {
		t.Fatalf("WriteConfig failed: %v", err)
	}

	data, err := os.ReadFile("/etc/hysteria/config.yaml")
	if err != nil {
		t.Fatalf("cannot read config: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, ":8443") {
		t.Error("config missing port 8443")
	}
	if !strings.Contains(content, "test-password") {
		t.Error("config missing password")
	}
	if strings.Contains(content, "masquerade:") {
		t.Error("config should not contain masquerade when disabled")
	}
	t.Logf("config:\n%s", content)
}

func TestWriteConfigWithMasquerade(t *testing.T) {
	cfg := config.Default()
	cfg.Port = "443"
	cfg.Password = "test-pass"
	cfg.MasqueradeEnabled = true
	cfg.MasqueradeURL = "https://example.com"
	cfg.MasqueradeHTTPPort = ":80"
	cfg.MasqueradeHTTPSPort = ":443"
	cfg.MasqueradeForceHTTPS = true
	cfg.BandwidthUp = "50 mbps"
	cfg.BandwidthDown = "100 mbps"

	if err := WriteConfig(cfg); err != nil {
		t.Fatalf("WriteConfig failed: %v", err)
	}

	data, err := os.ReadFile("/etc/hysteria/config.yaml")
	if err != nil {
		t.Fatalf("cannot read config: %v", err)
	}

	content := string(data)
	checks := []struct {
		name, want string
	}{
		{"masquerade block", "masquerade:"},
		{"proxy type", "type: proxy"},
		{"proxy url", "url: https://example.com"},
		{"rewriteHost", "rewriteHost: true"},
		{"listenHTTP", "listenHTTP: :80"},
		{"listenHTTPS", "listenHTTPS: :443"},
		{"forceHTTPS", "forceHTTPS: true"},
		{"bandwidth up", "up: 50 mbps"},
		{"bandwidth down", "down: 100 mbps"},
	}
	for _, c := range checks {
		if !strings.Contains(content, c.want) {
			t.Errorf("config missing %s (%s)", c.name, c.want)
		}
	}
	t.Logf("config:\n%s", content)
}
