package cmd

import (
	"github.com/joeblew999/goup-util/pkg/self"
	"github.com/spf13/cobra"
)

var selfCmd = &cobra.Command{
	Use:   "self",
	Short: "Manage goup-util itself",
	Long:  "Commands for building, installing, upgrading, and releasing goup-util itself.",
}

var (
	buildLocal    bool // Flag for local mode
	buildObfuscate bool // Flag for garble obfuscation
)

var selfBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build goup-util binaries",
	Long:  "Cross-compile goup-util for all supported architectures and generate bootstrap scripts.",
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
	Short: "Test bootstrap scripts locally",
	Long:  "Generate and test bootstrap scripts in local mode to verify they work before releasing.",
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
	Short: "Upgrade goup-util to latest release",
	Long:  "Download and install the latest goup-util release from GitHub.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return self.DownloadAndInstallLatest(self.FullRepoName)
	},
}

var selfReleaseCmd = &cobra.Command{
	Use:   "release [patch|minor|major|v1.2.3]",
	Short: "Release goup-util",
	Long:  "Complete release process: test, build, commit, push, and tag.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return self.Release(args[0])
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
	selfBuildCmd.Flags().BoolVar(&buildObfuscate, "obfuscate", false, "Use garble to obfuscate binaries (requires: go run . install garble)")
}
