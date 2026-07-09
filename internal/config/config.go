package config

import (
	"net/url"
	"strings"
)

// ParseMasqueradeHost extracts the hostname from a masquerade URL for use as SNI.
// Falls back to "hy2-server" if URL is empty or unparseable.
func ParseMasqueradeHost(masqueradeURL string) string {
	if masqueradeURL == "" {
		return "hy2-server"
	}
	u, err := url.Parse(masqueradeURL)
	if err != nil {
		return "hy2-server"
	}
	host := u.Hostname()
	if host == "" {
		return "hy2-server"
	}
	// Strip "www." prefix for cleaner SNI
	return strings.TrimPrefix(host, "www.")
}

type Config struct {
	Port     string
	Password string
	CertPath string
	KeyPath  string

	MasqueradeEnabled  bool
	MasqueradeURL      string
	MasqueradeHTTPPort string
	MasqueradeHTTPSPort string
	MasqueradeForceHTTPS bool

	BandwidthUp   string
	BandwidthDown string
}

func Default() Config {
	return Config{
		Port:                "443",
		CertPath:            "/etc/ssl/hysteria/cert.pem",
		KeyPath:             "/etc/ssl/hysteria/key.pem",
		MasqueradeURL:       "https://www.bing.com",
		MasqueradeHTTPPort:  ":80",
		MasqueradeHTTPSPort: ":443",
		MasqueradeForceHTTPS: true,
		BandwidthUp:         "30 mbps",
		BandwidthDown:       "80 mbps",
	}
}
