#!/bin/bash

# Define variables
BIN_PATH="/usr/local/bin/"
CONFIG_DIR="/etc/wallpaper"

# Create necessary directories
echo "Creating installation directories..."
sudo mkdir -p "$CONFIG_DIR"

# Download the Go binary and config file
echo "Downloading Go binary and config file..."
wget https://github.com/Toymakerftw/wallpaper/raw/refs/heads/main/wallpaper -O "$BIN_PATH/wallpaper"
wget https://github.com/Toymakerftw/wallpaper/raw/refs/heads/main/config.json -O "$CONFIG_DIR/config.json"

# Make the binary executable
sudo chmod +x /usr/local/bin/wallpaper

# Create a systemd service file to run the app
SERVICE_FILE="/etc/systemd/system/wallpaper.service"

echo "Creating systemd service..."
sudo bash -c "cat > $SERVICE_FILE" <<EOF
[Unit]
Description=Wallpaper Update Service
After=network.target

[Service]
ExecStart=/usr/local/bin/wallpaper
WorkingDirectory=/etc/wallpaper
Restart=always
User=root
Group=root

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd, enable and start the service
sudo systemctl daemon-reload
sudo systemctl enable wallpaper.service
sudo systemctl start wallpaper.service

echo "Wallpaper service has been installed and started."

# Check if the service is running
sudo systemctl status wallpaper.service
