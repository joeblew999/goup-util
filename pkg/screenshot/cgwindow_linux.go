// +build linux

package screenshot

import (
	"fmt"
	"time"
)

// GetCGWindowBoundsByPID - Linux implementation placeholder
// TODO: Implement using X11/XCB or Wayland protocols
func GetCGWindowBoundsByPID(pid int) (int, int, int, int, error) {
	return 0, 0, 0, 0, fmt.Errorf("CoreGraphics not available on Linux, will fall back to robotgo")
}

// GetCGWindowIDByPID - Linux implementation placeholder
func GetCGWindowIDByPID(pid int) (int, error) {
	return 0, fmt.Errorf("CoreGraphics not available on Linux")
}

// WaitForCGWindow - Linux implementation placeholder
// On Linux, robotgo's native window detection should work
func WaitForCGWindow(pid int, timeout time.Duration) error {
	return fmt.Errorf("CoreGraphics not available on Linux, will fall back to robotgo")
}

// CaptureWindowByCGBounds - Linux implementation placeholder
func CaptureWindowByCGBounds(pid int, output string, quality int) error {
	return fmt.Errorf("CoreGraphics not available on Linux, use robotgo methods")
}
