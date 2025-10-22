package self

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/joeblew999/goup-util/pkg/installer"
	"github.com/joeblew999/goup-util/pkg/self/output"
)

// BuildOptions contains options for the Build function.
type BuildOptions struct {
	UseLocal   bool // If true, generate bootstrap scripts for local testing
	Obfuscate  bool // If true, use garble to obfuscate the binary
}

// Build cross-compiles goup-util for all supported architectures.
// Generates binaries in the current directory and bootstrap scripts in scripts/.
func Build(opts BuildOptions) error {
	result := output.BuildResult{
		Binaries:        []string{},
		ScriptsGenerated: false,
		LocalMode:       opts.UseLocal,
	}

	// Get current directory (where goup-util source is)
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	result.OutputDir = currentDir

	// Check if garble is needed and available
	if opts.Obfuscate {
		if !installer.IsGarbleInstalled() {
			return fmt.Errorf("--obfuscate flag requires garble to be installed. Run: go run . install garble")
		}
		fmt.Println("🔒 Building with garble obfuscation...")
	}

	// Build for all supported architectures
	for _, arch := range SupportedArchitectures() {
		outputPath := filepath.Join(currentDir, fmt.Sprintf("goup-util-%s", arch.Suffix))

		var buildCmd *exec.Cmd
		if opts.Obfuscate {
			// Use garble build for obfuscation
			buildCmd = exec.Command("garble", "build", "-o", outputPath, ".")
		} else {
			// Normal go build
			buildCmd = exec.Command("go", "build", "-o", outputPath, ".")
		}

		buildCmd.Env = append(os.Environ(),
			"GOOS="+arch.GOOS,
			"GOARCH="+arch.GOARCH,
			"CGO_ENABLED=0", // Disable CGO for cross-platform builds
		)
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr

		if err := buildCmd.Run(); err != nil {
			return fmt.Errorf("failed to build goup-util for %s/%s: %w", arch.GOOS, arch.GOARCH, err)
		}

		result.Binaries = append(result.Binaries, fmt.Sprintf("goup-util-%s", arch.Suffix))
	}

	// Generate bootstrap scripts with correct binary names
	if err := generateBootstrapScripts(currentDir, opts); err != nil {
		return fmt.Errorf("failed to generate bootstrap scripts: %w", err)
	}

	result.ScriptsGenerated = true
	output.Print(result, "self build")
	return nil
}

// generateBootstrapScripts creates bootstrap shell/PowerShell scripts
func generateBootstrapScripts(baseDir string, opts BuildOptions) error {
	scriptsDir := filepath.Join(baseDir, "scripts")

	// Get supported architectures
	allArchs := SupportedArchitectures()
	macOSArchs := ArchsToGoArchList(FilterByOS(allArchs, "darwin"))
	windowsArchs := ArchsToGoArchList(FilterByOS(allArchs, "windows"))

	// Create bootstrap script config
	config := Config{
		Repo:           FullRepoName,
		SupportedArchs: ArchsToString(append(macOSArchs, windowsArchs...)),
		MacOSArchs:     macOSArchs,
		WindowsArchs:   windowsArchs,
		UseLocal:       opts.UseLocal,
		SetupCommand:   SetupCommand, // Single source of truth
	}

	// If local mode, set LocalBinDir
	if opts.UseLocal {
		config.LocalBinDir = baseDir
	}

	return Generate(scriptsDir, config)
}
