# Cloudflare-DDNS
This is a cloudflare ddns integration.

## Usage

This will only work if you host your dns with cloudflare.

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

#### SystemD Service
```bash
[unit]
Description=Cloudflare DDNS
After=network.target

[service]
Type=simple
User=cfddns
Group=cfddns
Environment=CF_TOKEN=your_token
Environment=RECORD_NAME=foo.bar.example.com
Environment=INTERVAL=5m
Environment=PROXIED=1
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
