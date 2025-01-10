#!/bin/bash
echo "Installing cloudflare-ddns..."

# Get latest release
echo "Getting latest release..."
latest=$(curl https://api.github.com/repos/deltxprt/cloudflare-ddns/tags | jq -r '.[].name' | head -n1)
arch=$(uname -m)

# Download latest release
echo "Downloading latest release..."
wget -O /tmp/cloudflare-ddns.tar.gz "https://github.com/deltxprt/cloudflare-ddns/releases/download/${latest}/cloudflare-ddns_Linux_${arch}.tar.gz"

# Extract files
echo "Extracting files..."
sudo tar -xzf /tmp/cloudflare-ddns.tar.gz -C /usr/local/bin

echo "Cleaning up..."
sudo rm /tmp/cloudflare-ddns.tar.gz

# create service account linux
echo "Creating service account..."
sudo useradd -m -s /bin/bash -d /home/cfddns cfddns

# set permissions
echo "Setting permissions..."
sudo chown cfddns:cfddns /usr/local/bin/cloudflare-ddns

# create config file
echo "Creating config file..."
sudo mkdir /etc/cfddns

sudo touch /home/cfddns/.env

sudo tee -a /etc/cfddns/.env <<EOF
CF_TOKEN=REPLACEME
RECORD_NAME=foo.bar.example.com
INTERVAL=5m
PROXIED=1
EOF

sudo chown cfddns:cfddns -R /home/cfddns

# create systemd service
echo "Creating systemd service..."
wget -O /etc/systemd/system/cf-ddns.service "https://raw.githubusercontent.com/deltxprt/cloudflare-ddns/refs/tags/${latest}/cf-ddns.service"

echo "Installation complete. Please edit /etc/cfddns/.env with your Cloudflare API token and desired settings."
echo "Then run 'systemctl enable cf-ddns && systemctl start cf-ddns' to start the service."