# <img src="https://raw.githubusercontent.com/Toymakerftw/wallpaper/refs/heads/main/banner.png" > Ubuntu Wallpaper Manager

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/Toymakerftw/wallpaper/release.yml?branch=main)](https://github.com/Toymakerftw/wallpaper/actions) [![GitHub Release](https://img.shields.io/github/v/release/Toymakerftw/wallpaper)](https://github.com/Toymakerftw/wallpaper/releases) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![Go Report Card](https://goreportcard.com/badge/github.com/Toymakerftw/wallpaper)](https://goreportcard.com/report/github.com/Toymakerftw/wallpaper) [![Go Version](https://img.shields.io/badge/Go-1.18-blue.svg)](https://golang.org/doc/go1.18) [![Ubuntu Version](https://img.shields.io/badge/Ubuntu-24.04-orange.svg)](https://ubuntu.com/) ![Visitor Count](https://visitor-badge.laobi.icu/badge?page_id=Toymakerftw.wallpaper)






## Introduction

The Ubuntu Wallpaper Manager is a Go-based application designed to run as a service on Ubuntu 24 laptops. It periodically checks for new wallpapers on a specified GitHub repository and updates the system wallpaper accordingly.

## Features

* Periodically checks for new wallpapers on a GitHub repository
* Updates the system wallpaper with the latest image
* Supports both light and dark mode wallpapers
* Sends a notification when a new wallpaper is set
* Allows users to disable the service via a notification action

## Installation

To install the Ubuntu Wallpaper Manager, run the following command:
```bash
curl  -sSL  https://github.com/Toymakerftw/wallpaper/raw/refs/heads/main/install.sh | sudo  bash
```
This script will:

1. Install the required dependencies
2. Create a systemd service file and enable the service
3. Configure the application to run at startup

## Uninstallation

To uninstall the Ubuntu Wallpaper Manager, run the following command:
```bash
curl  -sSL  https://github.com/Toymakerftw/wallpaper/raw/refs/heads/main/uninstall.sh | sudo  bash
```
This script will:

1. Stop and disable the systemd service
2. Remove the systemd service file
3. Remove the application and its dependencies

## Configuration

The application uses a JSON configuration file (`config.json`) to store settings. You can modify this file to change the GitHub repository, update frequency, and other settings.

Example `config.json` file:
```json
{
  "github": {
    "repository_owner": "Toymakerftw",
    "repository_name": "wallpaper",
    "asset_pattern": "wallpaper-*"
  },
  "metadata_file": "wallpaper_metadata.json",
  "download_path": "/etc/wallpaper",
  "update_frequency": 86400,
  "allowed_file_types": [".jpg", ".png", ".jpeg"]
}
```
## Contributing

Contributions are welcome! If you'd like to contribute to the Ubuntu Wallpaper Manager, please fork the repository and submit a pull request.

## License

The Ubuntu Wallpaper Manager is licensed under the MIT License. See [LICENSE](LICENSE) for details.
