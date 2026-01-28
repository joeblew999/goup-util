// Package utm provides UTM virtual machine management for goup-util.
//
// Architecture:
//   - This package handles VM configuration, paths, and gallery management
//   - cmd/utm.go provides the CLI interface using this package
//   - Taskfile.utm.yml handles environment setup (downloading UTM app, ISOs)
//
// Paths (all local to project):
//   - .bin/UTM.app       - UTM application
//   - .data/utm/vms/     - Virtual machine files (.utm)
//   - .data/utm/iso/     - ISO images for VM creation
//   - .data/utm/share/   - Shared folder for file transfer
package utm

import (
	"os"
	"path/filepath"
	"runtime"
)

// Paths holds all UTM-related paths
type Paths struct {
	// Root is the base directory for all UTM data (default: .data/utm)
	Root string

	// App is where UTM.app is installed (default: .bin/UTM.app)
	App string

	// VMs is where virtual machines are stored
	VMs string

	// ISO is where ISO images are stored
	ISO string

	// Share is the shared folder for host<->guest file transfer
	Share string
}

// DefaultPaths returns the default UTM paths relative to project root
func DefaultPaths() Paths {
	return Paths{
		Root:  ".data/utm",
		App:   ".bin/UTM.app",
		VMs:   ".data/utm/vms",
		ISO:   ".data/utm/iso",
		Share: ".data/utm/share",
	}
}

// GetPaths returns UTM paths, using defaults if not configured
func GetPaths() Paths {
	// TODO: Load from config file if present
	return DefaultPaths()
}

// GetUTMCtlPath returns the path to the utmctl binary
func GetUTMCtlPath() string {
	if runtime.GOOS != "darwin" {
		return "" // UTM is macOS only
	}

	paths := GetPaths()

	// Check common locations in order of preference
	locations := []string{
		// Local install (preferred)
		filepath.Join(paths.App, "Contents/MacOS/utmctl"),
		// Homebrew
		"/opt/homebrew/bin/utmctl",
		"/usr/local/bin/utmctl",
		// System install
		"/Applications/UTM.app/Contents/MacOS/utmctl",
	}

	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			return loc
		}
	}

	// Fallback to PATH lookup
	return "utmctl"
}

// IsUTMInstalled checks if UTM is available
func IsUTMInstalled() bool {
	path := GetUTMCtlPath()
	if path == "utmctl" {
		// Check if utmctl is in PATH
		_, err := os.Stat(path)
		return err == nil
	}
	_, err := os.Stat(path)
	return err == nil
}

// GetVMPath returns the full path for a VM by name
func GetVMPath(vmName string) string {
	paths := GetPaths()
	return filepath.Join(paths.VMs, vmName+".utm")
}

// GetISOPath returns the full path for an ISO by name
func GetISOPath(isoName string) string {
	paths := GetPaths()
	return filepath.Join(paths.ISO, isoName)
}

// EnsureDirectories creates all required UTM directories
func EnsureDirectories() error {
	paths := GetPaths()

	dirs := []string{
		paths.Root,
		paths.VMs,
		paths.ISO,
		paths.Share,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}
