package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed sdk-android-list.json
var AndroidSdkList []byte

//go:embed sdk-ios-list.json
var IosSdkList []byte

// Platform defines the structure for platform-specific SDK details.
type Platform struct {
	DownloadURL string `json:"downloadUrl"`
	Checksum    string `json:"checksum"`
}

// SdkItem defines the structure for an SDK entry in the JSON file.
type SdkItem struct {
	Version        string              `json:"version"`
	GoupName       string              `json:"goupName"`
	DownloadURL    string              `json:"downloadUrl,omitempty"`
	Checksum       string              `json:"checksum,omitempty"`
	InstallPath    string              `json:"installPath"`
	ApiLevel       int                 `json:"apiLevel"`
	Abi            string              `json:"abi"`
	Vendor         string              `json:"vendor"`
	Platforms      map[string]Platform `json:"platforms,omitempty"`
	SdkManagerName string              `json:"sdkmanagerName,omitempty"`
}

// SdkFile defines the top-level structure of the JSON file.
type SdkFile struct {
	SDKs map[string][]SdkItem `json:"sdks"`
}

// MetaFile defines the structure for setup configurations
type MetaFile struct {
	Meta struct {
		Setups map[string][]string `json:"setups"`
	} `json:"meta"`
}

// GetCacheDir returns the OS-appropriate cache directory for goup-util
func GetCacheDir() string {
	switch runtime.GOOS {
	case "darwin": // macOS
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, "goup-util-cache")
		}
	case "linux":
		if cacheHome := os.Getenv("XDG_CACHE_HOME"); cacheHome != "" {
			return filepath.Join(cacheHome, "goup-util")
		}
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, ".cache", "goup-util")
		}
	case "windows":
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			return filepath.Join(localAppData, "goup-util")
		}
	}

	// Fallback to the old behavior if we can't determine OS-specific path
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".goup-util")
	}

	// Last resort fallback
	return ".goup-util"
}

// GetSDKDir returns the OS-appropriate SDK storage directory for goup-util
func GetSDKDir() string {
	switch runtime.GOOS {
	case "darwin": // macOS
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, "goup-util-sdks")
		}
	case "linux":
		if dataHome := os.Getenv("XDG_DATA_HOME"); dataHome != "" {
			return filepath.Join(dataHome, "goup-util", "sdks")
		}
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, ".local", "share", "goup-util", "sdks")
		}
	case "windows":
		if appData := os.Getenv("APPDATA"); appData != "" {
			return filepath.Join(appData, "goup-util", "sdks")
		}
	}

	// Fallback - use cache dir + sdks subdirectory
	return filepath.Join(GetCacheDir(), "sdks")
}

// GetCachePath returns the full path to the cache.json file
func GetCachePath() string {
	return filepath.Join(GetCacheDir(), "cache.json")
}

// DirectoryInfo contains information about goup-util directories
type DirectoryInfo struct {
	CacheDir    string `json:"cache_dir"`
	SDKDir      string `json:"sdk_dir"`
	CacheExists bool   `json:"cache_exists"`
	SDKExists   bool   `json:"sdk_exists"`
	CacheSize   int64  `json:"cache_size,omitempty"`
	SDKSize     int64  `json:"sdk_size,omitempty"`
}

// EnsureDirectories creates cache and SDK directories if they don't exist
func EnsureDirectories() error {
	cacheDir := GetCacheDir()
	sdkDir := GetSDKDir()

	// Create cache directory
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory %s: %w", cacheDir, err)
	}

	// Create SDK directory
	if err := os.MkdirAll(sdkDir, 0755); err != nil {
		return fmt.Errorf("failed to create SDK directory %s: %w", sdkDir, err)
	}

	return nil
}

// CleanDirectories removes all cache and SDK directories
func CleanDirectories() error {
	var errors []error

	// Remove SDK directory
	sdkDir := GetSDKDir()
	if _, err := os.Stat(sdkDir); err == nil {
		if err := os.RemoveAll(sdkDir); err != nil {
			errors = append(errors, fmt.Errorf("failed to remove SDK directory %s: %w", sdkDir, err))
		}
	}

	// Remove cache directory
	cacheDir := GetCacheDir()
	if _, err := os.Stat(cacheDir); err == nil {
		if err := os.RemoveAll(cacheDir); err != nil {
			errors = append(errors, fmt.Errorf("failed to remove cache directory %s: %w", cacheDir, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("cleanup errors: %v", errors)
	}

	return nil
}

// CleanCache removes only the cache directory
func CleanCache() error {
	cacheDir := GetCacheDir()
	if _, err := os.Stat(cacheDir); err == nil {
		if err := os.RemoveAll(cacheDir); err != nil {
			return fmt.Errorf("failed to remove cache directory %s: %w", cacheDir, err)
		}
	}
	return nil
}

// GetDirectoryInfo returns size and health information about directories
func GetDirectoryInfo() DirectoryInfo {
	cacheDir := GetCacheDir()
	sdkDir := GetSDKDir()

	info := DirectoryInfo{
		CacheDir: cacheDir,
		SDKDir:   sdkDir,
	}

	// Check if directories exist and get sizes
	if stat, err := os.Stat(cacheDir); err == nil && stat.IsDir() {
		info.CacheExists = true
		if size, err := getDirSize(cacheDir); err == nil {
			info.CacheSize = size
		}
	}

	if stat, err := os.Stat(sdkDir); err == nil && stat.IsDir() {
		info.SDKExists = true
		if size, err := getDirSize(sdkDir); err == nil {
			info.SDKSize = size
		}
	}

	return info
}

// getDirSize calculates the total size of a directory
func getDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}
