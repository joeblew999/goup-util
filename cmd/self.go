package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/joeblew999/goup-util/pkg/bootstrap"
	"github.com/spf13/cobra"
)

var selfCmd = &cobra.Command{
	Use:   "self",
	Short: "Manage goup-util itself",
	Long:  "Commands for building and updating goup-util itself.",
}

var selfBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build goup-util itself",
	Long:  "Build the goup-util binary from source. Useful for IAC and GitHub Actions.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return buildSelf()
	},
}

var selfUpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade goup-util to latest release",
	Long:  "Download and install the latest goup-util release from GitHub.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return updateSelf()
	},
}

var selfReleaseCmd = &cobra.Command{
	Use:   "release [patch|minor|major|v1.2.3]",
	Short: "Release goup-util",
	Long:  "Complete release process: test, build, commit, push, and tag.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return releaseSelf(args[0])
	},
}

var selfBootstrapTestCmd = &cobra.Command{
	Use:   "bootstrap-test",
	Short: "Test bootstrap scripts with local binaries",
	Long: `Test bootstrap scripts by generating them in LOCAL mode and running them.

This command:
1. Builds binaries if needed (go run . self build)
2. Generates bootstrap scripts in LOCAL mode
3. Runs the appropriate bootstrap script for current platform
4. Verifies installation worked

Use this to test bootstrap scripts before release.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return testBootstrap()
	},
}

func init() {
	rootCmd.AddCommand(selfCmd)
	selfCmd.AddCommand(selfBuildCmd)
	selfCmd.AddCommand(selfUpgradeCmd)
	selfCmd.AddCommand(selfReleaseCmd)
	selfCmd.AddCommand(selfBootstrapTestCmd)
}

func buildSelf() error {
	fmt.Println("🔨 Building goup-util for all architectures...")
	
	// Get current directory (where goup-util source is)
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	
	// Build for multiple architectures
	architectures := []struct {
		goos, goarch, suffix string
	}{
		{"darwin", "arm64", "darwin-arm64"},
		{"darwin", "amd64", "darwin-amd64"},
		{"linux", "amd64", "linux-amd64"},
		{"linux", "arm64", "linux-arm64"},
		{"windows", "amd64", "windows-amd64.exe"},
		{"windows", "arm64", "windows-arm64.exe"},
	}
	
	for _, arch := range architectures {
		outputPath := filepath.Join(currentDir, fmt.Sprintf("goup-util-%s", arch.suffix))

		buildCmd := exec.Command("go", "build", "-o", outputPath, ".")
		buildCmd.Env = append(os.Environ(),
			"GOOS="+arch.goos,
			"GOARCH="+arch.goarch,
			"CGO_ENABLED=0", // Disable CGO for cross-platform builds
		)
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr
		
		if err := buildCmd.Run(); err != nil {
			return fmt.Errorf("failed to build goup-util for %s/%s: %w", arch.goos, arch.goarch, err)
		}
		
		fmt.Printf("✅ Built goup-util-%s\n", arch.suffix)
	}

	// Generate bootstrap scripts with correct binary names
	fmt.Println("\n📝 Generating bootstrap scripts...")
	if err := generateBootstrapScripts(currentDir, architectures); err != nil {
		return fmt.Errorf("failed to generate bootstrap scripts: %w", err)
	}

	return nil
}

func updateSelf() error {
	fmt.Println("🔄 Checking for latest release...")
	
	// Get latest release info
	resp, err := http.Get("https://api.github.com/repos/joeblew999/goup-util/releases/latest")
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to fetch release info: %s", resp.Status)
	}
	
	var release struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name string `json:"name"`
			URL  string `json:"browser_download_url"`
		} `json:"assets"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("failed to parse release info: %w", err)
	}
	
	fmt.Printf("📦 Latest version: %s\n", release.TagName)
	
	// Determine platform-specific binary name
	var binaryName string
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	
	switch {
	case goos == "darwin" && goarch == "arm64":
		binaryName = "goup-util-darwin-arm64"
	case goos == "darwin" && goarch == "amd64":
		binaryName = "goup-util-darwin-amd64"
	case goos == "linux" && goarch == "amd64":
		binaryName = "goup-util-linux-amd64"
	case goos == "linux" && goarch == "arm64":
		binaryName = "goup-util-linux-arm64"
	case goos == "windows" && goarch == "amd64":
		binaryName = "goup-util-windows-amd64.exe"
	case goos == "windows" && goarch == "arm64":
		binaryName = "goup-util-windows-arm64.exe"
	default:
		return fmt.Errorf("unsupported platform: %s/%s", goos, goarch)
	}
	
	// Find matching asset
	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == binaryName {
			downloadURL = asset.URL
			break
		}
	}
	
	if downloadURL == "" {
		return fmt.Errorf("no binary found for %s/%s", goos, goarch)
	}
	
	fmt.Printf("⬇️  Downloading %s...\n", binaryName)
	
	// Download binary
	downloadResp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download binary: %w", err)
	}
	defer downloadResp.Body.Close()
	
	// Get install path
	installPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	
	// Create backup
	backupPath := installPath + ".backup"
	if err := os.Rename(installPath, backupPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}
	
	// Write new binary
	file, err := os.OpenFile(installPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		// Restore backup on error
		os.Rename(backupPath, installPath)
		return fmt.Errorf("failed to write binary: %w", err)
	}
	defer file.Close()
	
	if _, err := io.Copy(file, downloadResp.Body); err != nil {
		// Restore backup on error
		os.Rename(backupPath, installPath)
		return fmt.Errorf("failed to save binary: %w", err)
	}
	
	// Remove backup on success
	os.Remove(backupPath)
	
	fmt.Printf("✅ Updated to %s\n", release.TagName)
	fmt.Printf("🔄 Run 'goup-util --version' to verify\n")
	
	return nil
}

