package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/joeblew999/goup-util/pkg/config"
	"github.com/joeblew999/goup-util/pkg/icons"
	"github.com/joeblew999/goup-util/pkg/project"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build [platform] [app-directory]",
	Short: "Build Gio applications for different platforms",
	Long:  "High-level command to generate icons and build Gio applications for various platforms.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		platform := args[0]
		appDir := args[1]

		// Validate platform
		validPlatforms := []string{"macos", "android", "ios", "ios-simulator", "windows", "all"}
		if !contains(validPlatforms, platform) {
			return fmt.Errorf("invalid platform: %s. Valid platforms: %v", platform, validPlatforms)
		}

		// Create and validate project
		proj, err := project.NewGioProject(appDir)
		if err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}

		if err := proj.Validate(); err != nil {
			return fmt.Errorf("invalid project: %w", err)
		}

		// Get paths from project
		paths := proj.Paths()
		outputDir := paths.Output

		// Get flags
		skipIcons, _ := cmd.Flags().GetBool("skip-icons")

		switch platform {
		case "macos":
			return buildMacOS(proj.RootDir, proj.Name, outputDir, skipIcons)
		case "android":
			return buildAndroid(proj.RootDir, proj.Name, outputDir, skipIcons)
		case "ios":
			return buildIOS(proj.RootDir, proj.Name, outputDir, skipIcons, false)
		case "ios-simulator":
			return buildIOS(proj.RootDir, proj.Name, outputDir, skipIcons, true)
		case "windows":
			return buildWindows(proj.RootDir, proj.Name, outputDir, skipIcons)
		case "all":
			return buildAll(proj.RootDir, proj.Name, outputDir, skipIcons)
		}

		return nil
	},
}

