#!/bin/bash

echo "Uninstalling cloudflare-ddns..."

# remove files

echo "Removing files..."
sudo rm /usr/local/bin/cloudflare-ddns
sudo rm /etc/systemd/system/cf-ddns.service
sudo rm -rdf /etc/cfddns

# remove service account
echo "Removing service account..."
sudo deluser cfddns

echo "Uninstallation complete."