func releaseSelf(version string) error {
	fmt.Printf("🚀 Starting release process for %s...\n", version)
	
	// Run tests
	fmt.Println("✅ Running tests...")
	testCmd := exec.Command("go", "test", "./...")
	testCmd.Stdout = os.Stdout
	testCmd.Stderr = os.Stderr
	if err := testCmd.Run(); err != nil {
		return fmt.Errorf("tests failed: %w", err)
	}
	
	// Run race tests
	fmt.Println("✅ Running race tests...")
	raceCmd := exec.Command("go", "test", "-race", "./...")
	raceCmd.Stdout = os.Stdout
	raceCmd.Stderr = os.Stderr
	if err := raceCmd.Run(); err != nil {
		return fmt.Errorf("race tests failed: %w", err)
	}
	
	// Build self
	fmt.Println("🔨 Building goup-util...")
	if err := buildSelf(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}
	
	// Validate version format
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}
	if !regexp.MustCompile(`^v\d+\.\d+\.\d+$`).MatchString(version) {
		return fmt.Errorf("invalid version format: %s (use v1.2.3, patch, minor, or major)", version)
	}
	
	// Handle bump types
	if version == "patch" || version == "minor" || version == "major" {
		// Get current version
		currentTag, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
		if err != nil {
			currentTag = []byte("v0.0.0")
		}
		
		current := strings.TrimSpace(string(currentTag))
		version = bumpVersion(current, version)
	}
	
	// Check if working directory is clean
	if err := exec.Command("git", "diff-index", "--quiet", "HEAD", "--").Run(); err != nil {
		return fmt.Errorf("working directory is not clean. Please commit changes first")
	}
	
	// Create and push tag
	fmt.Printf("🏷️  Creating release %s...\n", version)
	if err := exec.Command("git", "tag", "-a", version, "-m", "Release "+version).Run(); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}
	
	fmt.Println("📤 Pushing tag...")
	if err := exec.Command("git", "push", "origin", version).Run(); err != nil {
		return fmt.Errorf("failed to push tag: %w", err)
	}
	
	fmt.Printf("✅ Release %s created and pushed!\n", version)
	fmt.Println("🔄 GitHub Actions will now build and create the release automatically")
	
	return nil
}

func bumpVersion(current, bumpType string) string {
	current = strings.TrimPrefix(current, "v")
	parts := strings.Split(current, ".")
	if len(parts) != 3 {
		return "v1.0.0"
	}
	
	major, _ := strconv.Atoi(parts[0])
	minor, _ := strconv.Atoi(parts[1])
	patch, _ := strconv.Atoi(parts[2])
	
	switch bumpType {
	case "major":
		major++
		minor = 0
		patch = 0
	case "minor":
		minor++
		patch = 0
	case "patch":
		patch++
	}
	
	return fmt.Sprintf("v%d.%d.%d", major, minor, patch)
}

