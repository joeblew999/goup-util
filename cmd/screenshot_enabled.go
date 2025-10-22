// +build screenshot

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/joeblew999/goup-util/pkg/screenshot"
	"github.com/spf13/cobra"
)

var screenshotCmd = &cobra.Command{
	Use:   "screenshot [output-file]",
	Short: "Capture screenshots using robotgo",
	Long: `Capture screenshots using robotgo for full multi-display support.

Features:
  - Multi-display capture and information
  - Precise region selection
  - Delayed capture (for menus/tooltips)
  - JPEG/PNG output formats

Note: Requires CGO. On macOS 10.15+, grant Screen Recording permission:
System Settings → Privacy & Security → Screen Recording

Examples:
  # Capture full screen
  goup-util screenshot output.png

  # Capture specific region
  goup-util screenshot --x 100 --y 100 -w 800 -H 600 region.png

  # Capture all displays
  goup-util screenshot --all --prefix display

  # Delayed capture (useful for menus/tooltips)
  goup-util screenshot --delay 3000 output.png

  # Get display information
  goup-util screenshot --info

  # JPEG with custom quality
  goup-util screenshot -q 95 output.jpg

Using: ` + screenshot.Platform(),
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flags
		x, _ := cmd.Flags().GetInt("x")
		y, _ := cmd.Flags().GetInt("y")
		width, _ := cmd.Flags().GetInt("width")
		height, _ := cmd.Flags().GetInt("height")
		allDisplays, _ := cmd.Flags().GetBool("all")
		prefix, _ := cmd.Flags().GetString("prefix")
		delay, _ := cmd.Flags().GetInt("delay")
		quality, _ := cmd.Flags().GetInt("quality")
		info, _ := cmd.Flags().GetBool("info")
		force, _ := cmd.Flags().GetBool("force")

		// Show display info
		if info {
			return showDisplayInfo()
		}

		// Determine output file
		var output string
		if len(args) > 0 {
			output = args[0]
		} else if allDisplays {
			// Use prefix for multi-display
			if prefix == "" {
				prefix = "display"
			}
		} else {
			// Generate default filename with timestamp
			timestamp := time.Now().Format("20060102-150405")
			output = fmt.Sprintf("screenshot-%s.png", timestamp)
		}

		// Check if output already exists (idempotency)
		if !force && output != "" && !allDisplays {
			if info, err := os.Stat(output); err == nil {
				absPath, _ := filepath.Abs(output)
				fmt.Printf("✓ Screenshot already exists: %s\n", absPath)
				fmt.Printf("  Size: %d bytes, Modified: %s\n", info.Size(), info.ModTime().Format("2006-01-02 15:04:05"))
				fmt.Println("  Use --force to overwrite")
				return nil
			}
		}

		// Create config
		cfg := screenshot.Config{
			Output:      output,
			X:           x,
			Y:           y,
			Width:       width,
			Height:      height,
			AllDisplays: allDisplays,
			Prefix:      prefix,
			Quality:     quality,
			Delay:       delay,
		}

		// Validate
		if !allDisplays && output == "" {
			return fmt.Errorf("output file required (or use --all with --prefix)")
		}

		// Capture
		if err := screenshot.Capture(cfg); err != nil {
			return fmt.Errorf("screenshot failed: %w", err)
		}

		// Success message
		if allDisplays {
			fmt.Printf("✓ Captured all displays with prefix: %s\n", prefix)
		} else {
			absPath, _ := filepath.Abs(output)
			fmt.Printf("✓ Screenshot saved: %s\n", absPath)
		}

		return nil
	},
}

func showDisplayInfo() error {
	displays, err := screenshot.GetInfo()
	if err != nil {
		return err
	}

	if len(displays) == 0 {
		fmt.Println("No display information available")
		fmt.Println("(Rebuild with -tags screenshot for full features)")
		return nil
	}

	fmt.Printf("Found %d display(s):\n\n", len(displays))
	for _, d := range displays {
		fmt.Printf("Display %d:\n", d.ID)
		fmt.Printf("  Position: (%d, %d)\n", d.X, d.Y)
		fmt.Printf("  Size: %dx%d\n", d.Width, d.Height)
		fmt.Println()
	}

	return nil
}

func init() {
	screenshotCmd.Flags().Int("x", 0, "X coordinate of region to capture")
	screenshotCmd.Flags().Int("y", 0, "Y coordinate of region to capture")
	screenshotCmd.Flags().IntP("width", "w", 0, "Width of region to capture")
	screenshotCmd.Flags().IntP("height", "H", 0, "Height of region to capture")
	screenshotCmd.Flags().Bool("all", false, "Capture all displays (requires robotgo)")
	screenshotCmd.Flags().String("prefix", "display", "Prefix for multi-display captures")
	screenshotCmd.Flags().IntP("delay", "d", 0, "Delay before capture in milliseconds")
	screenshotCmd.Flags().IntP("quality", "q", 90, "JPEG quality (1-100)")
	screenshotCmd.Flags().Bool("info", false, "Show display information")
	screenshotCmd.Flags().Bool("force", false, "Overwrite existing screenshot")

	rootCmd.AddCommand(screenshotCmd)
}
