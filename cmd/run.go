package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/joeblew999/goup-util/pkg/project"
	"github.com/joeblew999/goup-util/pkg/utils"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [platform] [app-directory]",
	Short: "Build and run a Gio application",
	Long: `Build and run a Gio application for the specified platform.

This command builds the app (if needed) and launches it. The app is automatically
opened using the platform-specific path, so you don't need to know where it's built.

Platforms: macos, ios-simulator (macOS only for now)

Examples:
  goup-util run macos ./myapp
  goup-util run macos examples/hybrid-dashboard`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		platform := args[0]
		appDir := args[1]

		// Validate platform - only support platforms we can run locally
		validPlatforms := []string{"macos"}
		if runtime.GOOS == "darwin" {
			validPlatforms = append(validPlatforms, "ios-simulator")
		}

		if !utils.Contains(validPlatforms, platform) {
			return fmt.Errorf("cannot run %s apps on %s. Valid platforms: %v", platform, runtime.GOOS, validPlatforms)
		}

		// Create and validate project
		proj, err := project.NewGioProject(appDir)
		if err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}

		if err := proj.Validate(); err != nil {
			return fmt.Errorf("invalid project: %w", err)
		}

		// Get build flags
		force, _ := cmd.Flags().GetBool("force")
		skipIcons, _ := cmd.Flags().GetBool("skip-icons")
		schemes, _ := cmd.Flags().GetString("schemes")

		// Build the app
		opts := BuildOptions{
			Force:     force,
			SkipIcons: skipIcons,
			Schemes:   schemes,
		}

		switch platform {
		case "macos":
			if err := buildMacOS(proj, platform, opts); err != nil {
				return fmt.Errorf("build failed: %w", err)
			}
		case "ios-simulator":
			if err := buildIOS(proj, platform, opts, true); err != nil {
				return fmt.Errorf("build failed: %w", err)
			}
		}

		// Launch the app
		appPath := proj.GetOutputPath(platform)
		fmt.Printf("Launching %s...\n", appPath)

		switch platform {
		case "macos":
			return launchMacOSApp(appPath)
		case "ios-simulator":
			return launchIOSSimulator(appPath)
		}

		return nil
	},
}

func launchMacOSApp(appPath string) error {
	cmd := exec.Command("open", appPath)
	return cmd.Run()
}

func launchIOSSimulator(appPath string) error {
	// Boot a simulator if needed, then install and launch
	// First, try to open simulator
	bootCmd := exec.Command("open", "-a", "Simulator")
	if err := bootCmd.Run(); err != nil {
		fmt.Printf("Note: Could not open Simulator app: %v\n", err)
	}

	// Install the app
	installCmd := exec.Command("xcrun", "simctl", "install", "booted", appPath)
	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("failed to install on simulator: %w", err)
	}

	// Get the bundle ID from the app
	// For now, assume the app name matches the bundle ID base
	// In a real implementation, we'd read Info.plist
	fmt.Printf("App installed. Launch it from the Simulator.\n")
	return nil
}

func init() {
	runCmd.Flags().Bool("force", false, "Force rebuild even if up-to-date")
	runCmd.Flags().Bool("skip-icons", false, "Skip icon generation")
	runCmd.Flags().String("schemes", "", "Deep linking URI schemes")

	// Group for help organization
	runCmd.GroupID = "build"

	rootCmd.AddCommand(runCmd)
}
