package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/joeblew99/goup-util/pkg/gitignore"
	"github.com/spf13/cobra"
)

var gitignoreCmd = &cobra.Command{
	Use:   "gitignore [project-path]",
	Short: "Manage .gitignore files for Gio projects",
	Long:  `Manage .gitignore files for Gio projects. Shows status and can generate appropriate gitignore patterns.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath := "."
		if len(args) > 0 {
			projectPath = args[0]
		}

		// Check if project path exists
		if _, err := os.Stat(projectPath); os.IsNotExist(err) {
			return fmt.Errorf("project path does not exist: %s", projectPath)
		}

		gi := gitignore.New(projectPath)
		if err := gi.Load(); err != nil {
			return fmt.Errorf("failed to load .gitignore: %w", err)
		}

		// Show status
		fmt.Printf("📁 Project: %s\n", projectPath)
		info := gi.Info()

		if gi.Exists {
			fmt.Printf("✅ .gitignore exists with %d lines\n", info["lines"])
			if managedSection, ok := info["managed_section"].(bool); ok && managedSection {
				fmt.Printf("🔧 Has goup-util managed section\n")
			} else {
				fmt.Printf("⚠️  No goup-util managed section\n")
			}
		} else {
			fmt.Printf("❌ No .gitignore file found\n")
		}

		// Show current patterns
		if gi.Exists && len(gi.Lines) > 0 {
			fmt.Printf("\n📝 Current patterns:\n")
			for _, line := range gi.Lines {
				if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "#") {
					fmt.Printf("   %s\n", line)
				}
			}
		}

		// Show recommended patterns for Gio projects
		fmt.Printf("\n💡 Recommended patterns for Gio projects:\n")
		recommended := []string{
			"# Build artifacts",
			".bin/",
			"*.exe",
			"*.app",
			"*.apk",
			"*.ipa",
			"*.msix",
			"",
			"# Generated icons and assets",
			"icon.png",
			"icon.ico",
			"icon.icns",
			"Assets.xcassets/",
			"drawable-*/",
			"*.syso",
			"",
			"# OS files",
			".DS_Store",
			"Thumbs.db",
		}

		for _, pattern := range recommended {
			if pattern == "" {
				fmt.Println()
				continue
			}

			status := "✨"
			if gi.Exists && gi.HasPattern(strings.TrimPrefix(pattern, "# ")) {
				status = "✅"
			}
			fmt.Printf("   %s %s\n", status, pattern)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(gitignoreCmd)
}
