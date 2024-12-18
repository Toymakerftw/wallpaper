package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/go-github/v42/github"
	"github.com/m1cr0man/go-wallpaper"
)

// GitHubConfig holds GitHub repository-specific configuration
type GitHubConfig struct {
	RepositoryOwner string `json:"repository_owner"`
	RepositoryName  string `json:"repository_name"`
	AssetPattern    string `json:"asset_pattern"`
}

// Config holds the entire program configuration
type Config struct {
	GitHub           GitHubConfig  `json:"github"`
	MetadataFile     string        `json:"metadata_file"`
	DownloadPath     string        `json:"download_path"`
	UpdateFrequency  time.Duration `json:"update_frequency"`
	AllowedFileTypes []string      `json:"allowed_file_types"`
}

// Wallpaper represents a downloaded wallpaper's metadata
type Wallpaper struct {
	URL       string    `json:"url"`
	Filename  string    `json:"filename"`
	UpdatedAt time.Time `json:"updated_at"`
}

// loadConfig reads and validates configuration from a JSON file
func loadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults for missing configurations
	if config.UpdateFrequency == 0 {
		config.UpdateFrequency = 24 * time.Hour
	}
	if len(config.AllowedFileTypes) == 0 {
		config.AllowedFileTypes = []string{".jpg", ".png", ".jpeg"}
	}
	if config.DownloadPath == "" {
		config.DownloadPath = filepath.Join(os.TempDir(), "wallpapers")
	}

	if err := os.MkdirAll(config.DownloadPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create download path: %w", err)
	}

	// Environment variable overrides
	overrideConfigFromEnv(&config)

	// Validate required fields
	if config.GitHub.RepositoryOwner == "" || config.GitHub.RepositoryName == "" {
		return nil, errors.New("GitHub repository owner and name must be specified")
	}

	return &config, nil
}

func overrideConfigFromEnv(config *Config) {
	if owner := os.Getenv("GITHUB_REPO_OWNER"); owner != "" {
		config.GitHub.RepositoryOwner = owner
	}
	if name := os.Getenv("GITHUB_REPO_NAME"); name != "" {
		config.GitHub.RepositoryName = name
	}
	if pattern := os.Getenv("GITHUB_ASSET_PATTERN"); pattern != "" {
		config.GitHub.AssetPattern = pattern
	}
}

// loadOrCreateMetadata retrieves or initializes the metadata file
func loadOrCreateMetadata(metadataFile string) (*Wallpaper, error) {
	if _, err := os.Stat(metadataFile); os.IsNotExist(err) {
		return initializeMetadataFile(metadataFile)
	}

	data, err := os.ReadFile(metadataFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var wallpaper Wallpaper
	if err := json.Unmarshal(data, &wallpaper); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}
	return &wallpaper, nil
}

func initializeMetadataFile(metadataFile string) (*Wallpaper, error) {
	emptyMetadata := &Wallpaper{}
	data, err := json.MarshalIndent(emptyMetadata, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize metadata: %w", err)
	}
	if err := os.WriteFile(metadataFile, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write metadata: %w", err)
	}
	return emptyMetadata, nil
}

// findWallpaperAsset retrieves a suitable asset from the GitHub release
func findWallpaperAsset(release *github.RepositoryRelease, config *Config) *github.ReleaseAsset {
	for _, asset := range release.Assets {
		if isAssetMatching(asset.GetName(), config) {
			return asset
		}
	}
	return nil
}

func isAssetMatching(filename string, config *Config) bool {
	if config.GitHub.AssetPattern != "" {
		match, _ := filepath.Match(config.GitHub.AssetPattern, filename)
		return match
	}
	for _, ext := range config.AllowedFileTypes {
		if filepath.Ext(filename) == ext {
			return true
		}
	}
	return false
}

// needsUpdate checks if a new wallpaper is available
func needsUpdate(ctx context.Context, current *Wallpaper, config *Config) (bool, *github.RepositoryRelease, error) {
	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(ctx, config.GitHub.RepositoryOwner, config.GitHub.RepositoryName)
	if err != nil {
		return false, nil, fmt.Errorf("failed to fetch latest release: %w", err)
	}

	if current.UpdatedAt.IsZero() || release.GetPublishedAt().After(current.UpdatedAt) {
		return true, release, nil
	}
	return false, nil, nil
}

// downloadWallpaper downloads the specified wallpaper asset
func downloadWallpaper(ctx context.Context, asset *github.ReleaseAsset, config *Config) (*Wallpaper, error) {
	resp, err := http.Get(asset.GetBrowserDownloadURL())
	if err != nil {
		return nil, fmt.Errorf("failed to download wallpaper: %w", err)
	}
	defer resp.Body.Close()

	filePath := filepath.Join(config.DownloadPath, asset.GetName())
	out, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return &Wallpaper{
		URL:       asset.GetBrowserDownloadURL(),
		Filename:  asset.GetName(),
		UpdatedAt: time.Now(),
	}, nil
}

// setWallpaper sets the downloaded file as the system wallpaper
func setWallpaper(filePath string) error {
	// Set wallpaper for light mode
	if err := wallpaper.SetFromFile(filePath); err != nil {
		return fmt.Errorf("failed to set wallpaper for light mode: %w", err)
	}

	// Check if Gnome's dark mode wallpaper is supported
	out, err := exec.Command("gsettings", "writable", "org.gnome.desktop.background", "picture-uri-dark").Output()
	if err == nil && string(out) == "true\n" {
		// Set wallpaper for dark mode
		err := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri-dark", strconv.Quote("file://"+filePath)).Run()
		if err != nil {
			return fmt.Errorf("failed to set wallpaper for dark mode: %w", err)
		}
		log.Println("Dark mode wallpaper successfully set.")
	} else {
		log.Println("Dark mode wallpaper is not supported or writable on this system.")
	}

	log.Println("Light mode wallpaper successfully set.")
	return nil
}

// saveMetadata persists updated wallpaper metadata to a file
func saveMetadata(metadataFile string, wallpaper *Wallpaper) error {
	data, err := json.MarshalIndent(wallpaper, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	return os.WriteFile(metadataFile, data, 0644)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Load configuration
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Load metadata
	currentWallpaper, err := loadOrCreateMetadata(config.MetadataFile)
	if err != nil {
		log.Fatalf("Error loading metadata: %v", err)
	}

	// Check for updates
	needsUpdate, release, err := needsUpdate(ctx, currentWallpaper, config)
	if err != nil {
		log.Fatalf("Update check failed: %v", err)
	}

	if !needsUpdate {
		log.Println("No new wallpaper updates available.")
		return
	}

	// Find and download the wallpaper
	asset := findWallpaperAsset(release, config)
	if asset == nil {
		log.Fatalf("No suitable wallpaper asset found.")
	}

	newWallpaper, err := downloadWallpaper(ctx, asset, config)
	if err != nil {
		log.Fatalf("Error downloading wallpaper: %v", err)
	}

	// Set the wallpaper
	filePath := filepath.Join(config.DownloadPath, newWallpaper.Filename)
	if err := setWallpaper(filePath); err != nil {
		log.Fatalf("Failed to set wallpaper: %v", err)
	}

	// Save updated metadata
	if err := saveMetadata(config.MetadataFile, newWallpaper); err != nil {
		log.Fatalf("Failed to save metadata: %v", err)
	}

	log.Printf("Wallpaper updated successfully: %s", newWallpaper.Filename)
}