func generateBootstrapScripts(baseDir string, architectures []struct{ goos, goarch, suffix string }) error {
	scriptsDir := filepath.Join(baseDir, "scripts")

	// Find macOS and Windows architectures
	var macOSArchs, windowsArchs []string
	for _, arch := range architectures {
		switch arch.goos {
		case "darwin":
			macOSArchs = append(macOSArchs, arch.goarch)
		case "windows":
			windowsArchs = append(windowsArchs, arch.goarch)
		}
	}

	// Use bootstrap package to generate scripts
	config := bootstrap.Config{
		Repo:           "joeblew999/goup-util",
		SupportedArchs: bootstrap.ArchsToString(append(macOSArchs, windowsArchs...)),
		MacOSArchs:     macOSArchs,
		WindowsArchs:   windowsArchs,
	}

	return bootstrap.Generate(scriptsDir, config)
}

func testBootstrap() error {
	fmt.Println("🧪 Testing bootstrap scripts...")
	fmt.Println()

	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Check if binaries exist, build if not
	binaryPattern := filepath.Join(currentDir, "goup-util-*")
	matches, err := filepath.Glob(binaryPattern)
	if err != nil || len(matches) == 0 {
		fmt.Println("📦 Binaries not found, building first...")
		if err := buildSelf(); err != nil {
			return err
		}
		fmt.Println()
	}

	// Generate bootstrap scripts in LOCAL mode
	fmt.Println("📝 Generating bootstrap scripts in LOCAL mode...")
	
	architectures := []struct {
		goos, goarch, suffix string
	}{
		{"darwin", "arm64", "darwin-arm64"},
		{"darwin", "amd64", "darwin-amd64"},
		{"linux", "amd64", "linux-amd64"},
		{"linux", "arm64", "linux-arm64"},
		{"windows", "amd64", "windows-amd64.exe"},
		{"windows", "arm64", "windows-arm64.exe"},
	}

	if err := generateBootstrapScriptsLocal(currentDir, architectures); err != nil {
		return fmt.Errorf("failed to generate bootstrap scripts: %w", err)
	}
	fmt.Println()

	// Run appropriate bootstrap for current platform
	fmt.Println("🚀 Running bootstrap script for current platform...")
	fmt.Println()

	scriptsDir := filepath.Join(currentDir, "scripts")
	
	switch runtime.GOOS {
	case "darwin":
		scriptPath := filepath.Join(scriptsDir, "macos-bootstrap.sh")
		fmt.Printf("Executing: %s\n\n", scriptPath)
		
		cmd := exec.Command("bash", scriptPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("bootstrap script failed: %w", err)
		}

	case "windows":
		scriptPath := filepath.Join(scriptsDir, "windows-bootstrap.ps1")
		fmt.Printf("Executing: %s\n\n", scriptPath)
		
		cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", scriptPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("bootstrap script failed: %w", err)
		}

	case "linux":
		return fmt.Errorf("linux bootstrap testing not yet implemented - run scripts/linux-bootstrap.sh manually")

	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	fmt.Println()
	fmt.Println("✅ Bootstrap test completed successfully!")
	return nil
}

func generateBootstrapScriptsLocal(baseDir string, architectures []struct{ goos, goarch, suffix string }) error {
	scriptsDir := filepath.Join(baseDir, "scripts")

	// Find macOS and Windows architectures
	var macOSArchs, windowsArchs []string
	for _, arch := range architectures {
		switch arch.goos {
		case "darwin":
			macOSArchs = append(macOSArchs, arch.goarch)
		case "windows":
			windowsArchs = append(windowsArchs, arch.goarch)
		}
	}

	// Use bootstrap package in LOCAL mode
	config := bootstrap.Config{
		Repo:           "joeblew999/goup-util",
		SupportedArchs: bootstrap.ArchsToString(append(macOSArchs, windowsArchs...)),
		MacOSArchs:     macOSArchs,
		WindowsArchs:   windowsArchs,
		UseLocal:       true,
		LocalBinDir:    baseDir, // Use current directory where binaries are
	}

	return bootstrap.Generate(scriptsDir, config)
}
