package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	repoBaseURL  = "https://raw.githubusercontent.com/toymakerftw/wallpaper/main" // Replace with your GitHub raw URL
	metadataURL  = repoBaseURL + "/metadata.txt"
	localMeta    = "/tmp/metadata.txt"
	wallpaperDir = "/tmp"
	pollInterval = 30 * time.Second // Polling interval
)

func main() {
	for {
		err := checkForUpdates()
		if err != nil {
			fmt.Println("Error checking for updates:", err)
		}
		time.Sleep(pollInterval)
	}
}

// checkForUpdates fetches the metadata file and checks if a new wallpaper needs to be applied.
func checkForUpdates() error {
	// Download the metadata file
	err := downloadFile(metadataURL, localMeta)
	if err != nil {
		return fmt.Errorf("failed to download metadata.txt: %v", err)
	}

	// Parse the metadata file
	wallpaper, err := parseMetadata(localMeta)
	if err != nil {
		return fmt.Errorf("failed to parse metadata.txt: %v", err)
	}

	// Check if the wallpaper is already applied
	localWallpaper := wallpaperDir + "/current-wallpaper.jpg"
	if _, err := os.Stat(localWallpaper); err == nil {
		current, _ := os.ReadFile(localWallpaper)
		if wallpaper == string(current) {
			fmt.Println("Wallpaper is already up to date.")
			return nil
		}
	}

	// Download and apply the new wallpaper
	wallpaperURL := repoBaseURL + "/" + wallpaper
	localWallpaperPath := wallpaperDir + "/" + wallpaper

	err = downloadFile(wallpaperURL, localWallpaperPath)
	if err != nil {
		return fmt.Errorf("failed to download wallpaper: %v", err)
	}

	err = applyWallpaper(localWallpaperPath)
	if err != nil {
		return fmt.Errorf("failed to apply wallpaper: %v", err)
	}

	// Save the current wallpaper name locally
	err = os.WriteFile(localWallpaper, []byte(wallpaper), 0644)
	if err != nil {
		return fmt.Errorf("failed to update local wallpaper record: %v", err)
	}

	fmt.Println("Wallpaper updated to:", wallpaper)
	return nil
}

// downloadFile downloads a file from a given URL to a specified local path.
func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %v (status: %d)", url, resp.StatusCode)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	return nil
}

// parseMetadata reads and parses the metadata.txt file to extract the latest wallpaper file name.
func parseMetadata(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var wallpaper string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "wallpaper:") {
			wallpaper = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}

	if wallpaper == "" {
		return "", fmt.Errorf("wallpaper not specified in metadata.txt")
	}

	return wallpaper, nil
}

// applyWallpaper applies the wallpaper using gsettings (GNOME-based systems).
func applyWallpaper(filePath string) error {
	// Set wallpaper for the light theme
	cmd := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file://"+filePath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to set light theme wallpaper: %v", err)
	}

	// Set wallpaper for the dark theme
	cmdDark := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri-dark", "file://"+filePath)
	err = cmdDark.Run()
	if err != nil {
		return fmt.Errorf("failed to set dark theme wallpaper: %v", err)
	}

	fmt.Println("Wallpaper applied for both light and dark themes:", filePath)
	return nil
}
