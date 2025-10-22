package cmd

import (
	"github.com/joeblew999/goup-util/pkg/self"
	"github.com/spf13/cobra"
)

var selfCmd = &cobra.Command{
	Use:   "self",
	Short: "Manage goup-util itself",
	Long: `Commands for building, installing, upgrading, and releasing goup-util itself.

Information Commands:
  version  - Show current version
  status   - Check installation and available updates
  doctor   - Validate dependencies (Homebrew, git, go, task)

Installation Commands:
  setup     - Full setup: install dependencies + goup-util to system PATH
  upgrade   - Download and install latest release from GitHub
  uninstall - Remove goup-util from system PATH

Development Commands:
  build   - Cross-compile binaries for all platforms (outputs to .dist/)
  test    - Test bootstrap scripts locally before releasing
  release - Create git tag and push (triggers GitHub Actions to build and release)`,
}

var (
	buildLocal     bool // Flag for local mode
	buildObfuscate bool // Flag for garble obfuscation
)

var selfBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build goup-util binaries for all platforms",
	Long: `Cross-compile goup-util for all supported architectures and generate bootstrap scripts.

Output: All artifacts are placed in .dist/ directory
- Binaries: .dist/goup-util-<platform>
- Scripts: .dist/macos-bootstrap.sh, .dist/windows-bootstrap.ps1

Flags:
  --local      Generate scripts that use local binaries (for testing)
  --obfuscate  Use garble to obfuscate code (auto-installs if needed)

This is a LOCAL build command - it does NOT create releases or push to GitHub.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := self.BuildOptions{
			UseLocal:  buildLocal,
			Obfuscate: buildObfuscate,
		}
		return self.Build(opts)
	},
}

var selfUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove goup-util from system path",
	Long:  "Remove goup-util binary from system path. Dependencies (Homebrew, git, go, task) are NOT removed.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return self.UninstallSelf()
	},
}

var selfVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show goup-util version",
	Long:  "Display the currently installed version of goup-util.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return self.ShowVersion()
	},
}

var selfStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check installation status and updates",
	Long:  "Check if goup-util is installed, show version, and check for available updates.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return self.ShowStatus()
	},
}

var selfDoctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Validate installation and dependencies",
	Long:  "Check that goup-util and all dependencies (Homebrew, git, go, task) are properly installed.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return self.Doctor()
	},
}

var selfTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test bootstrap scripts locally before releasing",
	Long: `Generate and test bootstrap scripts in local mode to verify they work.

This command:
1. Builds goup-util with --local flag (scripts use local binaries)
2. Verifies bootstrap scripts exist and contain expected content
3. Tests that the binary executes correctly

Run this before 'self release' to catch issues early.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return self.TestBootstrap()
	},
}

var selfSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Full setup: install dependencies + goup-util",
	Long:  "Install system dependencies (Homebrew/winget, git, go, task) and install goup-util to system path. This is what bootstrap scripts call.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Install dependencies first
		if err := self.InstallDeps(); err != nil {
			return err
		}
		// Then install binary
		return self.InstallSelf()
	},
}

var selfUpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Download and install latest release from GitHub",
	Long: `Download and install the latest goup-util release from GitHub.

This downloads the pre-built binary for your platform from the GitHub Releases page
and installs it to your system PATH (~/.local/bin/ or ~/bin/).

Use this command to update goup-util after a new release has been published.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return self.DownloadAndInstallLatest(self.FullRepoName)
	},
}

var selfReleaseCmd = &cobra.Command{
	Use:   "release [patch|minor|major|v1.2.3]",
	Short: "Create git tag and trigger GitHub Actions release",
	Long: `Prepare and trigger a release by creating a git tag.

This command does the following locally:
1. Runs tests (self test)
2. Builds with obfuscation (self build --obfuscate)
3. Commits .dist/ artifacts
4. Creates a git tag (e.g., v1.1.0)
5. Pushes tag to GitHub

GitHub Actions then automatically:
- Builds obfuscated binaries for all platforms
- Creates a GitHub Release
- Uploads artifacts (.dist/goup-util-*, .dist/*.sh, .dist/*.ps1)

Version options (defaults to 'minor'):
  patch      - Increment patch version (1.0.0 → 1.0.1)
  minor      - Increment minor version (1.0.0 → 1.1.0)
  major      - Increment major version (1.0.0 → 2.0.0)
  v1.2.3     - Use specific version

IMPORTANT: This does NOT upload files directly. It triggers GitHub Actions by pushing a tag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		version := "minor" // Default to minor release
		if len(args) == 1 {
			version = args[0]
		}
		return self.Release(version)
	},
}

func init() {
	rootCmd.AddCommand(selfCmd)

	// Information commands
	selfCmd.AddCommand(selfVersionCmd)
	selfCmd.AddCommand(selfStatusCmd)
	selfCmd.AddCommand(selfDoctorCmd)

	// Installation commands
	selfCmd.AddCommand(selfSetupCmd)
	selfCmd.AddCommand(selfUpgradeCmd)
	selfCmd.AddCommand(selfUninstallCmd)

	// Development commands
	selfCmd.AddCommand(selfBuildCmd)
	selfCmd.AddCommand(selfTestCmd)
	selfCmd.AddCommand(selfReleaseCmd)

	// Add flags
	selfBuildCmd.Flags().BoolVar(&buildLocal, "local", false, "Generate bootstrap scripts for local testing (uses local binaries instead of GitHub releases)")
	selfBuildCmd.Flags().BoolVar(&buildObfuscate, "obfuscate", false, "Use garble to obfuscate binaries (auto-installs garble if needed)")
}
