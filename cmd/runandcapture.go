// +build screenshot

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/joeblew999/goup-util/pkg/screenshot"
	"github.com/spf13/cobra"
)

var runAndCaptureCmd = &cobra.Command{
	Use:   "run-and-capture <app-dir> <output-file>",
	Short: "Run Gio app and capture screenshot",
	Long: `Run a Gio application, wait for its window to appear, and capture a screenshot.

This automates the workflow of:
1. Launch the app
2. Wait for window to appear
3. Optionally resize window (if --preset or --width/--height specified)
4. Capture screenshot of the window
5. Stop the app

Examples:
  # Run app and capture screenshot
  goup-util run-and-capture examples/hybrid-dashboard screenshot.png

  # Use App Store preset size
  goup-util run-and-capture --preset macos-retina examples/hybrid-dashboard screenshot.png

  # Custom size
  goup-util run-and-capture --width 1280 --height 800 examples/hybrid-dashboard screenshot.png`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		appDir := args[0]
		output := args[1]

		// Get flags
		presetName, _ := cmd.Flags().GetString("preset")
		width, _ := cmd.Flags().GetInt("width")
		height, _ := cmd.Flags().GetInt("height")
		quality, _ := cmd.Flags().GetInt("quality")
		waitTime, _ := cmd.Flags().GetInt("wait")

		// Resolve absolute paths
		absAppDir, err := filepath.Abs(appDir)
		if err != nil {
			return fmt.Errorf("failed to resolve app directory: %w", err)
		}

		absOutput, err := filepath.Abs(output)
		if err != nil {
			return fmt.Errorf("failed to resolve output path: %w", err)
		}

		// Ensure output directory exists
		outputDir := filepath.Dir(absOutput)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		// Handle preset
		if presetName != "" {
			preset, ok := screenshot.GetPreset(presetName)
			if !ok {
				return fmt.Errorf("unknown preset: %s (use --list-presets to see available)", presetName)
			}
			fmt.Printf("Using preset: %s (%dx%d)\n", preset.Name, preset.Width, preset.Height)
			width = preset.Width
			height = preset.Height
		}

		// Build the app first to get a direct binary
		fmt.Printf("Building app in %s...\n", absAppDir)
		binaryPath := filepath.Join(absAppDir, "app-temp")
		cmdBuild := exec.Command("go", "build", "-o", binaryPath, ".")
		cmdBuild.Dir = absAppDir
		cmdBuild.Stdout = os.Stdout
		cmdBuild.Stderr = os.Stderr

		if err := cmdBuild.Run(); err != nil {
			return fmt.Errorf("failed to build app: %w", err)
		}

		// Launch the binary directly
		fmt.Printf("Launching %s...\n", binaryPath)
		cmdRun := exec.Command(binaryPath)
		cmdRun.Dir = absAppDir
		cmdRun.Stdout = os.Stdout
		cmdRun.Stderr = os.Stderr

		if err := cmdRun.Start(); err != nil {
			os.Remove(binaryPath)
			return fmt.Errorf("failed to launch app: %w", err)
		}

		pid := cmdRun.Process.Pid
		fmt.Printf("✓ Launched app (PID %d)\n", pid)

		// Ensure we kill the app and clean up when done
		defer func() {
			if cmdRun.Process != nil {
				cmdRun.Process.Kill()
				fmt.Printf("✓ Stopped app\n")
			}
			os.Remove(binaryPath)
		}()

		// Wait for window to appear
		fmt.Printf("Waiting for app to initialize...\n")
		timeout := time.Duration(waitTime) * time.Millisecond
		err = screenshot.WaitForWindow(pid, timeout)

		// If window detection fails, fall back to full screen capture
		if err != nil {
			fmt.Printf("⚠ Window detection failed (robotgo may not support Gio windows on this platform)\n")
			fmt.Printf("⚠ Falling back to full screen capture\n")
			fmt.Printf("⚠ Please manually position the app window before screenshot\n")

			// Give user time to position window
			fmt.Printf("Waiting 3 seconds for window positioning...\n")
			time.Sleep(3 * time.Second)

			// Capture full screen instead
			fmt.Printf("Capturing full screen...\n")
			if err := screenshot.CaptureDesktop(absOutput, quality); err != nil {
				return fmt.Errorf("failed to capture screenshot: %w", err)
			}
		} else {
			fmt.Printf("✓ Window appeared\n")

			// Give app time to render
			time.Sleep(1 * time.Second)

			// Note: robotgo doesn't have window positioning functions
			// So we just capture the window as-is
			if width > 0 && height > 0 {
				fmt.Printf("Note: Window sizing requested (%dx%d) but robotgo doesn't support resizing\n", width, height)
				fmt.Printf("      Will capture window at its current size\n")
			}

			// Capture screenshot by PID
			fmt.Printf("Capturing window screenshot...\n")
			if err := screenshot.CaptureWindowByPID(pid, absOutput, quality); err != nil {
				return fmt.Errorf("failed to capture screenshot: %w", err)
			}
		}

		fmt.Printf("✓ Screenshot saved: %s\n", absOutput)
		return nil
	},
}

func init() {
	runAndCaptureCmd.Flags().String("preset", "", "App Store preset size (e.g., macos-retina, iphone-6.9)")
	runAndCaptureCmd.Flags().Int("width", 0, "Window width (note: robotgo doesn't support resizing)")
	runAndCaptureCmd.Flags().Int("height", 0, "Window height (note: robotgo doesn't support resizing)")
	runAndCaptureCmd.Flags().IntP("quality", "q", 90, "JPEG quality (1-100)")
	runAndCaptureCmd.Flags().IntP("wait", "w", 5000, "Max wait time for window in milliseconds")

	rootCmd.AddCommand(runAndCaptureCmd)
}
