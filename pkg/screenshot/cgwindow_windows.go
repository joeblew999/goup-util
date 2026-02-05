// +build windows

package screenshot

import (
	"fmt"
	"time"
)

// GetCGWindowBoundsByPID - Windows implementation placeholder
// TODO: Implement using Win32 EnumWindows/GetWindowRect
func GetCGWindowBoundsByPID(pid int) (int, int, int, int, error) {
	return 0, 0, 0, 0, fmt.Errorf("CoreGraphics not available on Windows, will fall back to robotgo")
}

// GetCGWindowIDByPID - Windows implementation placeholder
func GetCGWindowIDByPID(pid int) (int, error) {
	return 0, fmt.Errorf("CoreGraphics not available on Windows")
}

// WaitForCGWindow - Windows implementation placeholder
// On Windows, robotgo's native window detection should work better
func WaitForCGWindow(pid int, timeout time.Duration) error {
	return fmt.Errorf("CoreGraphics not available on Windows, will fall back to robotgo")
}

// CaptureWindowByCGBounds - Windows implementation placeholder
func CaptureWindowByCGBounds(pid int, output string, quality int) error {
	return fmt.Errorf("CoreGraphics not available on Windows, use robotgo methods")
}
