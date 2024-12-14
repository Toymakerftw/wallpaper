#!/bin/bash

# Define variables
BIN_PATH="/usr/local/bin/"
CONFIG_DIR="/etc/wallpaper"
SERVICE_FILE="/etc/systemd/system/wallpaper.service"

# Stop and disable the systemd service
echo "Stopping and disabling wallpaper service..."
sudo systemctl stop wallpaper.service
sudo systemctl disable wallpaper.service

# Remove the systemd service file
echo "Removing systemd service file..."
sudo rm -f "$SERVICE_FILE"

# Remove the binary and configuration files
echo "Removing binary and config files..."
sudo rm -f "$BIN_PATH/wallpaper"
sudo rm -f "$CONFIG_DIR/config.json"

# Remove the wallpaper configuration directory
echo "Removing configuration directory..."
sudo rmdir "$CONFIG_DIR"

# Reload systemd daemon to apply changes
sudo systemctl daemon-reload

echo "Uninstallation complete."

# Check if the service is still installed
sudo systemctl status wallpaper.service
