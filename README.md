# hy2tool

Hysteria 2 server installer and manager.

## Install

```bash
curl -fsSL -o hy2tool https://github.com/kreuger97/hy2-installer/releases/latest/download/hy2tool_linux_amd64
chmod +x hy2tool
sudo ./hy2tool install
```

## Commands

| Command | Description |
|---------|-------------|
| `hy2tool install` | Interactive installer |
| `hy2tool link` | Show connection link + QR code |
| `hy2tool status` | Service status |
| `hy2tool start` | Start service |
| `hy2tool stop` | Stop service |
| `hy2tool restart` | Restart service |

## Features

- Self-signed SSL certificate
- Masquerade proxy (optional)
- Bandwidth limits
- Firewall configuration (ufw)
- Random password generation
