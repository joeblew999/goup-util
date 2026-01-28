// Package screenshot provides cross-platform screenshot capabilities
// using robotgo for testing and documentation purposes.
package screenshot

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/go-vgo/robotgo"
)

// DisplayInfo contains information about a display
type DisplayInfo struct {
	ID     int
	X      int
	Y      int
	Width  int
	Height int
}

// Config contains screenshot configuration
type Config struct {
	// Output file path
	Output string

	// Region to capture (optional)
	X, Y, Width, Height int

	// Capture all displays
	AllDisplays bool

	// Prefix for multi-display captures
	Prefix string

	// Quality for JPEG (1-100, only if output is .jpg/.jpeg)
	Quality int

	// Delay before capture (milliseconds)
	Delay int
}

// Capture performs a screenshot with the given config
func Capture(cfg Config) error {
	// Validate output path
	if cfg.Output == "" && !cfg.AllDisplays {
		return fmt.Errorf("output path is required")
	}

	// Ensure output directory exists
	if cfg.Output != "" {
		outputDir := filepath.Dir(cfg.Output)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Apply delay if specified
	if cfg.Delay > 0 {
		fmt.Printf("Waiting %dms before capture...\n", cfg.Delay)
		robotgo.MilliSleep(cfg.Delay)
	}

	// Capture all displays
	if cfg.AllDisplays {
		return CaptureAllDisplays(cfg.Prefix)
	}

	// Capture region
	if cfg.Width > 0 && cfg.Height > 0 {
		return CaptureRegion(cfg.X, cfg.Y, cfg.Width, cfg.Height, cfg.Output, cfg.Quality)
	}

	// Capture desktop
	return CaptureDesktop(cfg.Output, cfg.Quality)
}

// CaptureDesktop captures the entire primary display
func CaptureDesktop(output string, quality int) error {
	img, err := robotgo.CaptureImg()
	if err != nil {
		return fmt.Errorf("robotgo capture failed: %w\n\nNote: On macOS 10.15+, grant Screen Recording permission:\nSystem Settings → Privacy & Security → Screen Recording", err)
	}

	return SaveImage(img, output, quality)
}

// CaptureRegion captures a specific region of the screen
func CaptureRegion(x, y, w, h int, output string, quality int) error {
	img, err := robotgo.CaptureImg(x, y, w, h)
	if err != nil {
		return fmt.Errorf("robotgo region capture failed: %w", err)
	}

	return SaveImage(img, output, quality)
}

// CaptureAllDisplays captures all connected displays
func CaptureAllDisplays(prefix string) error {
	if prefix == "" {
		prefix = "display"
	}

	num := robotgo.DisplaysNum()
	fmt.Printf("Found %d display(s)\n", num)

	for i := 0; i < num; i++ {
		robotgo.DisplayID = i

		// Get display bounds
		x, y, w, h := robotgo.GetDisplayBounds(i)
		fmt.Printf("Capturing display %d: %dx%d at (%d,%d)\n", i, w, h, x, y)

		// Capture entire display
		img, err := robotgo.CaptureImg(x, y, w, h)
		if err != nil {
			return fmt.Errorf("failed to capture display %d: %w", i, err)
		}

		// Save with display number
		path := fmt.Sprintf("%s_%d.png", prefix, i)
		if err := SaveImage(img, path, 90); err != nil {
			return fmt.Errorf("failed to save display %d: %w", i, err)
		}

		absPath, _ := filepath.Abs(path)
		fmt.Printf("✓ Saved %s\n", absPath)
	}

	return nil
}

// GetInfo returns information about all displays
func GetInfo() ([]DisplayInfo, error) {
	num := robotgo.DisplaysNum()
	displays := make([]DisplayInfo, num)

	for i := 0; i < num; i++ {
		x, y, w, h := robotgo.GetDisplayBounds(i)
		displays[i] = DisplayInfo{
			ID:     i,
			X:      x,
			Y:      y,
			Width:  w,
			Height: h,
		}
	}

	return displays, nil
}

// SaveImage saves an image to disk in the appropriate format
func SaveImage(img image.Image, path string, quality int) error {
	// Determine format from extension
	ext := filepath.Ext(path)

	// Default quality if not set
	if quality <= 0 {
		quality = 90
	}

	switch ext {
	case ".jpg", ".jpeg":
		return saveJPEG(img, path, quality)
	case ".png", ".PNG":
		return savePNG(img, path)
	default:
		return savePNG(img, path) // Default to PNG
	}
}

// savePNG saves an image as PNG
func savePNG(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	return nil
}

// saveJPEG saves an image as JPEG
func saveJPEG(img image.Image, path string, quality int) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	opts := &jpeg.Options{Quality: quality}
	if err := jpeg.Encode(file, img, opts); err != nil {
		return fmt.Errorf("failed to encode JPEG: %w", err)
	}

	return nil
}

// Sleep pauses execution for the specified number of milliseconds
func Sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// Platform returns a descriptive string of the screenshot method being used
func Platform() string {
	return fmt.Sprintf("robotgo on %s", runtime.GOOS)
}

// WindowInfo contains information about a window
type WindowInfo struct {
	Title  string
	PID    int
	Bounds image.Rectangle
}

// CaptureActiveWindow captures the currently focused window
func CaptureActiveWindow(output string, quality int) error {
	// Get active window PID
	pid := robotgo.GetPid()
	if pid == 0 {
		return fmt.Errorf("no active window found")
	}

	// Get window bounds for active PID
	x, y, w, h := robotgo.GetBounds(pid)

	if w <= 0 || h <= 0 {
		return fmt.Errorf("invalid window bounds: %dx%d at (%d,%d)", w, h, x, y)
	}

	// Capture that region
	return CaptureRegion(x, y, w, h, output, quality)
}

// CaptureWindowByPID captures window owned by process ID
func CaptureWindowByPID(pid int, output string, quality int) error {
	// Activate window first
	if err := robotgo.ActivePid(pid); err != nil {
		return fmt.Errorf("failed to activate window for PID %d: %w", pid, err)
	}

	// Give it time to come to front
	robotgo.MilliSleep(200)

	// Get window bounds
	x, y, w, h := robotgo.GetBounds(pid)

	if w <= 0 || h <= 0 {
		return fmt.Errorf("invalid window bounds for PID %d: %dx%d at (%d,%d)", pid, w, h, x, y)
	}

	// Capture that region
	return CaptureRegion(x, y, w, h, output, quality)
}

// WaitForWindow polls until window appears for PID
func WaitForWindow(pid int, timeout time.Duration) error {
	start := time.Now()

	for time.Since(start) < timeout {
		// Check if window has valid bounds (width and height > 0)
		_, _, w, h := robotgo.GetBounds(pid)
		if w > 0 && h > 0 {
			return nil
		}

		robotgo.MilliSleep(100)
	}

	return fmt.Errorf("window for PID %d did not appear within %v", pid, timeout)
}
