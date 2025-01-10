# Cloudflare-DDNS
This is a cloudflare ddns integration.

## Usage

This will only work if you host your dns with cloudflare.

The app will run on an interval that you set or default to 5 minutes. 
It will check your public ip via the url 'https://cloudflare.com/cdn-cgi/trace' and update the record if it has changed.

## Configurations

### Cloudflare Token
You will need to create a token with the following permissions:
- DNS: Edit

### Environment Variables
- `CF_TOKEN` - Your cloudflare token
- `RECORD_NAME` - The record you want to update (foo.bar.example.com)
- `INTERVAL` - The interval you want to update the record (5m, 1h, 1d) default: 5m
- `PROXIED` - If you want to proxy the record (1 or 0) default: 0

## Installation

### Linux
quick install
```bash
curl -s https://raw.githubusercontent.com/deltxprt/cloudflare-ddns/refs/heads/master/install.sh | bash
```

#### SystemD Service
```bash
[unit]
Description=Cloudflare DDNS
After=network.target

[service]
Type=simple
User=cfddns
EnvironmentFile=/etc/cfddns/.env
ExecStart=/usr/bin/cloudflare-ddns

[install]
WantedBy=multi-user.target

```
soon

### Windows

soon

### Container

#### Podman Adhoc
```bash
podman run -d \
    --name=cloudflare-ddns \
    --restart unless-stopped \
    -e CF_TOKEN=your_token \
    -e RECORD_NAME=foo.bar.example.com \
    -e INTERVAL=5m \
    -e PROXIED=1 \
    ghcr.io/robertscherzer/cloudflare-ddns:latest
```

#### Podman Quadlet

soon

#### Docker Adhoc
```bash
docker run -d \
    --name=cloudflare-ddns \
    --restart unless-stopped \
    -e CF_TOKEN=your_token \
    -e RECORD_NAME=foo.bar.example.com \
    -e INTERVAL=5m \
    -e PROXIED=1 \
    ghcr.io/robertscherzer/cloudflare-ddns:latest
```

#### Docker Compose
```yaml
services:
    cloudflare-ddns:
        image: ghcr.io/robertscherzer/cloudflare-ddns:latest
        container_name: cloudflare-ddns
        restart: unless-stopped
        environment:
          - CF_TOKEN=your_token
          - RECORD_NAME=foo.bar.example.com
          - INTERVAL=5m
          - PROXIED=1
```

# Troubleshooting

## Logs
You can view the logs by running `journalctl -u cloudflare-ddns`

## Issues
If you have any issues, please open an issue on the github repo.