func buildMacOS(appDir, appName, outputDir string, skipIcons bool) error {
	fmt.Printf("Building %s for macOS...\n", appName)

	// Generate icons
	if !skipIcons {
		if err := generateIcons(appDir, "macos"); err != nil {
			return fmt.Errorf("failed to generate icons: %w", err)
		}
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Remove existing app bundle
	appPath := filepath.Join(outputDir, appName+".app")
	os.RemoveAll(appPath)

	// Build with gogio
	iconPath := filepath.Join(appDir, "icon-source.png")
	gogioCmd := exec.Command("gogio", "-target", "macos", "-arch", "arm64", "-icon", iconPath, "-o", appPath, appDir)
	gogioCmd.Stdout = os.Stdout
	gogioCmd.Stderr = os.Stderr

	if err := gogioCmd.Run(); err != nil {
		return fmt.Errorf("gogio build failed: %w", err)
	}

	fmt.Printf("✓ Built %s for macOS: %s\n", appName, appPath)
	return nil
}

func buildAndroid(appDir, appName, outputDir string, skipIcons bool) error {
	fmt.Printf("Building %s for Android...\n", appName)

	// Generate icons
	if !skipIcons {
		if err := generateIcons(appDir, "android"); err != nil {
			return fmt.Errorf("failed to generate icons: %w", err)
		}
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Use OS-specific SDK directory only
	sdkRoot := config.GetSDKDir()

	// Set Android environment variables with absolute paths
	env := os.Environ()
	javaHome := filepath.Join(sdkRoot, "openjdk", "17", "jdk-17.0.11+9", "Contents", "Home")
	env = append(env, "JAVA_HOME="+javaHome)
	env = append(env, "ANDROID_SDK_ROOT="+sdkRoot)

	// Build with gogio
	apkPath := filepath.Join(outputDir, appName+".apk")
	gogioCmd := exec.Command("gogio", "-target", "android", "-o", apkPath, appDir)
	gogioCmd.Env = env
	gogioCmd.Stdout = os.Stdout
	gogioCmd.Stderr = os.Stderr

	if err := gogioCmd.Run(); err != nil {
		return fmt.Errorf("gogio build failed: %w", err)
	}

	fmt.Printf("✓ Built %s for Android: %s\n", appName, apkPath)
	return nil
}

func buildIOS(appDir, appName, outputDir string, skipIcons bool, simulator bool) error {
	target := "iOS device"
	if simulator {
		target = "iOS simulator"
	}
	fmt.Printf("Building %s for %s...\n", appName, target)

	// Generate icons
	if !skipIcons {
		if err := generateIcons(appDir, "ios"); err != nil {
			return fmt.Errorf("failed to generate icons: %w", err)
		}
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Build with gogio
	appPath := filepath.Join(outputDir, appName+".app")
	gogioCmd := exec.Command("gogio", "-target", "ios", "-o", appPath, appDir)
	gogioCmd.Stdout = os.Stdout
	gogioCmd.Stderr = os.Stderr

	if err := gogioCmd.Run(); err != nil {
		return fmt.Errorf("gogio build failed: %w", err)
	}

	fmt.Printf("✓ Built %s for %s: %s\n", appName, target, appPath)
	return nil
}

func buildWindows(appDir, appName, outputDir string, skipIcons bool) error {
	fmt.Printf("Building %s for Windows...\n", appName)

	// Generate icons
	if !skipIcons {
		if err := generateIcons(appDir, "windows"); err != nil {
			return fmt.Errorf("failed to generate icons: %w", err)
		}
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Set Windows environment
	env := os.Environ()
	env = append(env, "GOOS=windows")
	env = append(env, "GOARCH=arm64")

	// Build with gogio
	exePath := filepath.Join(outputDir, appName+".exe")
	iconPath := filepath.Join(appDir, "icon-source.png")
	gogioCmd := exec.Command("gogio", "-o", exePath, "-target", "windows", "-icon", iconPath, appDir)
	gogioCmd.Env = env
	gogioCmd.Stdout = os.Stdout
	gogioCmd.Stderr = os.Stderr
	gogioCmd.Dir = appDir

	if err := gogioCmd.Run(); err != nil {
		return fmt.Errorf("gogio build failed: %w", err)
	}

	fmt.Printf("✓ Built %s for Windows: %s\n", appName, exePath)
	return nil
}

func buildAll(appDir, appName, outputDir string, skipIcons bool) error {
	fmt.Printf("Building %s for all platforms...\n", appName)

	platforms := []string{"macos", "android", "ios-simulator", "windows"}

	for _, platform := range platforms {
		fmt.Printf("\n--- Building for %s ---\n", platform)
		switch platform {
		case "macos":
			if err := buildMacOS(appDir, appName, outputDir, skipIcons); err != nil {
				fmt.Printf("❌ Failed to build for %s: %v\n", platform, err)
			}
		case "android":
			if err := buildAndroid(appDir, appName, outputDir, skipIcons); err != nil {
				fmt.Printf("❌ Failed to build for %s: %v\n", platform, err)
			}
		case "ios-simulator":
			if err := buildIOS(appDir, appName, outputDir, skipIcons, true); err != nil {
				fmt.Printf("❌ Failed to build for %s: %v\n", platform, err)
			}
		case "windows":
			if err := buildWindows(appDir, appName, outputDir, skipIcons); err != nil {
				fmt.Printf("❌ Failed to build for %s: %v\n", platform, err)
			}
		}
	}

	fmt.Printf("\n✓ Build complete for all platforms\n")
	return nil
}

func generateIcons(appDir, platform string) error {
	// Ensure source icon exists
	sourceIconPath, err := icons.EnsureSourceIcon(appDir)
	if err != nil {
		return err
	}

	// Generate platform-specific icons
	var outputPath string
	switch platform {
	case "android":
		outputPath = appDir
	case "ios", "macos":
		outputPath = filepath.Join(appDir, "Assets.xcassets")
	case "windows":
		outputPath = filepath.Join(appDir, ".bin")
		platform = "windows-msix" // Use the correct platform name
	default:
		return nil // Skip icon generation for unknown platforms
	}

	fmt.Printf("Generating %s icons...\n", platform)
	return icons.Generate(icons.Config{
		InputPath:  sourceIconPath,
		OutputPath: outputPath,
		Platform:   platform,
	})
}

// Remove the old generateTestIcon function since it's now in the icons package

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func init() {
	buildCmd.Flags().BoolVar(&skipIcons, "skip-icons", false, "Skip icon generation")
	rootCmd.AddCommand(buildCmd)
}

var skipIcons bool
