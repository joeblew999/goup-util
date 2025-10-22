package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/joeblew999/goup-util/pkg/buildcache"
	"github.com/joeblew999/goup-util/pkg/config"
	"github.com/joeblew999/goup-util/pkg/constants"
	"github.com/joeblew999/goup-util/pkg/icons"
	"github.com/joeblew999/goup-util/pkg/installer"
	"github.com/joeblew999/goup-util/pkg/project"
	"github.com/spf13/cobra"
)

// BuildOptions contains options for build commands
type BuildOptions struct {
	Force     bool
	CheckOnly bool
	SkipIcons bool
}

// Global build cache
var globalBuildCache *buildcache.Cache

// getBuildCache returns the global build cache, initializing if needed
func getBuildCache() *buildcache.Cache {
	if globalBuildCache == nil {
		cache, err := buildcache.NewCache(buildcache.GetDefaultCachePath())
		if err != nil {
			// If cache fails, create empty one (won't save)
			cache = &buildcache.Cache{}
		}
		globalBuildCache = cache
	}
	return globalBuildCache
}

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

		// Check for custom output directory flag first
		customOutput, _ := cmd.Flags().GetString("output")

		// Create and validate project with potential custom output
		var proj *project.GioProject
		var err error

		if customOutput != "" {
			// Use custom output directory
			proj, err = project.NewGioProjectWithOutput(appDir, customOutput)
		} else {
			// Use default behavior (artifacts in project directory)
			proj, err = project.NewGioProject(appDir)
		}

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
		force, _ := cmd.Flags().GetBool("force")
		checkOnly, _ := cmd.Flags().GetBool("check")

		// Create build options
		opts := BuildOptions{
			Force:     force,
			CheckOnly: checkOnly,
			SkipIcons: skipIcons,
		}

		switch platform {
		case "macos":
			return buildMacOS(proj.RootDir, proj.Name, outputDir, platform, opts)
		case "android":
			return buildAndroid(proj.RootDir, proj.Name, outputDir, platform, opts)
		case "ios":
			return buildIOS(proj.RootDir, proj.Name, outputDir, platform, opts, false)
		case "ios-simulator":
			return buildIOS(proj.RootDir, proj.Name, outputDir, "ios-simulator", opts, true)
		case "windows":
			return buildWindows(proj.RootDir, proj.Name, outputDir, platform, opts)
		case "all":
			return buildAll(proj.RootDir, proj.Name, outputDir, opts)
		}

		return nil
	},
}

