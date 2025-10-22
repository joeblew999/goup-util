package bootstrap

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/macos-bootstrap.sh.tmpl
var macosTemplate string

//go:embed templates/windows-bootstrap.ps1.tmpl
var windowsTemplate string

// Config holds the configuration for generating bootstrap scripts
type Config struct {
	Repo            string   // e.g., "joeblew999/goup-util"
	SupportedArchs  string   // e.g., "arm64, amd64"
	MacOSArchs      []string // e.g., ["arm64", "amd64"]
	WindowsArchs    []string // e.g., ["amd64", "arm64"]
	LocalBinDir     string   // Optional: local directory with binaries for testing
	UseLocal        bool     // If true, use local binaries instead of GitHub releases
}

// Generate creates bootstrap scripts from templates
func Generate(outputDir string, config Config) error {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate macOS bootstrap
	if len(config.MacOSArchs) > 0 {
		macosPath := filepath.Join(outputDir, "macos-bootstrap.sh")
		if err := generateScript(macosPath, macosTemplate, config, 0755); err != nil {
			return fmt.Errorf("failed to generate macOS bootstrap: %w", err)
		}
		fmt.Printf("  ✓ Generated %s\n", macosPath)
	}

	// Generate Windows bootstrap
	if len(config.WindowsArchs) > 0 {
		windowsPath := filepath.Join(outputDir, "windows-bootstrap.ps1")
		if err := generateScript(windowsPath, windowsTemplate, config, 0644); err != nil {
			return fmt.Errorf("failed to generate Windows bootstrap: %w", err)
		}
		fmt.Printf("  ✓ Generated %s\n", windowsPath)
	}

	return nil
}

func generateScript(path, tmplContent string, config Config, perm os.FileMode) error {
	tmpl, err := template.New(filepath.Base(path)).Parse(tmplContent)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, config); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// ArchsToString converts a slice of architectures to a comma-separated string
func ArchsToString(archs []string) string {
	return strings.Join(archs, ", ")
}
