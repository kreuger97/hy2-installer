package config

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
