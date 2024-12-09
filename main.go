package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	repoBaseURL  = "https://raw.githubusercontent.com/toymakerftw/wallpaper/main" // Replace with your GitHub raw URL
	metadataURL  = repoBaseURL + "/metadata.txt"
	localMeta    = "/tmp/metadata.txt"
	wallpaperDir = "/tmp"
	logFilePath  = "/tmp/wallpaper-service.log" // Path to the log file
	pollInterval = 30 * time.Second             // Polling interval
)

func main() {
	// Initialize logging
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		os.Exit(1)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.Println("Service started...")

	// Main loop
	for {
		err := checkForUpdates()
		if err != nil {
			log.Println("Error checking for updates:", err)
		}
		time.Sleep(pollInterval)
	}
}

// checkForUpdates fetches the metadata file and checks if a new wallpaper needs to be applied.
func checkForUpdates() error {
	log.Println("Checking for updates...")
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

	// Check if the wallpaper is already applied by comparing image hashes
	localWallpaper := wallpaperDir + "/current-wallpaper.jpg"
	if _, err := os.Stat(localWallpaper); err == nil {
		// Calculate the hash of the current and new wallpaper
		existingHash, err := getImageHash(localWallpaper)
		if err != nil {
			log.Println("Error calculating hash of the current wallpaper:", err)
		} else {
			// Download the new wallpaper temporarily to compare
			newWallpaperPath := wallpaperDir + "/new-wallpaper.jpg"
			wallpaperURL := repoBaseURL + "/" + wallpaper
			err = downloadFile(wallpaperURL, newWallpaperPath)
			if err != nil {
				return fmt.Errorf("failed to download wallpaper: %v", err)
			}

			newHash, err := getImageHash(newWallpaperPath)
			if err != nil {
				log.Println("Error calculating hash of the new wallpaper:", err)
			} else if existingHash == newHash {
				log.Println("Wallpaper is already up to date.")
				return nil
			}
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

	log.Println("Wallpaper updated to:", wallpaper)
	return nil
}

// downloadFile downloads a file from a given URL to a specified local path.
func downloadFile(url, filepath string) error {
	log.Printf("Downloading file from %s to %s\n", url, filepath)
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

	log.Printf("File downloaded successfully: %s\n", filepath)
	return nil
}

// parseMetadata reads and parses the metadata.txt file to extract the latest wallpaper file name.
func parseMetadata(filepath string) (string, error) {
	log.Printf("Parsing metadata file: %s\n", filepath)
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

	log.Printf("Parsed wallpaper: %s\n", wallpaper)
	return wallpaper, nil
}

// applyWallpaper applies the wallpaper using gsettings (GNOME-based systems).
func applyWallpaper(filePath string) error {
	log.Printf("Applying wallpaper: %s\n", filePath)

	// Ensure absolute file path
	filePath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %v", err)
	}

	// Ensure the DBUS_SESSION_BUS_ADDRESS environment variable is set
	os.Setenv("DBUS_SESSION_BUS_ADDRESS", getDBusSession())
	os.Setenv("DISPLAY", ":0") // Replace with your DISPLAY value if different

	// Set wallpaper for the light theme
	cmd := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file://"+filePath)
	lightOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to set light theme wallpaper: %s\nOutput: %s\n", err, string(lightOutput))
		return fmt.Errorf("failed to set light theme wallpaper: %v", err)
	}

	// Set wallpaper for the dark theme
	cmdDark := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri-dark", "file://"+filePath)
	darkOutput, err := cmdDark.CombinedOutput()
	if err != nil {
		log.Printf("Failed to set dark theme wallpaper: %s\nOutput: %s\n", err, string(darkOutput))
		return fmt.Errorf("failed to set dark theme wallpaper: %v", err)
	}

	log.Println("Wallpaper applied for both light and dark themes:", filePath)
	return nil
}

// getDBusSession retrieves the DBUS_SESSION_BUS_ADDRESS environment variable.
func getDBusSession() string {
	cmd := exec.Command("bash", "-c", "echo $DBUS_SESSION_BUS_ADDRESS")
	output, err := cmd.Output()
	if err != nil {
		log.Println("Failed to get DBUS_SESSION_BUS_ADDRESS:", err)
		return ""
	}
	return strings.TrimSpace(string(output))
}

// getImageHash calculates the SHA-256 hash of the image at the given path.
func getImageHash(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Calculate the hash of the file content
	hash := sha256.New()
	_, err = io.Copy(hash, file) // Changed from ioutil.Copy to io.Copy
	if err != nil {
		return "", fmt.Errorf("failed to calculate file hash: %v", err)
	}

	// Return the hexadecimal representation of the hash
	return hex.EncodeToString(hash.Sum(nil)), nil
}
