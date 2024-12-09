#!/bin/bash

# Set the binary download URL
WALLPAPER_BIN_URL="https://github.com/Toymakerftw/wallpaper/raw/refs/heads/main/wallpaper"
WALLPAPER_BIN_PATH="$HOME/wallpaper"
SERVICE_NAME="wallpaper-service"
SYSTEMD_SERVICE_PATH="/etc/systemd/system/${SERVICE_NAME}.service"

# Step 1: Download the wallpaper binary
echo "Downloading the wallpaper service binary..."
curl -L -o "$WALLPAPER_BIN_PATH" "$WALLPAPER_BIN_URL"

# Step 2: Make the binary executable
echo "Making the wallpaper binary executable..."
chmod +x "$WALLPAPER_BIN_PATH"

# Step 3: Create the systemd service for the current user
echo "Creating systemd service..."

# Create a systemd service file to run the wallpaper binary as the current user
cat <<EOF | sudo tee "$SYSTEMD_SERVICE_PATH" > /dev/null
[Unit]
Description=Wallpaper Service
After=network.target

[Service]
ExecStart=$WALLPAPER_BIN_PATH
User=$USER
WorkingDirectory=$HOME
Environment=DISPLAY=:0
Environment=DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/$(id -u $USER)/bus
Restart=always
RestartSec=10

[Install]
WantedBy=default.target
EOF

# Step 4: Reload systemd, enable, and start the service
echo "Enabling and starting the wallpaper service..."
sudo systemctl daemon-reload
sudo systemctl enable "$SERVICE_NAME"
sudo systemctl start "$SERVICE_NAME"

# Step 5: Confirm the service is running
echo "Checking status of the wallpaper service..."
sudo systemctl status "$SERVICE_NAME"

echo "Installation complete. The wallpaper service is now running."
