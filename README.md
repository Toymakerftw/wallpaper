# Wallpaper Service

This project provides a wallpaper service that periodically updates the wallpaper on your system based on metadata stored in a remote GitHub repository. It downloads the latest wallpaper image and applies it automatically.

### Features:
- Downloads the latest wallpaper from a remote GitHub repository.
- Periodically checks for wallpaper updates.
- Automatically sets the wallpaper for both light and dark themes.
- Runs as a background service, ensuring the wallpaper stays updated.
  
## Installation

To install the wallpaper service, follow these steps:

1. **Clone or Download the Repository**

   If you haven’t already, you can download the wallpaper repository and binary from the GitHub link provided. The following script automates the entire installation process.

2. **Install the Service Using the Script**

   Download and run the installation script to set up the wallpaper service:

   ```bash
   curl -LO https://raw.githubusercontent.com/Toymakerftw/wallpaper/main/install_wallpaper_service.sh
   chmod +x install_wallpaper_service.sh
   ./install_wallpaper_service.sh
   ```

   This script will:
   - Download the wallpaper service binary.
   - Make the binary executable.
   - Set up a systemd service to run the wallpaper service as the current user.
   - Start the wallpaper service immediately.
   - Enable the service to run automatically on system boot.

3. **Service Configuration**

   The wallpaper service will run under the current user’s session, which ensures that the wallpaper is applied without needing root permissions.

   The systemd service will also be configured to handle both the light and dark themes of your desktop environment.

## Usage

Once the installation is complete, the wallpaper service will automatically start and check for wallpaper updates at a regular interval (default: every 30 seconds). It will download the latest wallpaper based on the metadata and apply it.

To check the status of the service, you can use:

```bash
systemctl status wallpaper-service
```

To manually restart the service:

```bash
sudo systemctl restart wallpaper-service
```

To stop the service:

```bash
sudo systemctl stop wallpaper-service
```

### How It Works

1. The service checks for updates by downloading a metadata file from a GitHub repository.
2. The metadata file contains the filename of the latest wallpaper.
3. The service compares the current wallpaper to the one specified in the metadata file.
4. If the wallpaper has changed, it downloads the new wallpaper and sets it as the desktop background.
5. The service runs as the current user, ensuring that it has the necessary permissions to change the wallpaper.

### Requirements

- **Linux-based OS** with `systemd` for service management.
- **Gnome Desktop Environment** or any other compatible environment that supports `gsettings` for applying wallpapers.
- **curl** to download the binary.
- **DBUS_SESSION_BUS_ADDRESS** to interact with the user’s session.

## Troubleshooting

- If the wallpaper is not updating, check the log files for any errors:

  ```bash
  sudo journalctl -u wallpaper-service
  ```

- Make sure that the environment variables `DISPLAY` and `DBUS_SESSION_BUS_ADDRESS` are correctly set. These are necessary for applying the wallpaper.

## Uninstalling the Service

To remove the wallpaper service:

1. Stop and disable the service:

   ```bash
   sudo systemctl stop wallpaper-service
   sudo systemctl disable wallpaper-service
   ```

2. Remove the systemd service file:

   ```bash
   sudo rm /etc/systemd/system/wallpaper-service.service
   ```

3. Delete the wallpaper binary:

   ```bash
   rm ~/wallpaper
   ```

4. Optionally, delete any log files or metadata:

   ```bash
   rm /tmp/wallpaper-service.log
   rm /tmp/metadata.txt
   ```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
