package installer

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/joeblew99/goup-util/pkg/config"
)

// SDK represents a software development kit.
type SDK struct {
	Name        string
	Version     string
	URL         string
	Checksum    string
	InstallPath string
}

// Install downloads and installs an SDK.
func Install(sdk *SDK, cache *Cache) error {
	// Check if the SDK is already cached
	if cache.IsCached(sdk) {
		fmt.Printf("%s %s is already installed and up-to-date.\n", sdk.Name, sdk.Version)
		return nil
	}

	// First, resolve the installation path and check if the SDK already exists.
	dest, err := ResolveInstallPath(sdk.InstallPath)
	if err != nil {
		return err
	}

	if _, err := os.Stat(dest); err == nil {
		fmt.Printf("SDK %s is already present at %s. Adding to cache.\n", sdk.Name, dest)
		cache.Add(sdk)
		if err := cache.Save(); err != nil {
			return fmt.Errorf("failed to save cache: %w", err)
		}
		return nil
	}

	// If the SDK doesn't exist and there's no URL, it's a manual installation.
	if sdk.URL == "" {
		return fmt.Errorf("cannot automatically install SDK %s. Please install it manually (e.g., by installing or updating Xcode) and ensure it is available at %s", sdk.Name, dest)
	}

	fmt.Printf("Downloading %s %s from %s\n", sdk.Name, sdk.Version, sdk.URL)

	// Create a temporary file with the correct extension from the URL
	fileExt := filepath.Ext(sdk.URL)
	tmpFile, err := os.CreateTemp("", "sdk-download-*"+fileExt)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up the temp file

	// Get the data
	resp, err := http.Get(sdk.URL)
	if err != nil {
		return fmt.Errorf("failed to download SDK: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download SDK: received status code %d", resp.StatusCode)
	}

	// Create a tee reader to write to the file and the hash simultaneously
	hasher := sha256.New()
	teeReader := io.TeeReader(resp.Body, hasher)

	// Write the body to the temporary file
	_, err = io.Copy(tmpFile, teeReader)
	if err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write to temporary file: %w", err)
	}
	tmpFile.Close()

	// Verify the checksum
	calculatedChecksum := hex.EncodeToString(hasher.Sum(nil))
	expectedChecksum := strings.TrimPrefix(sdk.Checksum, "sha256:")

	if expectedChecksum != "" && calculatedChecksum != expectedChecksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, calculatedChecksum)
	}
	fmt.Println("Checksum verified.")

	fmt.Printf("Downloaded %s %s to %s\n", sdk.Name, sdk.Version, tmpFile.Name())

	// The destination is already resolved, so we don't need to call ResolveInstallPath again.

	// Extract the archive
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	fmt.Printf("Extracting to %s\n", dest)
	if err := Extract(tmpFile.Name(), dest); err != nil {
		return fmt.Errorf("failed to extract SDK: %w", err)
	}

	// Add to cache and save
	cache.Add(sdk)
	if err := cache.Save(); err != nil {
		return fmt.Errorf("failed to save cache: %w", err)
	}

	fmt.Printf("Successfully installed %s %s to %s\n", sdk.Name, sdk.Version, dest)

	// If OpenJDK was installed, print instructions for setting JAVA_HOME
	if strings.Contains(sdk.Name, "openjdk") {
		fmt.Println("\n---------------------------------------------------------------------")
		fmt.Println("IMPORTANT: To use this JDK for Android development with Gio,")
		fmt.Println("you need to set the JAVA_HOME environment variable.")
		fmt.Println("\nFor your current shell session, run:")
		fmt.Printf("export JAVA_HOME=\"%s\"\n", dest)
		fmt.Println("\nTo make this change permanent, add the line above to your")
		fmt.Println("shell profile file (e.g., ~/.zshrc, ~/.bash_profile).")
		fmt.Println("---------------------------------------------------------------------")
	}

	return nil
}

func ResolveInstallPath(path string) (string, error) {
	if path == "" {
		// Default to OS-appropriate SDK directory
		return config.GetSDKDir(), nil
	}

	// Expand environment variables first
	expandedPath := os.ExpandEnv(path)

	// If it's already absolute, return as-is
	if filepath.IsAbs(expandedPath) {
		return expandedPath, nil
	}

	// For relative paths starting with "sdks/", use OS-specific SDK directory
	if strings.HasPrefix(expandedPath, "sdks/") {
		return filepath.Join(config.GetSDKDir(), strings.TrimPrefix(expandedPath, "sdks/")), nil
	}

	// For other relative paths, make them relative to current directory
	if !filepath.IsAbs(expandedPath) {
		if cwd, err := os.Getwd(); err == nil {
			return filepath.Join(cwd, expandedPath), nil
		}
	}

	return expandedPath, nil
}