func buildMacOS(appDir, appName, outputDir, platform string, opts BuildOptions) error {
	appPath := filepath.Join(outputDir, appName+".app")
	cache := getBuildCache()

	// Check if rebuild is needed
	if !opts.Force {
		needsRebuild, reason := cache.NeedsRebuild(appName, platform, appDir, appPath)

		if opts.CheckOnly {
			if needsRebuild {
				fmt.Printf("Rebuild needed: %s\n", reason)
				os.Exit(1)
			} else {
				fmt.Printf("Up to date: %s\n", appPath)
				os.Exit(0)
			}
		}

		if !needsRebuild {
			fmt.Printf("✓ %s for %s is up-to-date (use --force to rebuild)\n", appName, platform)
			return nil
		}

		fmt.Printf("Rebuilding: %s\n", reason)
	}

	fmt.Printf("Building %s for macOS...\n", appName)

	// Generate icons
	if !opts.SkipIcons {
		if err := generateIcons(appDir, "macos"); err != nil {
			cache.RecordBuild(appName, platform, appDir, appPath, false)
			return fmt.Errorf("failed to generate icons: %w", err)
		}
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Remove existing app bundle only if it exists
	if _, err := os.Stat(appPath); err == nil {
		os.RemoveAll(appPath)
	}

	// Build with gogio
	iconPath := filepath.Join(appDir, "icon-source.png")
	gogioCmd := exec.Command("gogio", "-target", "macos", "-arch", "arm64", "-icon", iconPath, "-o", appPath, appDir)
	gogioCmd.Stdout = os.Stdout
	gogioCmd.Stderr = os.Stderr

	if err := gogioCmd.Run(); err != nil {
		cache.RecordBuild(appName, platform, appDir, appPath, false)
		return fmt.Errorf("gogio build failed: %w", err)
	}

	// Record successful build
	cache.RecordBuild(appName, platform, appDir, appPath, true)

	fmt.Printf("✓ Built %s for macOS: %s\n", appName, appPath)
	return nil
}

func buildAndroid(appDir, appName, outputDir, platform string, opts BuildOptions) error {
	apkPath := filepath.Join(outputDir, appName+".apk")
	cache := getBuildCache()

	// Check if rebuild is needed
	if !opts.Force {
		needsRebuild, reason := cache.NeedsRebuild(appName, platform, appDir, apkPath)

		if opts.CheckOnly {
			if needsRebuild {
				fmt.Printf("Rebuild needed: %s\n", reason)
				os.Exit(1)
			} else {
				fmt.Printf("Up to date: %s\n", apkPath)
				os.Exit(0)
			}
		}

		if !needsRebuild {
			fmt.Printf("✓ %s for %s is up-to-date (use --force to rebuild)\n", appName, platform)
			return nil
		}

		fmt.Printf("Rebuilding: %s\n", reason)
	}

	fmt.Printf("Building %s for Android...\n", appName)

	// Generate icons
	if !opts.SkipIcons {
		if err := generateIcons(appDir, "android"); err != nil {
			cache.RecordBuild(appName, platform, appDir, apkPath, false)
			return fmt.Errorf("failed to generate icons: %w", err)
		}
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Use OS-specific SDK directory only
	sdkRoot := config.GetSDKDir()

	// Check for required Android components
	ndkPath := filepath.Join(sdkRoot, "ndk-bundle")
	if _, err := os.Stat(ndkPath); os.IsNotExist(err) {
		fmt.Printf("⚠️  Android NDK not found. Installing...\n")
		// Auto-install NDK
		if err := installNDK(sdkRoot); err != nil {
			cache.RecordBuild(appName, platform, appDir, apkPath, false)
			return fmt.Errorf("failed to install NDK: %w", err)
		}
	}

	// Set Android environment variables with absolute paths
	env := os.Environ()
	javaHome := filepath.Join(sdkRoot, "openjdk", "17", "jdk-17.0.11+9", "Contents", "Home")
	env = append(env, "JAVA_HOME="+javaHome)
	env = append(env, "ANDROID_SDK_ROOT="+sdkRoot)
	env = append(env, "ANDROID_HOME="+sdkRoot)
	env = append(env, "ANDROID_NDK_ROOT="+ndkPath)

	// Build with gogio
	gogioCmd := exec.Command("gogio", "-target", "android", "-o", apkPath, appDir)
	gogioCmd.Env = env
	gogioCmd.Stdout = os.Stdout
	gogioCmd.Stderr = os.Stderr

	if err := gogioCmd.Run(); err != nil {
		cache.RecordBuild(appName, platform, appDir, apkPath, false)
		return fmt.Errorf("gogio build failed: %w", err)
	}

	// Record successful build
	cache.RecordBuild(appName, platform, appDir, apkPath, true)

	fmt.Printf("✓ Built %s for Android: %s\n", appName, apkPath)
	return nil
}

func buildIOS(appDir, appName, outputDir, platform string, opts BuildOptions, simulator bool) error {
	target := "iOS device"
	if simulator {
		target = "iOS simulator"
	}

	appPath := filepath.Join(outputDir, appName+".app")
	cache := getBuildCache()

	// Check if rebuild is needed
	if !opts.Force {
		needsRebuild, reason := cache.NeedsRebuild(appName, platform, appDir, appPath)

		if opts.CheckOnly {
			if needsRebuild {
				fmt.Printf("Rebuild needed: %s\n", reason)
				os.Exit(1)
			} else {
				fmt.Printf("Up to date: %s\n", appPath)
				os.Exit(0)
			}
		}

		if !needsRebuild {
			fmt.Printf("✓ %s for %s is up-to-date (use --force to rebuild)\n", appName, platform)
			return nil
		}

		fmt.Printf("Rebuilding: %s\n", reason)
	}

	fmt.Printf("Building %s for %s...\n", appName, target)

	// Generate icons
	if !opts.SkipIcons {
		if err := generateIcons(appDir, "ios"); err != nil {
			cache.RecordBuild(appName, platform, appDir, appPath, false)
			return fmt.Errorf("failed to generate icons: %w", err)
		}
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Build with gogio
	gogioCmd := exec.Command("gogio", "-target", "ios", "-o", appPath, appDir)
	gogioCmd.Stdout = os.Stdout
	gogioCmd.Stderr = os.Stderr

	if err := gogioCmd.Run(); err != nil {
		cache.RecordBuild(appName, platform, appDir, appPath, false)
		return fmt.Errorf("gogio build failed: %w", err)
	}

	// Record successful build
	cache.RecordBuild(appName, platform, appDir, appPath, true)

	fmt.Printf("✓ Built %s for %s: %s\n", appName, target, appPath)
	return nil
}

func buildWindows(appDir, appName, outputDir, platform string, opts BuildOptions) error {
	exePath := filepath.Join(outputDir, appName+".exe")
	cache := getBuildCache()

	// Check if rebuild is needed
	if !opts.Force {
		needsRebuild, reason := cache.NeedsRebuild(appName, platform, appDir, exePath)

		if opts.CheckOnly {
			if needsRebuild {
				fmt.Printf("Rebuild needed: %s\n", reason)
				os.Exit(1)
			} else {
				fmt.Printf("Up to date: %s\n", exePath)
				os.Exit(0)
			}
		}

		if !needsRebuild {
			fmt.Printf("✓ %s for %s is up-to-date (use --force to rebuild)\n", appName, platform)
			return nil
		}

		fmt.Printf("Rebuilding: %s\n", reason)
	}

	fmt.Printf("Building %s for Windows...\n", appName)

	// Generate icons
	if !opts.SkipIcons {
		if err := generateIcons(appDir, "windows"); err != nil {
			cache.RecordBuild(appName, platform, appDir, exePath, false)
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
	iconPath := filepath.Join(appDir, "icon-source.png")
	gogioCmd := exec.Command("gogio", "-o", exePath, "-target", "windows", "-icon", iconPath, appDir)
	gogioCmd.Env = env
	gogioCmd.Stdout = os.Stdout
	gogioCmd.Stderr = os.Stderr
	gogioCmd.Dir = appDir

	if err := gogioCmd.Run(); err != nil {
		cache.RecordBuild(appName, platform, appDir, exePath, false)
		return fmt.Errorf("gogio build failed: %w", err)
	}

	// Record successful build
	cache.RecordBuild(appName, platform, appDir, exePath, true)

	fmt.Printf("✓ Built %s for Windows: %s\n", appName, exePath)
	return nil
}

// installNDK installs the Android NDK if not present
func installNDK(sdkRoot string) error {
	fmt.Printf("📦 Installing Android NDK...\n")
	
	// Use the installer package to install NDK
	ndkSDK := &installer.SDK{
		Name:        "Android NDK",
		Version:     "latest",
		InstallPath: "ndk-bundle",
	}
	
	cache, err := installer.NewCache(filepath.Join(config.GetCacheDir(), "cache.json"))
	if err != nil {
		return fmt.Errorf("failed to create cache: %w", err)
	}
	
	return installer.Install(ndkSDK, cache)
}

func buildAll(appDir, appName, outputDir string, opts BuildOptions) error {
	fmt.Printf("Building %s for all platforms...\n", appName)

	platforms := []string{"macos", "android", "ios-simulator", "windows"}

	for _, platform := range platforms {
		fmt.Printf("\n--- Building for %s ---\n", platform)
		switch platform {
		case "macos":
			if err := buildMacOS(appDir, appName, outputDir, platform, opts); err != nil {
				fmt.Printf("❌ Failed to build for %s: %v\n", platform, err)
			}
		case "android":
			if err := buildAndroid(appDir, appName, outputDir, platform, opts); err != nil {
				fmt.Printf("❌ Failed to build for %s: %v\n", platform, err)
			}
		case "ios-simulator":
			if err := buildIOS(appDir, appName, outputDir, platform, opts, true); err != nil {
				fmt.Printf("❌ Failed to build for %s: %v\n", platform, err)
			}
		case "windows":
			if err := buildWindows(appDir, appName, outputDir, platform, opts); err != nil {
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
		outputPath = filepath.Join(appDir, constants.BuildDir)
	case "ios", "macos":
		outputPath = filepath.Join(appDir, constants.BuildDir, "Assets.xcassets")
	case "windows":
		outputPath = filepath.Join(appDir, constants.BuildDir)
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
	buildCmd.Flags().String("output", "", "Custom output directory for build artifacts")
	buildCmd.Flags().Bool("force", false, "Force rebuild even if up-to-date")
	buildCmd.Flags().Bool("check", false, "Check if rebuild needed (exit 0=no, 1=yes)")
	rootCmd.AddCommand(buildCmd)
}

var skipIcons bool
