package cmd

import (
	"github.com/joeblew999/goup-util/pkg/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joeblew999/goup-util/pkg/constants"
	"github.com/joeblew999/goup-util/pkg/packaging"
	"github.com/joeblew999/goup-util/pkg/project"
	"github.com/spf13/cobra"
)

var bundleCmd = &cobra.Command{
	Use:   "bundle [platform] [app-directory]",
	Short: "Create signed app bundles for distribution",
	Long: `Create properly signed and structured app bundles for distribution.
This includes:
- macOS: .app bundle with Info.plist, code signing, and entitlements
- Android: Signed APK (future)
- iOS: Signed IPA (future)
- Windows: Installer (future)

This is different from 'package' which just creates archives of built apps.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		platform := args[0]
		appDir := args[1]

		// Validate platform
		validPlatforms := []string{"macos", "android", "ios", "windows"}
		if !utils.Contains(validPlatforms, platform) {
			return fmt.Errorf("invalid platform: %s. Valid platforms: %v", platform, validPlatforms)
		}

		// Get flags
		bundleID, _ := cmd.Flags().GetString("bundle-id")
		version, _ := cmd.Flags().GetString("version")
		signingIdentity, _ := cmd.Flags().GetString("sign")
		outputDir, _ := cmd.Flags().GetString("output")
		entitlements, _ := cmd.Flags().GetBool("entitlements")

		// Create and validate project
		proj, err := project.NewGioProject(appDir)
		if err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}

		if err := proj.Validate(); err != nil {
			return fmt.Errorf("invalid project: %w", err)
		}

		switch platform {
		case "macos":
			return bundleMacOS(proj, bundleID, version, signingIdentity, outputDir, entitlements)
		case "android":
			return fmt.Errorf("android bundling not yet implemented")
		case "ios":
			return fmt.Errorf("ios bundling not yet implemented")
		case "windows":
			publisher, _ := cmd.Flags().GetString("publisher")
			createMSIX, _ := cmd.Flags().GetBool("create-msix")
			return bundleWindows(proj, bundleID, version, publisher, outputDir, createMSIX)
		}

		return nil
	},
}

func bundleMacOS(proj *project.GioProject, bundleID, version, signingIdentity, outputDir string, useEntitlements bool) error {
	fmt.Printf("Creating macOS bundle for %s...\n", proj.Name)

	// Set defaults
	if bundleID == "" {
		bundleID = fmt.Sprintf("com.example.%s", proj.Name)
	}
	if version == "" {
		version = "1.0.0"
	}
	if outputDir == "" {
		outputDir = filepath.Join(proj.RootDir, constants.DistDir)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Find the built binary - check multiple locations
	binDir := filepath.Join(proj.RootDir, constants.BinDir)
	var binaryPath string

	// Check if there's a standalone binary
	standaloneBinary := filepath.Join(binDir, proj.Name)
	if _, err := os.Stat(standaloneBinary); err == nil {
		binaryPath = standaloneBinary
	} else {
		// Check if there's an existing .app bundle with the binary inside
		appBundlePath := filepath.Join(binDir, proj.Name+".app")
		binaryInApp := filepath.Join(appBundlePath, "Contents", "MacOS", proj.Name)
		if _, err := os.Stat(binaryInApp); err == nil {
			binaryPath = binaryInApp
			fmt.Println("ℹ️  Found binary in existing .app bundle, will create new signed bundle")
		} else {
			return fmt.Errorf("binary not found in %s or %s\nRun 'goup-util build macos %s' first",
				standaloneBinary, binaryInApp, proj.RootDir)
		}
	}

	// Create bundle config
	config := packaging.MacOSBundleConfig{
		Name:            proj.Name,
		DisplayName:     toDisplayName(proj.Name),
		BundleID:        bundleID,
		Version:         version,
		BuildNumber:     "1",
		BinaryPath:      binaryPath,
		OutputDir:       outputDir,
		SigningIdentity: signingIdentity,
		Entitlements:    useEntitlements,
	}

	// Create the bundle
	if err := packaging.CreateMacOSBundle(config); err != nil {
		return fmt.Errorf("failed to create bundle: %w", err)
	}

	fmt.Println()
	fmt.Println("🎯 Next steps:")
	fmt.Println("   1. Test the app: open", filepath.Join(outputDir, proj.Name+".app"))
	fmt.Println("   2. Grant permissions if needed (System Settings → Privacy & Security)")
	fmt.Println("   3. Package for distribution:")
	fmt.Printf("      task package:macos:dmg OR\n")
	fmt.Printf("      goup-util package macos %s\n", proj.RootDir)

	return nil
}

func bundleWindows(proj *project.GioProject, bundleID, version, publisher, outputDir string, createMSIX bool) error {
	fmt.Printf("Creating Windows bundle for %s...\n", proj.Name)

	// Set defaults
	if bundleID == "" {
		bundleID = proj.Name
	}
	if version == "" {
		version = "1.0.0.0"
	}
	if publisher == "" {
		publisher = fmt.Sprintf("CN=%s", proj.Name)
	}
	if outputDir == "" {
		outputDir = filepath.Join(proj.RootDir, constants.DistDir)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Find the built binary
	binDir := filepath.Join(proj.RootDir, constants.BinDir)
	binaryPath := filepath.Join(binDir, proj.Name+".exe")

	// Check if binary exists
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return fmt.Errorf("binary not found: %s\nRun 'goup-util build windows %s' first", binaryPath, proj.RootDir)
	}

	// Check for assets directory
	assetsDir := filepath.Join(proj.RootDir, "assets")
	if _, err := os.Stat(assetsDir); os.IsNotExist(err) {
		assetsDir = "" // Will generate placeholders
	}

	// Create bundle config
	config := packaging.WindowsBundleConfig{
		Name:                 bundleID,
		Publisher:            publisher,
		PublisherDisplayName: toDisplayName(proj.Name),
		DisplayName:          toDisplayName(proj.Name),
		Description:          fmt.Sprintf("%s application", toDisplayName(proj.Name)),
		Version:              version,
		BinaryPath:           binaryPath,
		OutputDir:            outputDir,
		AssetsDir:            assetsDir,
		CreateMSIX:           createMSIX,
	}

	// Create the bundle
	if err := packaging.CreateWindowsBundle(config); err != nil {
		return fmt.Errorf("failed to create bundle: %w", err)
	}

	fmt.Println()
	fmt.Println("🎯 Next steps:")
	if createMSIX {
		fmt.Println("   1. Test the MSIX:", filepath.Join(outputDir, proj.Name+".msix"))
		fmt.Println("   2. Install: Add-AppxPackage", filepath.Join(outputDir, proj.Name+".msix"))
	} else {
		fmt.Println("   1. Copy .staging directory to Windows machine")
		fmt.Println("   2. Run: goup-util bundle --create-msix windows", proj.RootDir)
	}
	fmt.Println("   3. Package for distribution:")
	fmt.Printf("      goup-util package windows %s\n", proj.RootDir)

	return nil
}

// toDisplayName converts a name like "goup-util" to "Goup Util"
func toDisplayName(name string) string {
	// Simple title case - can be improved
	return name
}

func init() {
	bundleCmd.Flags().String("bundle-id", "", "Bundle identifier (e.g., com.example.myapp)")
	bundleCmd.Flags().String("version", "1.0.0", "Version string")
	bundleCmd.Flags().String("sign", "", "Code signing identity (empty for auto-detect)")
	bundleCmd.Flags().String("output", "", "Output directory (default: .dist/)")
	bundleCmd.Flags().Bool("entitlements", true, "Use entitlements for hardened runtime (macOS)")
	bundleCmd.Flags().String("publisher", "", "Publisher for Windows MSIX (e.g., CN=MyCompany)")
	bundleCmd.Flags().Bool("create-msix", false, "Create MSIX package (Windows-only, requires msix toolkit)")

	rootCmd.AddCommand(bundleCmd)
}
