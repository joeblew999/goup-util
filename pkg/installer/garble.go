package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	GarbleVersion     = "latest" // Always use latest version for best Go compatibility
	GarblePackage     = "mvdan.cc/garble"
	GarbleInstallPath = "sdks/tools/garble" // Relative to SDK root
)

// InstallGarble installs garble to SDK directory using go install
func InstallGarble(cache *Cache) error {
	fmt.Printf("📥 Installing garble %s to SDK directory...\n", GarbleVersion)

	// Resolve SDK install path
	installPath, err := ResolveInstallPath(GarbleInstallPath)
	if err != nil {
		return fmt.Errorf("failed to resolve install path: %w", err)
	}

	// Check if already installed in SDK directory
	garbleBinary := "garble"
	if runtime.GOOS == "windows" {
		garbleBinary = "garble.exe"
	}
	garblePath := filepath.Join(installPath, garbleBinary)

	if entry, ok := cache.Entries["garble"]; ok {
		if _, err := os.Stat(garblePath); err == nil {
			fmt.Printf("✅ garble %s is already installed at: %s\n", entry.Version, garblePath)
			return nil
		}
	}

	// Check if Go is installed
	if _, err := exec.LookPath("go"); err != nil {
		return fmt.Errorf("go command not found. Please install Go first")
	}

	// Create install directory
	if err := os.MkdirAll(installPath, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}

	// Run go install with GOBIN set to SDK directory
	cmd := exec.Command("go", "install", GarblePackage+"@"+GarbleVersion)
	cmd.Env = append(os.Environ(), "GOBIN="+installPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install garble: %w", err)
	}

	// Verify installation
	if _, err := os.Stat(garblePath); err != nil {
		return fmt.Errorf("garble binary not found at %s after installation", garblePath)
	}

	fmt.Printf("✅ garble installed successfully at: %s\n", garblePath)

	// Test garble version
	versionCmd := exec.Command(garblePath, "version")
	output, err := versionCmd.Output()
	if err != nil {
		fmt.Printf("⚠️  Warning: Could not verify garble version: %v\n", err)
	} else {
		fmt.Printf("   Version: %s\n", strings.TrimSpace(string(output)))
	}

	// Add to cache
	cache.Add(&SDK{
		Name:        "garble",
		Version:     GarbleVersion,
		Checksum:    "go-install", // Special marker for go-install tools
		InstallPath: GarbleInstallPath,
	})

	if err := cache.Save(); err != nil {
		fmt.Printf("⚠️  Warning: Could not update cache: %v\n", err)
	}

	return nil
}

// IsGarbleInstalled checks if garble is available in SDK directory
func IsGarbleInstalled() bool {
	garblePath, err := GetGarblePath()
	if err != nil {
		return false
	}
	_, err = os.Stat(garblePath)
	return err == nil
}

// GetGarblePath returns the path to garble binary in SDK directory
func GetGarblePath() (string, error) {
	installPath, err := ResolveInstallPath(GarbleInstallPath)
	if err != nil {
		return "", err
	}

	garbleBinary := "garble"
	if runtime.GOOS == "windows" {
		garbleBinary = "garble.exe"
	}

	return filepath.Join(installPath, garbleBinary), nil
}
