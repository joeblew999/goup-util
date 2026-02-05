// Package utm provides UTM VM management functionality.
// Driver interface for version-specific UTM operations.
//
// Driver pattern adapted from:
// https://github.com/naveenrajm7/packer-plugin-utm
// See builder/utm/common/driver*.go in that repository.
package utm

import (
	"fmt"
	"strconv"
	"strings"
)

// Driver defines the interface for UTM operations.
// Different UTM versions have different capabilities, so we use
// version-specific drivers that implement this interface.
type Driver interface {
	// Version returns the UTM version string (e.g., "4.6.4")
	Version() string

	// SupportsExport returns true if this UTM version supports VM export
	SupportsExport() bool

	// SupportsImport returns true if this UTM version supports VM import
	SupportsImport() bool

	// SupportsGuestTools returns true if this UTM version provides guest tools
	SupportsGuestTools() bool

	// Export exports a VM to a .utm file (requires UTM 4.6+)
	Export(vmName, outputPath string) error

	// Import imports a VM from a .utm file (requires UTM 4.6+)
	Import(utmPath string) (string, error)

	// GuestToolsISOPath returns the path to the guest tools ISO (requires UTM 4.6+)
	GuestToolsISOPath() (string, error)

	// ExecuteOsaScript executes an embedded AppleScript
	ExecuteOsaScript(command ...string) (string, error)

	// Utmctl executes a utmctl command
	Utmctl(args ...string) (string, error)
}

// UTMVersion represents a parsed UTM version
type UTMVersion struct {
	Major int
	Minor int
	Patch int
	Raw   string
}

// ParseVersion parses a version string like "4.6.4" into components
func ParseVersion(version string) (*UTMVersion, error) {
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid version format: %s", version)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", parts[0])
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %s", parts[1])
	}

	patch := 0
	if len(parts) >= 3 {
		patch, _ = strconv.Atoi(parts[2]) // Ignore error, default to 0
	}

	return &UTMVersion{
		Major: major,
		Minor: minor,
		Patch: patch,
		Raw:   version,
	}, nil
}

// AtLeast returns true if this version is >= the given version
func (v *UTMVersion) AtLeast(major, minor int) bool {
	if v.Major > major {
		return true
	}
	if v.Major == major && v.Minor >= minor {
		return true
	}
	return false
}

// baseDriver provides common functionality for all UTM versions
type baseDriver struct {
	version *UTMVersion
}

func (d *baseDriver) Version() string {
	return d.version.Raw
}

func (d *baseDriver) ExecuteOsaScript(command ...string) (string, error) {
	return ExecuteOsaScript(command...)
}

func (d *baseDriver) Utmctl(args ...string) (string, error) {
	return RunUTMCtl(args...)
}

// driver45 implements Driver for UTM 4.5.x
type driver45 struct {
	baseDriver
}

func (d *driver45) SupportsExport() bool     { return false }
func (d *driver45) SupportsImport() bool     { return false }
func (d *driver45) SupportsGuestTools() bool { return false }

func (d *driver45) Export(vmName, outputPath string) error {
	return fmt.Errorf("UTM %s does not support VM export. Please upgrade to UTM 4.6 or later", d.version.Raw)
}

func (d *driver45) Import(utmPath string) (string, error) {
	return "", fmt.Errorf("UTM %s does not support VM import. Please upgrade to UTM 4.6 or later", d.version.Raw)
}

func (d *driver45) GuestToolsISOPath() (string, error) {
	return "", fmt.Errorf("UTM %s does not provide guest tools. Please upgrade to UTM 4.6 or later", d.version.Raw)
}

// driver46 implements Driver for UTM 4.6.x and later
type driver46 struct {
	baseDriver
}

func (d *driver46) SupportsExport() bool     { return true }
func (d *driver46) SupportsImport() bool     { return true }
func (d *driver46) SupportsGuestTools() bool { return true }

func (d *driver46) Export(vmName, outputPath string) error {
	return ExportVM(vmName, outputPath)
}

func (d *driver46) Import(utmPath string) (string, error) {
	return ImportVM(utmPath)
}

func (d *driver46) GuestToolsISOPath() (string, error) {
	return GetGuestToolsISOPath()
}

// NewDriver creates a new Driver instance based on the installed UTM version.
// Returns an error if UTM is not installed or version cannot be determined.
func NewDriver() (Driver, error) {
	versionStr, err := GetUTMVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get UTM version: %w", err)
	}

	version, err := ParseVersion(versionStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse UTM version: %w", err)
	}

	return NewDriverForVersion(version), nil
}

// NewDriverForVersion creates a Driver for a specific UTM version
func NewDriverForVersion(version *UTMVersion) Driver {
	base := baseDriver{version: version}

	// UTM 4.6+ has export/import support
	if version.AtLeast(4, 6) {
		return &driver46{baseDriver: base}
	}

	// UTM 4.5 and earlier
	return &driver45{baseDriver: base}
}

// GetDriver returns a cached driver instance (for convenience)
var cachedDriver Driver

// GetDriver returns a Driver instance, creating one if needed
func GetDriver() (Driver, error) {
	if cachedDriver != nil {
		return cachedDriver, nil
	}

	driver, err := NewDriver()
	if err != nil {
		return nil, err
	}

	cachedDriver = driver
	return driver, nil
}

// ResetDriver clears the cached driver (useful for testing)
func ResetDriver() {
	cachedDriver = nil
}

// CheckVersion checks if the installed UTM version meets minimum requirements
func CheckVersion(minMajor, minMinor int) error {
	driver, err := GetDriver()
	if err != nil {
		return err
	}

	version, err := ParseVersion(driver.Version())
	if err != nil {
		return err
	}

	if !version.AtLeast(minMajor, minMinor) {
		return fmt.Errorf("UTM %s is too old. Minimum required: %d.%d", version.Raw, minMajor, minMinor)
	}

	return nil
}
