// Package simctl provides a Go wrapper around Xcode's simctl tool for iOS simulator management.
// It mirrors the pattern from pkg/adb for Android.
package simctl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Client wraps simctl operations for iOS simulators.
type Client struct{}

// New creates a new simctl client.
func New() *Client {
	return &Client{}
}

// Available returns true if xcrun simctl is available.
func (c *Client) Available() bool {
	cmd := exec.Command("xcrun", "simctl", "help")
	return cmd.Run() == nil
}

// Device represents an iOS simulator device.
type Device struct {
	UDID    string `json:"udid"`
	Name    string `json:"name"`
	State   string `json:"state"` // "Booted", "Shutdown"
	Runtime string // e.g. "iOS 17.5"
}

// simctlDeviceJSON mirrors the simctl JSON output for a single device.
type simctlDeviceJSON struct {
	UDID       string `json:"udid"`
	Name       string `json:"name"`
	State      string `json:"state"`
	IsAvailable bool  `json:"isAvailable"`
	DeviceTypeIdentifier string `json:"deviceTypeIdentifier"`
}

// simctlListJSON mirrors the top-level simctl JSON output.
type simctlListJSON struct {
	Devices map[string][]simctlDeviceJSON `json:"devices"`
}

// run executes a simctl command and returns combined output.
func (c *Client) run(args ...string) (string, error) {
	fullArgs := append([]string{"simctl"}, args...)
	cmd := exec.Command("xcrun", fullArgs...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("xcrun simctl %s: %w\n%s", strings.Join(args, " "), err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// runPassthrough executes a simctl command with stdout/stderr connected to the terminal.
func (c *Client) runPassthrough(args ...string) error {
	fullArgs := append([]string{"simctl"}, args...)
	cmd := exec.Command("xcrun", fullArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Devices lists all available simulator devices.
func (c *Client) Devices() ([]Device, error) {
	out, err := c.run("list", "devices", "-j")
	if err != nil {
		return nil, err
	}

	var result simctlListJSON
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		return nil, fmt.Errorf("parse simctl output: %w", err)
	}

	var devices []Device
	for runtime, devs := range result.Devices {
		// Extract readable runtime name from key like "com.apple.CoreSimulator.SimRuntime.iOS-17-5"
		runtimeName := parseRuntimeName(runtime)
		for _, d := range devs {
			if !d.IsAvailable {
				continue
			}
			devices = append(devices, Device{
				UDID:    d.UDID,
				Name:    d.Name,
				State:   d.State,
				Runtime: runtimeName,
			})
		}
	}
	return devices, nil
}

// parseRuntimeName converts "com.apple.CoreSimulator.SimRuntime.iOS-17-5" to "iOS 17.5".
func parseRuntimeName(key string) string {
	// Try to get the last component after the prefix
	parts := strings.Split(key, ".")
	if len(parts) > 0 {
		last := parts[len(parts)-1]
		// Replace hyphens with dots for version, keep platform prefix
		components := strings.SplitN(last, "-", 2)
		if len(components) == 2 {
			version := strings.ReplaceAll(components[1], "-", ".")
			return components[0] + " " + version
		}
		return last
	}
	return key
}

// BootedDevices returns only devices that are currently booted.
func (c *Client) BootedDevices() ([]Device, error) {
	devices, err := c.Devices()
	if err != nil {
		return nil, err
	}
	var booted []Device
	for _, d := range devices {
		if d.State == "Booted" {
			booted = append(booted, d)
		}
	}
	return booted, nil
}

// HasBooted returns true if at least one simulator is booted.
func (c *Client) HasBooted() bool {
	devices, err := c.BootedDevices()
	if err != nil {
		return false
	}
	return len(devices) > 0
}

// Boot starts a simulator device. No-op if already booted.
func (c *Client) Boot(udid string) error {
	_, err := c.run("boot", udid)
	if err != nil && strings.Contains(err.Error(), "current state: Booted") {
		return nil // Already booted
	}
	return err
}

// Shutdown stops a simulator device.
func (c *Client) Shutdown(udid string) error {
	_, err := c.run("shutdown", udid)
	return err
}

// OpenSimulatorApp launches the Simulator.app GUI.
func (c *Client) OpenSimulatorApp() error {
	cmd := exec.Command("open", "-a", "Simulator")
	return cmd.Run()
}

// Install installs an .app bundle onto the booted simulator.
func (c *Client) Install(appPath string) error {
	return c.runPassthrough("install", "booted", appPath)
}

// Uninstall removes an app by bundle ID.
func (c *Client) Uninstall(bundleID string) error {
	return c.runPassthrough("uninstall", "booted", bundleID)
}

// Launch starts an app by bundle ID on the booted simulator.
func (c *Client) Launch(bundleID string) error {
	return c.runPassthrough("launch", "booted", bundleID)
}

// Terminate stops an app by bundle ID.
func (c *Client) Terminate(bundleID string) error {
	_, err := c.run("terminate", "booted", bundleID)
	return err
}

// Screenshot captures the booted simulator screen to a local file.
func (c *Client) Screenshot(outputPath string) error {
	return c.runPassthrough("io", "booted", "screenshot", outputPath)
}

// StatusBarOverride sets a clean status bar for screenshots (iOS 13+).
func (c *Client) StatusBarOverride() error {
	_, err := c.run("status_bar", "booted", "override",
		"--time", "9:41",
		"--batteryState", "charged",
		"--batteryLevel", "100",
		"--wifiBars", "3",
		"--cellularBars", "4",
	)
	return err
}

// StatusBarClear removes the status bar override.
func (c *Client) StatusBarClear() error {
	_, err := c.run("status_bar", "booted", "clear")
	return err
}

// GetAppContainer returns the data container path for an installed app.
func (c *Client) GetAppContainer(bundleID, containerType string) (string, error) {
	if containerType == "" {
		containerType = "app"
	}
	return c.run("get_app_container", "booted", bundleID, containerType)
}

// ListDeviceTypes returns available device types (iPhone 15, iPad Pro, etc.).
func (c *Client) ListDeviceTypes() (string, error) {
	return c.run("list", "devicetypes")
}

// ListRuntimes returns available runtimes (iOS versions).
func (c *Client) ListRuntimes() (string, error) {
	return c.run("list", "runtimes")
}
