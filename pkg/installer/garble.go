package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	GarbleVersion = "v0.15.0"
	GarblePackage = "mvdan.cc/garble"
)

// InstallGarble installs garble using go install and adds it to the cache
func InstallGarble(cache *Cache) error {
	fmt.Printf("📥 Installing garble %s...\n", GarbleVersion)

	// Check if already in cache and installed
	garblePath, _ := exec.LookPath("garble")
	if entry, ok := cache.Entries["garble"]; ok && garblePath != "" {
		fmt.Printf("✅ garble %s is already installed at: %s\n", entry.Version, garblePath)
		return nil
	}

	// Check if Go is installed
	if _, err := exec.LookPath("go"); err != nil {
		return fmt.Errorf("go command not found. Please install Go first")
	}

	// Run go install
	cmd := exec.Command("go", "install", GarblePackage+"@"+GarbleVersion)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install garble: %w", err)
	}

	// Verify installation
	garblePath, err := exec.LookPath("garble")
	if err != nil {
		return fmt.Errorf("garble installed but not found in PATH. Make sure GOBIN is in your PATH")
	}

	fmt.Printf("✅ garble installed successfully at: %s\n", garblePath)

	// Test garble version
	versionCmd := exec.Command("garble", "version")
	output, err := versionCmd.Output()
	if err != nil {
		fmt.Printf("⚠️  Warning: Could not verify garble version: %v\n", err)
	} else {
		fmt.Printf("   Version: %s", string(output))
	}

	// Add to cache
	cache.Add(&SDK{
		Name:        "garble",
		Version:     GarbleVersion,
		Checksum:    "go-install", // Special marker for go-install tools
		InstallPath: garblePath,
	})

	if err := cache.Save(); err != nil {
		fmt.Printf("⚠️  Warning: Could not update cache: %v\n", err)
	}

	return nil
}

// IsGarbleInstalled checks if garble is available in PATH
func IsGarbleInstalled() bool {
	_, err := exec.LookPath("garble")
	return err == nil
}

// GetGarblePath returns the path to garble binary
func GetGarblePath() (string, error) {
	return exec.LookPath("garble")
}

// GetGarbleCommand returns the command to run garble
// This handles both Unix and Windows paths
func GetGarbleCommand() string {
	if runtime.GOOS == "windows" {
		return "garble.exe"
	}
	return "garble"
}

// GarbleBuild runs garble build with the given arguments
func GarbleBuild(args ...string) error {
	cmd := exec.Command("garble", append([]string{"build"}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// CheckGarbleVersion verifies the installed garble version
func CheckGarbleVersion() (string, error) {
	cmd := exec.Command("garble", "version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to check garble version: %w", err)
	}

	return string(output), nil
}

// GetGOBIN returns the GOBIN directory where garble is installed
func GetGOBIN() string {
	// Check GOBIN environment variable
	if gobin := os.Getenv("GOBIN"); gobin != "" {
		return gobin
	}

	// Default to GOPATH/bin
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		// Default GOPATH is $HOME/go
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		gopath = filepath.Join(home, "go")
	}

	return filepath.Join(gopath, "bin")
}
